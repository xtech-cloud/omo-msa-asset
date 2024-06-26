package proxy

type LocationInfo struct {
	Left     float32 `json:"left" bson:"left"`
	Top      float32 `json:"top" bson:"top"`
	Width    uint32  `json:"width" bson:"width"`
	Height   uint32  `json:"height" bson:"height"`
	Rotation float32 `json:"rotation" bson:"rotation"`
}

type PairInfo struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
	Count uint32 `json:"count" bson:"count"`
}
