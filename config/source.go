package config

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/config/source/memory"
	logrusPlugin "github.com/micro/go-plugins/logger/logrus/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/encoder/yaml"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/config/source/file"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/config/source/consul/v2"
)

type EnvConfig struct {
	Source  string   `json:source`
	Prefix  string   `json:prefix`
	Key     string   `json:key`
	Address []string `json:address`
}

var envConfig EnvConfig

var Schema SchemaConfig

func setupEnvironment() {
	//registry plugin
	registryPlugin := os.Getenv("MSA_REGISTRY_PLUGIN")
	if "" == registryPlugin {
		registryPlugin = "consul"
	}
	os.Setenv("MICRO_REGISTRY", registryPlugin)

	//registry address
	registryAddress := os.Getenv("MSA_REGISTRY_ADDRESS")
	if "" == registryAddress {
		registryAddress = "127.0.0.1:8500"
	}
	_ = os.Setenv("MICRO_REGISTRY_ADDRESS", registryAddress)

	//config
	envConfigDefine := os.Getenv("MSA_CONFIG_DEFINE")

	if "" == envConfigDefine {
		logger.Warn("MSA_CONFIG_DEFINE is empty")
		return
	}

	logger.Infof("MSA_CONFIG_DEFINE is %v", envConfigDefine)
	err := json.Unmarshal([]byte(envConfigDefine), &envConfig)
	if err != nil {
		logger.Error(err)
	}
}

func mergeFile(_config config.Config) {
	filepath := envConfig.Prefix + envConfig.Key
	fileSource := file.NewSource(
		file.WithPath(filepath),
	)
	err := _config.Load(fileSource)
	if nil != err {
		logger.Errorf("load config %v failed: %v", filepath, err)
	} else {
		logger.Infof("load config %v success", filepath)
		_config.Scan(&Schema)
	}
}

func mergeConsul(_config config.Config) {
	consulKey := envConfig.Prefix + envConfig.Key
	consulSource := consul.NewSource(
		consul.WithPrefix(envConfig.Prefix),
		consul.StripPrefix(true),
		source.WithEncoder(yaml.NewEncoder()),
	)
Loop:
	for {
		select {
		case <-time.After(time.Second * time.Duration(2)):
			for _, addr := range envConfig.Address {
				consul.WithAddress(addr)
				err := _config.Load(consulSource)
				if nil == err {
					logger.Infof("load config %v from %v success", consulKey, addr)
					break Loop
				} else {
					logger.Errorf("load config %v from %v failed: %v", consulKey, addr, err)
				}
			}
		}
	}
	_config.Get(envConfig.Key).Scan(&Schema)
}

func mergeEtcd(_config config.Config) {
	etcdKey := envConfig.Prefix + envConfig.Key
	etcdSource := etcd.NewSource(
		etcd.WithPrefix(envConfig.Prefix),
		etcd.StripPrefix(true),
		source.WithEncoder(yaml.NewEncoder()),
	)
Loop:
	for {
		select {
		case <-time.After(time.Second * time.Duration(2)):
			for _, addr := range envConfig.Address {
				etcd.WithAddress(addr)
				err := _config.Load(etcdSource)
				if nil == err {
					logger.Infof("load config %v from %v success", etcdKey, addr)
					break Loop
				} else {
					logger.Errorf("load config %v from %v failed: %v", etcdKey, addr, err)
				}
			}
		}
	}
	_config.Get(envConfig.Key).Scan(&Schema)
}

func Setup() {
	mode := os.Getenv("MSA_MODE")
	if "" == mode {
		mode = "debug"
	}
	setupEnvironment()
	conf, err := config.NewConfig()
	if nil != err {
		panic(err)
	}

	// load default config
	logger.Infof("default config is: \n\r%v", defaultJson)
	memorySource := memory.NewSource(
		memory.WithJSON([]byte(defaultJson)),
	)
	conf.Load(memorySource)
	err1 := conf.Scan(&Schema)
	if err1 != nil {
		panic(err1)
		return
	}

	// merge others
	if "file" == envConfig.Source {
		mergeFile(conf)
	} else if "consul" == envConfig.Source {
		if mode != "debug" {
			mergeConsul(conf)
		}
	} else if "etcd" == envConfig.Source {
		if mode != "debug" {
			mergeEtcd(conf)
		}
	}

	ycd, err := json.Marshal(&Schema)
	if nil != err {
		logger.Error(err)
	} else {
		logger.Infof("current config is: \n\r%v", string(ycd))
	}

	// initialize logger
	initLogger(mode)
}

func getLoggerOut() io.Writer {
	path := Schema.Logger.File
	logger.Info("logger path = " + path)
	log :=&lumberjack.Logger{
		LocalTime:  true,
		Filename:   path,
		MaxSize:    20, // megabytes
		MaxBackups: 20,
		MaxAge:     0,    //days
		Compress:   false, // disabled by default
	}
	if Schema.Logger.Std {
		writers := []io.Writer{
			log,
			os.Stdout,
		}
		return io.MultiWriter(writers...)
	}else{
		writers := []io.Writer{
			log,
		}
		return io.MultiWriter(writers...)
	}
}

func initLogger(mode string)  {
	out := getLoggerOut()
	if "debug" == mode {
		logger.DefaultLogger = logrusPlugin.NewLogger(
			logger.WithOutput(out),
			logger.WithLevel(logger.TraceLevel),
			logrusPlugin.WithTextTextFormatter(new(logrus.TextFormatter)),
		)
		logger.Info("-------------------------------------------------------------")
		logger.Info("- Micro Service Agent -> Setup")
		logger.Info("-------------------------------------------------------------")
		logger.Warn("Running in \"debug\" mode. Switch to \"release\" mode in production.")
		logger.Warn("- using env:	export MSA_MODE=release")
	} else {
		logger.DefaultLogger = logrusPlugin.NewLogger(
			logger.WithOutput(out),
			logger.WithLevel(logger.TraceLevel),
			logrusPlugin.WithJSONFormatter(new(logrus.JSONFormatter)),
		)
		logger.Info("-------------------------------------------------------------")
		logger.Info("- Micro Service Agent -> Setup")
		logger.Info("-------------------------------------------------------------")
	}


	level, err := logger.GetLevel(Schema.Logger.Level)
	if nil != err {
		logger.Warnf("the level %v is invalid, just use info level", Schema.Logger.Level)
		level = logger.InfoLevel
	}

	if "debug" == mode {
		logger.Warn("Using \"MSA_DEBUG_LOG_LEVEL\" to switch log's level in \"debug\" mode.")
		logger.Warn("- using env:	export MSA_DEBUG_LOG_LEVEL=debug")
		debugLogLevel := os.Getenv("MSA_DEBUG_LOG_LEVEL")
		if "" == debugLogLevel {
			debugLogLevel = "trace"
		}
		level, _ = logger.GetLevel(debugLogLevel)
	}
	logger.Infof("level is %v now", level)
	logger.Init(
		logger.WithLevel(level),
	)
}
