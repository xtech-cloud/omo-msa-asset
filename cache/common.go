package cache

import (
	"fmt"
	"github.com/micro/go-micro/v2/logger"
	"omo.msa.asset/proxy/nosql"
	"time"
)

const MaxTask int = 8

type cacheContext struct {
	//assetPool *ants.Pool
	//thumbPool *ants.Pool
	assets     []*AssetInfo //待处理的人脸asset
	thumbs     []*ThumbInfo //待处理的人脸
	assetIndex int
	thumbIndex int
	//thumbCount int
}

var cacheCtx *cacheContext

func TestDetectFaces() {
	//url := "https://rdpdown.suii.cn/000c0f54-3dd7-40c6-aa2b-f67378947978"
	//url := "https://rdpdown.suii.cn/00278e27e030ac05"
	arr := []string{"667bbf07b6f93288c1761eec",
		"667bbf07b6f93288c1761eed",
		"667bbf07b6f93288c1761eee",
		"667bbf07b6f93288c1761eef",
		"667bbf07b6f93288c1761ef0",
		"667bbf07b6f93288c1761ef1",
		"667bbf07b6f93288c1761ef2",
		"667bbf07b6f93288c1761ef3",
		"667bbf07b6f93288c1761ef4",
		"667bbf07b6f93288c1761ef5",
		"667bbf08b6f93288c1761ef8",
		"667bbf08b6f93288c1761ef6",
		"667bbf08b6f93288c1761ef7",
		"667bbf08b6f93288c1761ef9",
		"667bbf08b6f93288c1761efb",
		"667bbf08b6f93288c1761efa",
		"667bbf08b6f93288c1761efc",
		"667bbf08b6f93288c1761efd",
		"667bbf08b6f93288c1761efe",
		"667bbf09b6f93288c1761eff",
		"667bbf09b6f93288c1761f00",
		"667bbf09b6f93288c1761f01",
		"667bbf09b6f93288c1761f02"}
	for _, uid := range arr {
		//thumb := cacheCtx.GetThumb(uid)
		asset := cacheCtx.GetAsset(uid)
		cacheCtx.addPendingAsset(asset)
		//_, url := asset.getMinURL()
		//_, buf, err := downloadAsset(url)
		//if err == nil {
		//	saveImage(buf.Bytes(), fmt.Sprintf("files/img/asset-%s.jpg", uid))
		//	_, bts, er := clipImageFace(buf, thumb.Location)
		//	if er == nil {
		//		saveImage(bts, fmt.Sprintf("files/img/thumb-%s.jpg", uid))
		//	}
		//}
	}

	//fmt.Println("test complete!!!!1")
	//asset := cacheCtx.GetAsset("6678ec2b92bf1c2e15633fc5")
	//validateAsset(asset)
}

func (mine *cacheContext) CheckThumbs() {
	dbs, _ := nosql.GetThumbsByUser("")
	for _, db := range dbs {
		if db.Status == 0 {
			info := new(ThumbInfo)
			info.initInfo(db)
			mine.addPendingThumb(info, false)
		}
	}
	for i := 0; i < MaxTask; i += 1 {
		go mine.checkPendingThumb()
	}
}

func (mine *cacheContext) initPool() {
	//mine.assetPool, _ = ants.NewPool(MaxTask, ants.WithOptions(ants.Options{PreAlloc: true}))
	//mine.thumbPool, _ = ants.NewPool(MaxTask, ants.WithOptions(ants.Options{PreAlloc: true}))
}

func (mine *cacheContext) addPendingAsset(asset *AssetInfo) {
	if asset == nil {
		return
	}
	if mine.assets == nil {
		mine.assetIndex = 0
		mine.assets = make([]*AssetInfo, 0, 100)
	}
	for _, item := range mine.assets {
		if item.UID == asset.UID {
			return
		}
	}
	mine.assets = append(mine.assets, asset)
	logger.Warn(fmt.Sprintf("addPendingAsset... the asset uid = %s", asset.UID))
	//err := mine.assetPool.Submit(mine.checkPendingAsset)
	//if err != nil {
	//	logger.Warn("pool submit error = " + err.Error())
	//}
	go mine.checkPendingAsset()
}

