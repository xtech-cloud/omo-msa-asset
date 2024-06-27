package cache

import (
	"errors"
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy"
	"omo.msa.asset/proxy/nosql"
	"omo.msa.asset/tool"
	"time"
)

type ThumbInfo struct {
	ID       uint64 `json:"-"`
	Created  int64
	Updated  int64
	Status   uint32
	Probably float32
	Similar  float32
	Blur     float32
	UID      string `json:"uid"`
	Creator  string
	Operator string

	Face     string
	File     string
	Owner    string
	Asset    string
	Meta     string
	User     string
	Quote    string
	Group    string
	data     []byte
	Location proxy.LocationInfo
}

func CreateThumb(asset, owner, bs64, quote, group, operator string, bts []byte, info *DetectFace) *ThumbInfo {
	db := new(ThumbInfo)
	db.UID = primitive.NewObjectID().Hex()
	db.ID = nosql.GetThumbNextID()
	db.Created = time.Now().Unix()
	db.data = bts
	db.Creator = operator
	db.Operator = operator
	db.User = ""
	db.Quote = quote
	db.File = ""
	db.Asset = asset
	db.Blur = info.Quality.Blur
	db.Owner = owner
	db.Probably = info.Probability
	db.Similar = 0
	db.Group = group
	db.Location = info.Location
	db.Meta = bs64
	db.Status = 0
	return db
}

func (mine *cacheContext) GetUserThumbsByQuote(quote string, assets []string) []*ThumbInfo {
	dbs, err := nosql.GetThumbsByQuote(quote)
	if err != nil {
		return nil
	}
	list := make([]*ThumbInfo, 0, len(dbs))
	users := make([]string, 0, len(dbs))
	for _, db := range dbs {
		if len(db.User) > 0 {
			if len(assets) > 0 {
				if tool.HasItem(assets, db.Asset) && !tool.HasItem(users, db.User) {
					users = append(users, db.User)
					info := new(ThumbInfo)
					info.initInfo(db)
					list = append(list, info)
				}
			} else {
				if !tool.HasItem(users, db.User) {
					users = append(users, db.User)
					info := new(ThumbInfo)
					info.initInfo(db)
					list = append(list, info)
				}
			}
		}
	}
	return list
}

