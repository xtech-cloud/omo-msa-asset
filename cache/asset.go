package cache

import (
	"errors"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"omo.msa.asset/config"
	"omo.msa.asset/proxy/nosql"
	"strings"
	"time"
)

const (
	AssetTypePerson       = 0
	AssetTypeGroup        = 1 //组织
	AssetTypeWindowModel  = 2 //windows 平台模型
	AssetTypeAndroidModel = 3 //web 平台模型
	AssetTypeAudio        = 4
	AssetTypeVideo        = 5
	AssetTypePortrait     = 6 //系统头像
	AssetTypeIcon         = 7 //图标库
	AssetTypeCertify      = 8 //证书
)

const (
	AssetScopePersonal = 0
	AssetScopeOrg      = 1
	AssetScopeSystem   = 2
)

const UP_QINIU = "qiniu"

const (
	StatusPrivate uint8 = 0
	StatusPending uint8 = 1 //待审
	StatusPublish uint8 = 2 //审核通过
	StatusVisible uint8 = 3 //可展示
)

type AssetInfo struct {
	Type     uint8
	Status   uint8
	Scope    uint8
	Width    uint32
	Height   uint32
	Weight   uint32
	Size     uint64
	Created  int64
	Updated  int64
	Code     int    //内部状态码
	ID       uint64 `json:"-"`
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Remark   string
	Meta     string
	Creator  string
	Operator string

	Owner    string
	UUID     string //file 云存储文件名
	Version  string
	Format   string
	MD5      string
	Language string
	Quote    string //引用的对象

	// 快照，中图
	Snapshot string
	// 封面小图
	Small string

	Links []string //关联的实体
	Tags  []string
}

