package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.asset/proxy/nosql"
	"time"
)

type ThumbInfo struct {
	ID       uint64 `json:"-"`
	Created  int64
	Updated  int64
	Probably float32
	Similar  float32
	Blur     float32
	UID      string `json:"uid"`
	Creator  string
	Operator string

	File  string
	Owner string
	Asset string
	Face  string
	Meta  string
	Token string
}

func CreateThumb(asset, key, owner, bs64, operator string, info *ImageFace) (*ThumbInfo, error) {
	db := new(nosql.Thumb)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetThumbNextID()
	db.Created = time.Now().Unix()
	db.Creator = operator
	db.Operator = operator
	db.Token = info.Token
	db.File = key
	db.Asset = asset
	db.Blur = info.Quality.Blur
	db.Owner = owner
	db.Probably = info.Probability
	db.Similar = 0
	db.Location = info.Location
	db.Meta = bs64
	err := nosql.CreateThumb(db)
	if err == nil {
		data := new(ThumbInfo)
		data.initInfo(db)
		return data, nil
	}
	return nil, err
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
	mine.Face = db.FaceID
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Meta = db.Meta
	mine.Token = db.Token
	mine.File = db.File
}

func (mine *ThumbInfo) UpdateBase(owner string, similar float32) error {
	err := nosql.UpdateThumbBase(mine.UID, owner, similar)
	if err == nil {
		mine.Owner = owner
		mine.Similar = similar
	}
	return err
}

func (mine *ThumbInfo) UpdateInfo(meta, operator string) error {
	err := nosql.UpdateThumbMeta(mine.UID, meta, operator)
	if err == nil {
		mine.Meta = meta
		mine.Operator = operator
	}
	return err
}
