package cache

import (
	"omo.msa.asset/proxy/nosql"
	"time"
)

type ThumbInfo struct {
	ID       uint64 `json:"-"`
	Probably float32
	UID      string `json:"uid"`
	Creator  string
	Operator string

	Owner      string
	Asset      string
	Face       string
	URL        string
	CreateTime time.Time
	UpdateTime time.Time
}

func (mine *ThumbInfo) initInfo(db *nosql.Thumb) {
	mine.ID = db.ID
	mine.UID = db.UID.Hex()
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Owner = db.Owner
	mine.Asset = db.Asset
	mine.Probably = db.Probably
	mine.Face = db.FaceID
	mine.URL = db.URL
	mine.Creator = db.Creator
	mine.Operator = db.Operator
}

func (mine *ThumbInfo)UpdateBase(owner string, probably float32) error {
	err := nosql.UpdateThumbBase(mine.UID, owner, probably)
	if err == nil {
		mine.Owner = owner
		mine.Probably = probably
	}
	return err
}