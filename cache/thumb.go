package cache

import (
	"omo.msa.asset/proxy/nosql"
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

	Owner string
	Asset string
	Face  string
	URL   string
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
	mine.URL = db.URL
	mine.Creator = db.Creator
	mine.Operator = db.Operator
}

func (mine *ThumbInfo) UpdateBase(owner string, similar float32) error {
	err := nosql.UpdateThumbBase(mine.UID, owner, similar)
	if err == nil {
		mine.Owner = owner
		mine.Similar = similar
	}
	return err
}
