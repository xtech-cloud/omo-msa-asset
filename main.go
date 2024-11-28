package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	proto "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"io"
	"omo.msa.asset/cache"
	"omo.msa.asset/config"
	"omo.msa.asset/grpc"
	"os"
	"path/filepath"
	"time"
)

var (
	BuildVersion string
	BuildTime    string
	CommitID     string
)

func main() {
	config.Setup()
	err := cache.InitData()
	if err != nil {
		panic(err)
	}
	// New Service
	service := micro.NewService(
		micro.Name("omo.msa.asset"),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)
	// Initialise service
	service.Init()
	// Register Handler
	_ = proto.RegisterAssetServiceHandler(service.Server(), new(grpc.AssetService))
	_ = proto.RegisterThumbServiceHandler(service.Server(), new(grpc.ThumbService))
	_ = proto.RegisterFolderServiceHandler(service.Server(), new(grpc.FolderService))
	_ = proto.RegisterLabelServiceHandler(service.Server(), new(grpc.LabelService))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service

	go delayRun()

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func delayRun() {
	time.Sleep(5 * time.Second)
	cache.Context().CheckThumbs()
	//cache.PublishSystemAssets()
	//cache.TestDetectFaces()
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()
	io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil))
}
