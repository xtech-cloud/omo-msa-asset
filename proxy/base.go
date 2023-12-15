package proxy

type LocationInfo struct {
	Left     float32 `json:"left"`
	Top      float32 `json:"top"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Rotation float32 `json:"rotation"`
}