func (mine *cacheContext) CreateAsset(info *pb.ReqAssetAdd) (*AssetInfo, error) {
	db := new(nosql.Asset)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAssetNextID()
	db.Created = time.Now().Unix()
	db.Creator = info.Operator
	db.Operator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = info.Owner
	db.Type = uint8(info.Type)
	db.Size = info.Size
	db.UUID = info.Uuid
	db.Format = info.Format
	db.MD5 = info.Md5
	db.Version = info.Version
	db.Language = info.Language
	db.Snapshot = info.Snapshot
	db.Small = info.Small
	db.Width = info.Width
	db.Height = info.Height
	db.Meta = info.Meta
	db.Weight = 0
	db.Tags = info.Tags
	db.Scope = uint8(info.Scope)
	db.Quote = info.Quote
	db.Status = StatusPrivate
	db.Links = make([]string, 0, 1)
	db.Meta = info.Meta
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}

	err := nosql.CreateAsset(db)
	if err == nil {
		tmp := new(AssetInfo)
		tmp.initInfo(db)
		if tmp.SupportFace() {
			go validateAsset(tmp)
		}
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetAsset(uid string) *AssetInfo {
	if uid == "" {
		return nil
	}
	db, err := nosql.GetAsset(uid)
	if err == nil {
		info := new(AssetInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetAssetByKey(key string) *AssetInfo {
	if key == "" {
		return nil
	}
	db, err := nosql.GetAssetByKey(key)
	if err == nil {
		info := new(AssetInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetAssetsByOwner(uid string) []*AssetInfo {
	array, err := nosql.GetAssetsByOwner(uid)
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetPublishAssetsByOwner(uid string, st uint32) []*AssetInfo {
	array, err := nosql.GetAssetsByOwnerStatus(uid, uint8(st))
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAssetCount(quote string) uint32 {
	return uint32(nosql.GetAssetsCountByQuote(quote))
}

func (mine *cacheContext) GetAssetCountByOwnerCreator(quote, creator string) uint32 {
	return uint32(nosql.GetAssetsCountByOwnerCreator(quote, creator))
}

func (mine *cacheContext) GetAssetCountByQuoteCreator(quote, creator string) uint32 {
	return uint32(nosql.GetAssetsCountByQuoteCreator(quote, creator))
}

func (mine *cacheContext) GetAssetsByQuote(owner, quote string, page, num uint32) (uint32, uint32, []*AssetInfo) {
	if quote == "" {
		return 0, 0, make([]*AssetInfo, 0, 1)
	}
	var dbs []*nosql.Asset
	var err error
	start, number := getPageStart(page, num)
	var total int64 = 0
	if len(owner) > 1 {
		dbs, err = nosql.GetAssetsByOwnerQuote(owner, quote, int64(start), int64(number))
		total = nosql.GetAssetsCountByOwnerQuote(owner, quote)
	} else {
		dbs, err = nosql.GetAssetsByQuote(quote, int64(start), int64(number))
		total = nosql.GetAssetsCountByQuote(quote)
	}

	pages := math.Ceil(float64(total) / float64(number))
	if err != nil {
		return 0, 0, make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, num)
	for _, asset := range dbs {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return uint32(total), uint32(pages), list
}

func (mine *cacheContext) GetAssetsByOwnerType(owner string, tp int) []*AssetInfo {
	array, err := nosql.GetAssetsByOwnerType(owner, uint8(tp))
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAssetsByQuoteStatus(quote string, st, page, num uint32) (uint32, uint32, []*AssetInfo) {
	start, number := getPageStart(page, num)
	total := nosql.GetAssetsCountByQuoteStatus(quote, st)
	array, err := nosql.GetAssetsByQuoteStatus(quote, uint8(st), int64(start), int64(number))
	pages := math.Ceil(float64(total) / float64(number))
	if err != nil {
		return 0, 0, make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return uint32(total), uint32(pages), list
}

func (mine *cacheContext) GetAssetsByQuoteCreator(quote, creator string, page, num uint32) (uint32, uint32, []*AssetInfo) {
	start, number := getPageStart(page, num)
	total := nosql.GetAssetsCountByQuoteCreator(quote, creator)
	pages := math.Ceil(float64(total) / float64(number))
	array, err := nosql.GetAssetsByQuoteCreator(quote, creator, int64(start), int64(number))

	if err != nil {
		return 0, 0, make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return uint32(total), uint32(pages), list
}

func (mine *cacheContext) GetAssetsByType(tp int) []*AssetInfo {
	array, err := nosql.GetAssetsByType(uint8(tp))
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) UpdateAssetsStatus(arr []string, st uint32, operator string) error {
	if arr == nil {
		return nil
	}
	for _, uid := range arr {
		_ = nosql.UpdateAssetStatus(uid, operator, uint8(st))
	}
	return nil
}

func (mine *cacheContext) PublishAssetsEntity(entity, operator string) error {
	if entity == "" {
		return errors.New("the entity is empty")
	}
	assets, _ := nosql.GetAssetsByOwner(entity)
	for _, asset := range assets {
		if asset.Status != StatusPublish {
			_ = nosql.UpdateAssetStatus(asset.UID.Hex(), operator, StatusPublish)
		}
	}
	return nil
}

func (mine *cacheContext) BatchUpdateScope(list []string) error {
	for _, owner := range list {
		assets, _ := nosql.GetAssetsByOwner(owner)
		for _, asset := range assets {
			_ = nosql.UpdateAssetScope(asset.UID.Hex(), asset.Operator, uint8(AssetScopeOrg))
		}
	}

	return nil
}

func (mine *cacheContext) GetAssetsByRegex(key string, from, to int64) []*AssetInfo {
	array, err := nosql.GetAssetsByRegex(key, from, to)
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAssetsByCreator(uid string) []*AssetInfo {
	array, err := nosql.GetAssetsByCreator(uid)
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAssetsByException(page, number uint32) (uint32, uint32, []*AssetInfo) {
	array, err := nosql.GetAllAssets()
	if err != nil {
		return 0, 0, make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		if asset.Creator == asset.Owner {
			info := new(AssetInfo)
			info.initInfo(asset)
			list = append(list, info)
		}
	}
	return checkPage(page, number, list)
}

func (mine *cacheContext) GetAssetsByLink(link string) []*AssetInfo {
	array, err := nosql.GetAssetsByLink(link)
	if err != nil {
		return make([]*AssetInfo, 0, 1)
	}
	list := make([]*AssetInfo, 0, len(array))
	for _, asset := range array {
		info := new(AssetInfo)
		info.initInfo(asset)
		list = append(list, info)
	}
	return list
}

func (mine *AssetInfo) initInfo(db *nosql.Asset) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Meta = db.Meta
	mine.Weight = db.Weight

	mine.Size = db.Size
	mine.UUID = db.UUID
	mine.Type = db.Type
	mine.Scope = db.Scope
	mine.Owner = db.Owner
	mine.Version = db.Version
	mine.MD5 = db.MD5
	mine.Format = db.Format
	mine.Language = db.Language
	mine.Snapshot = db.Snapshot
	mine.Small = db.Small
	mine.Quote = db.Quote

	mine.Width = db.Width
	mine.Height = db.Height
	mine.Weight = db.Weight
	mine.Status = db.Status
	mine.Links = db.Links
	mine.Tags = db.Tags
	mine.Code = db.Code

	//if mine.Code == BD_Conclusion {
	//	if mine.GetThumbCount() > 0 {
	//		_ = nosql.UpdateAssetCode(mine.UID, BD_Detection)
	//		mine.Code = BD_Detection
	//	} else {
	//		cacheCtx.addPendingAsset(mine)
	//	}
	//}
}

func (mine *AssetInfo) GetThumbCount() uint32 {
	return nosql.GetThumbCountByAsset(mine.UID)
}

func (mine *AssetInfo) GetThumbs() ([]*ThumbInfo, error) {
	array, err := nosql.GetThumbsByAsset(mine.UID)
	if err != nil {
		return nil, err
	}
	list := make([]*ThumbInfo, 0, len(array))
	for _, thumb := range array {
		tmp := new(ThumbInfo)
		tmp.initInfo(thumb)
		list = append(list, tmp)
	}
	return list, nil
}

func (mine *AssetInfo) CheckFaceGroup() string {
	group := FaceGroupDefault
	if mine.Scope == AssetScopeOrg {
		group = mine.Owner
	}
	return group
}

func (mine *AssetInfo) Remove(operator string) error {
	if mine.Type == AssetTypePortrait || mine.Type == AssetTypeIcon {
		return errors.New("the asset of type can not remove")
	}
	err := mine.ToRecycle(operator)
	if err == nil {
		_ = nosql.RemoveAsset(mine.UID)
		_ = nosql.RemoveThumbsByAsset(mine.UID, operator)
	}
	return err
}

func (mine *AssetInfo) getMinURL() (string, string) {
	if mine.Snapshot != "" {
		return mine.Snapshot, GetURL(mine.Snapshot, true)
	} else {
		return mine.UUID, GetURL(mine.UUID, true)
	}
}

func (mine *AssetInfo) SupportFace() bool {
	if mine.Type > AssetTypeWindowModel {
		return false
	}
	arr := []string{"png", "jpg", "jpeg", "bmp"}
	for _, s := range arr {
		format := strings.ToLower(mine.Format)
		if strings.Contains(format, s) {
			return true
		}
	}
	return false
}

func (mine *AssetInfo) ToRecycle(operator string) error {
	db := new(nosql.Recycle)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecycleNextID()
	db.Created = time.Now().Unix()
	db.Creator = mine.Creator
	db.Scavenger = operator
	db.Operator = mine.Operator
	db.Name = mine.Name
	db.Remark = mine.Remark
	db.Owner = mine.Owner
	db.Type = mine.Type
	db.Size = mine.Size
	db.UUID = mine.UUID
	db.Format = mine.Format
	db.MD5 = mine.MD5
	db.Version = mine.Version
	db.Language = mine.Language
	db.Snapshot = mine.Snapshot
	db.Small = mine.Small
	db.Width = mine.Width
	db.Height = mine.Height
	db.Meta = mine.Meta
	db.Weight = mine.Weight
	db.Status = mine.Status
	db.Quote = mine.Quote
	db.Links = mine.Links
	return nosql.CreateRecycle(db)
}

func (mine *AssetInfo) UpdateSnapshot(operator, snapshot string) error {
	err := nosql.UpdateAssetSnapshot(mine.UID, snapshot, operator)
	if err == nil {
		mine.Snapshot = snapshot
	}
	return err
}

func (mine *AssetInfo) UpdateSmall(operator, small string) error {
	err := nosql.UpdateAssetSmall(mine.UID, small, operator)
	if err == nil {
		mine.Small = small
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateBase(operator, name, remark string) error {
	err := nosql.UpdateAssetBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateMeta(operator, meta string) error {
	err := nosql.UpdateAssetMeta(mine.UID, meta, operator)
	if err == nil {
		mine.Meta = meta
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateWeight(weight uint32, operator string) error {
	err := nosql.UpdateAssetWeight(mine.UID, operator, weight)
	if err == nil {
		mine.Weight = weight
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateAssetStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateLinks(operator string, links []string) error {
	err := nosql.UpdateAssetLinks(mine.UID, operator, links)
	if err == nil {
		mine.Links = links
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateType(st uint8, operator string) error {
	err := nosql.UpdateAssetType(mine.UID, operator, st)
	if err == nil {
		mine.Type = st
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateAssetTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateOwner(operator, owner string) error {
	err := nosql.UpdateAssetOwner(mine.UID, owner, operator)
	if err == nil {
		mine.Owner = owner
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateQuote(operator, quote string) error {
	err := nosql.UpdateAssetQuote(mine.UID, quote, operator)
	if err == nil {
		mine.Quote = quote
		mine.Operator = operator
	}
	return err
}

func (mine *AssetInfo) UpdateLanguage(lan, operator string) error {
	err := nosql.UpdateAssetLanguage(mine.UID, operator, lan)
	if err == nil {
		mine.Language = lan
		mine.Operator = operator
	}
	return err
}

func GetURL(key string, cdn bool) string {
	if len(key) < 2 {
		return ""
	}
	if strings.Contains(key, "http") {
		return key
	}
	domain := config.Schema.Storage.Domain
	if !cdn {
		domain = config.Schema.Storage.Source
	}
	if config.Schema.Storage.Type == UP_QINIU {
		if config.Schema.Storage.ACM > 0 {
			mac := qbox.NewMac(config.Schema.Storage.AccessKey, config.Schema.Storage.SecretKey)
			return storage.MakePrivateURL(mac, domain, key, config.Schema.Storage.Period)
		} else {
			return storage.MakePublicURL(domain, key)
		}
	} else {
		return ""
	}
}

func (mine *AssetInfo) URL() string {
	return GetURL(mine.UUID, true)
}

func (mine *AssetInfo) SourceURL() string {
	if mine.Snapshot != "" {
		return GetURL(mine.Snapshot, false)
	}
	return GetURL(mine.UUID, false)
}

func (mine *AssetInfo) SnapshotURL() string {
	return GetURL(mine.Snapshot, true)
}

func (mine *AssetInfo) SmallImageURL() string {
	return GetURL(mine.Small, false)
}

func (mine *AssetInfo) HadThumbByFace(face string) bool {
	info := mine.GetThumbByFace(face)
	if info == nil {
		return false
	}
	return true
}

func (mine *AssetInfo) GetThumbByFace(face string) *ThumbInfo {
	db, err := nosql.GetThumbByFace(mine.UID, face)
	if err != nil {
		return nil
	}
	tmp := new(ThumbInfo)
	tmp.initInfo(db)
	return tmp
}

func (mine *AssetInfo) GetThumb(uid string) *ThumbInfo {
	db, err := nosql.GetThumb(uid)
	if err != nil {
		return nil
	}
	tmp := new(ThumbInfo)
	tmp.initInfo(db)
	return tmp
}

func (mine *AssetInfo) hadThumb(uid string) bool {
	info := mine.GetThumb(uid)
	if info == nil {
		return false
	}
	return true
}

func (mine *AssetInfo) RemoveThumb(uid, operator string) error {
	if !mine.hadThumb(uid) {
		return nil
	}
	return nosql.RemoveThumb(uid, operator)
}

func (mine *AssetInfo) CreateThumb(file, operator, owner string, score, similar, blur float32) (*ThumbInfo, error) {
	db := new(nosql.Thumb)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetThumbNextID()
	db.Created = time.Now().Unix()
	db.Creator = operator
	db.Operator = operator
	db.File = file
	db.Asset = mine.UID
	db.Blur = blur
	db.Owner = owner
	db.Probably = score
	db.Similar = similar
	err := nosql.CreateThumb(db)
	if err == nil {
		info := new(ThumbInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}