func (mine *cacheContext) checkPendingAsset() {
	if mine.assetIndex > MaxTask {
		return
	}
	info := mine.firstPendingAsset()
	if info == nil {
		return
	}

	logger.Warn(fmt.Sprintf("checkPendingAsset... the id = %d; and uid = %s", mine.assetIndex, info.UID))
	mine.assetIndex += 1
	//if info.Code == BD_Conclusion {
	_, url := info.getMinURL()
	group := info.CheckFaceGroup()
	_ = checkFaceGroup(group)
	er, code := checkFaces(info.UID, info.Owner, url, group, info.Quote, info.Creator)
	if er != nil {
		if code == ErrorCodeQPSLimit {
			mine.addPendingAsset(info)
		}
		logger.Warn("check faces failed that uid = " + info.UID + " and msg = " + er.Error())
	}
	//}
	mine.assetIndex -= 1
	//err := mine.assetPool.Submit(mine.checkPendingAsset)
	//if err != nil {
	//	logger.Warn("pool submit error = " + err.Error())
	//}
	time.Sleep(time.Second * 1)
	go mine.checkPendingAsset()
}

func (mine *cacheContext) firstPendingAsset() *AssetInfo {
	if len(mine.assets) < 1 {
		return nil
	}
	arr := mine.assets[:0]
	info := mine.assets[0]
	for i, v := range mine.assets {
		if i != 0 {
			arr = append(arr, v)
		}
	}
	mine.assets = arr
	return info
}

func (mine *cacheContext) addPendingThumb(info *ThumbInfo, check bool) {
	if info == nil {
		return
	}
	if mine.thumbs == nil {
		mine.thumbIndex = 0
		mine.thumbs = make([]*ThumbInfo, 0, 100)
	}
	for _, item := range mine.thumbs {
		if item.UID == info.UID {
			return
		}
	}
	mine.thumbs = append(mine.thumbs, info)
	if check {
		//mine.thumbPool.Submit(mine.checkPendingThumb)
		go mine.checkPendingThumb()
	}
}

func (mine *cacheContext) checkPendingThumb() {
	if mine.thumbIndex > MaxTask {
		return
	}

	info := mine.firstPendingThumb()
	if info == nil {
		return
	}

	logger.Warn(fmt.Sprintf("checkPendingThumb... the id = %d and uid = %s", mine.thumbIndex, info.UID))
	mine.thumbIndex += 1
	users, er := info.SearchUsers()
	if er != nil {
		//fmt.Println(fmt.Sprintf("search user face (%s) from group of %s, that err = %s", info.UID, info.Group, er.Error()))
		logger.Warn(fmt.Sprintf("search asset(%s) user face (%s) from group of %s, that err = %s", info.Asset, info.UID, info.Group, er.Error()))
	} else {
		if len(users) > 0 {
			_ = info.RegisterFace(users[0].ID, users[0].Group)
		} else {
			_ = info.RegisterFace(info.User, info.Group)
		}
		_ = info.save()
	}
	time.Sleep(time.Second * 1)
	mine.thumbIndex -= 1
	//mine.thumbPool.Submit(mine.checkPendingThumb)
	go mine.checkPendingThumb()
}

func (mine *cacheContext) firstPendingThumb() *ThumbInfo {
	if len(mine.thumbs) < 1 {
		return nil
	}
	arr := mine.thumbs[:0]
	info := mine.thumbs[0]
	for i, v := range mine.thumbs {
		if i != 0 {
			arr = append(arr, v)
		}
	}
	mine.thumbs = arr
	return info
}
