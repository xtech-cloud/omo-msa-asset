package config

type ServiceConfig struct {
	TTL      int64  `json:"ttl"`
	Interval int64  `json:"interval"`
	Address  string `json:"address"`
}

type LoggerConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
	Std   bool   `json:"std"`
}

type DBConfig struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

type BasicConfig struct {
	SynonymMax int `ini:"synonym"`
	TagMax     int `ini:"tag"`
}

type StorageConfig struct {
	Type string `json:"type"`
	/**
	最大尺寸大小
	*/
	Limit int32 `json:"limit"`
	/**
	token最大时效（秒）
	*/
	Expire int32 `json:"expire"`
	//CDN 域名
	Domain string `json:"domain"`
	//源域名
	Source string `json:"source"`
	// 公有库或者私有库
	ACM int `json:"acm"`
	// 私有地址的有效时间（秒）
	Period    int64  `json:"period"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access"`
	SecretKey string `json:"secret"`
}

type ExamineConfig struct {
	Type      string `json:"type"`
	App       string `json:"app"`
	AccessKey string `json:"access"`
	SecretKey string `json:"secret"`
}

type APIConfig struct {
	Type      string     `json:"type"`
	App       string     `json:"app"`
	Address   string     `json:"address"`
	AccessKey string     `json:"access"`
	SecretKey string     `json:"secret"`
	Face      *APIOption `json:"face"`
	User      *APIOption `json:"user"`
	Group     *APIOption `json:"group"`
}

type APIOption struct {
	Add    string `json:"add"`
	Delete string `json:"delete"`
	List   string `json:"list"`
	Get    string `json:"get"`
	Update string `json:"update"`
}

type SchemaConfig struct {
	Service   ServiceConfig `json:"service"`
	Logger    LoggerConfig  `json:"logger"`
	Database  DBConfig      `json:"database"`
	Basic     BasicConfig   `json:"basic"`
	Storage   StorageConfig `json:"storage"`
	Examine   ExamineConfig `json:"examine"`
	Detection APIConfig     `json:"detection"`
}