func (mine *cacheContext) GetThumbsByQuote(quote string) []*ThumbInfo {
	dbs, err := nosql.GetThumbsByQuote(quote)
	if err != nil {
		return nil
	}
	list := make([]*ThumbInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(ThumbInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) BindFaceEntity(user, entity, operator string) error {
	dbs, err := nosql.GetThumbsByUser(user)
	if err != nil {
		return err
	}
	for _, db := range dbs {
		er := nosql.UpdateThumbUser(db.UID.Hex(), entity, operator)
		if er != nil {
			return er
		}
	}

	return nil
}

func (mine *ThumbInfo) initInfo(db *nosql.Thumb) {
	mine.ID = db.ID
	mine.UID = db.UID.Hex()
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Owner = db.Owner
	mine.Asset = db.Asset
	mine.Probably = db.Probably
	mine.Similar = db.Similar
	mine.Blur = db.Blur
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Meta = db.Meta
	mine.User = db.User
	mine.Quote = db.Quote
	mine.File = db.File
	mine.Group = db.Group
	mine.Status = db.Status
	mine.Location = db.Location
	if len(mine.Group) < 1 {
		asset := cacheCtx.GetAsset(db.Asset)
		if asset != nil {
			mine.Group = asset.CheckFaceGroup()
			_ = nosql.UpdateThumbGroup(mine.UID, mine.Group)
		}
	}
}

func (mine *ThumbInfo) save() error {
	file := tool.CreateUUID()
	db := new(nosql.Thumb)
	db.UID, _ = primitive.ObjectIDFromHex(mine.UID)
	db.ID = mine.ID
	db.Created = mine.Created
	db.Creator = mine.Creator
	db.Operator = mine.Operator
	db.User = mine.User
	db.Quote = mine.Quote
	db.File = file
	db.Asset = mine.Asset
	db.Blur = mine.Blur
	db.Owner = mine.Owner
	db.Probably = mine.Probably
	db.Similar = 0
	db.Group = mine.Group
	db.Location = mine.Location
	db.Meta = mine.Meta
	db.Status = mine.Status
	er := nosql.CreateThumb(db)
	if er == nil {
		go uploadToQiNiu(file, mine.data)
	}
	logger.Warn("try save thumb of asset = " + db.Asset + " and thumb = " + mine.UID + "; user = " + db.User)
	return er
}

func (mine *ThumbInfo) UpdateBase(owner string, similar float32) error {
	err := nosql.UpdateThumbBase(mine.UID, owner, similar)
	if err == nil {
		mine.Owner = owner
		mine.Similar = similar
	}
	return err
}

//从人脸库里面搜索相似人脸的用户
func (mine *ThumbInfo) SearchUsers() ([]*UserResult, error) {
	req := new(FaceSearchReq)
	req.Type = ImageTypeBase64
	req.Image = mine.Meta
	req.Groups = mine.Group
	req.Quality = QualityNone
	req.MaxUser = 10
	req.Threshold = 80
	result, err, code := searchFaceByOne(req)
	if err != nil {
		if code == ErrorCodeQPSLimit {
			cacheCtx.addPendingThumb(mine, true)
		} else if code == ErrorCodeNotMatch {
			//_ = mine.RegisterFace(mine.User, mine.Group)
			//_ = mine.save()
			return nil, nil
		} else {
			mine.Status = uint32(BD_DetectFailed)
		}
		return nil, err
	}
	return result.Users, nil
}

//人脸认证：当前人脸和指定的用户的人脸是否一致
func (mine *ThumbInfo) Identification(user, group string) (*UserResult, error) {
	req := new(FaceSearchReq)
	req.Type = ImageTypeBase64
	req.Image = mine.Meta
	req.Groups = group
	req.Quality = QualityNone
	req.MaxUser = 1
	req.User = user
	req.Threshold = 80
	result, err, code := searchFaceByOne(req)
	if err != nil {
		if code == ErrorCodeQPSLimit {

		}
		return nil, err
	}
	if len(result.Users) < 1 {
		return nil, nil
	}
	return result.Users[0], nil
}

//把该人脸注册到人脸数据库中
func (mine *ThumbInfo) RegisterFace(user, group string) error {
	if len(group) < 1 {
		return errors.New("the group is empty")
	}
	mine.Status = uint32(BD_Detection)
	id := user
	if len(id) < 2 {
		id = fmt.Sprintf("temp_%d", mine.ID)
	}
	req := new(FaceAddReq)
	req.Type = ImageTypeBase64
	req.Image = mine.Meta
	req.Group = group
	req.User = id
	req.Quality = QualityNone
	req.Meta = fmt.Sprintf(`"user":"%s", "thumb":"%s"`, id, mine.UID)
	_, code, err := registerUserFace(req)
	if err != nil && code != ErrorCodeFaceExist {
		return err
	}

	if len(mine.User) < 2 {
		mine.User = id
	}
	return nil
}

func (mine *ThumbInfo) UpdateInfo(meta, operator string) error {
	err := nosql.UpdateThumbMeta(mine.UID, meta, operator)
	if err == nil {
		mine.Meta = meta
		mine.Operator = operator
	}
	return err
}

func (mine *ThumbInfo) BindEntity(entity, operator string) error {
	dbs, err := nosql.GetThumbsByUser(mine.User)
	if err != nil {
		return err
	}
	for _, db := range dbs {
		er := nosql.UpdateThumbUser(db.UID.Hex(), entity, operator)
		if er != nil {
			return er
		}
	}
	return err
}

func (mine *ThumbInfo) UpdateUser(user, operator string) error {
	err := nosql.UpdateThumbUser(mine.UID, user, operator)
	if err == nil {
		mine.User = user
		mine.Operator = operator
	}
	return err
}
