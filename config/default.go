package config

const defaultJson string = `{
	"service": {
		"address": ":7076",
		"ttl": 15,
		"interval": 10
	},
	"logger": {
		"level": "info",
		"file": "logs/server.log",
		"std": false
	},
	"database": {
		"name": "rgsCloud",
		"ip": "127.0.0.1",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	},
	"storage": {
		"type": "qiniu",
		"limit": 500,
		"expire": 360000,
		"acm": 0,
		"period": 600,
		"domain": "http://testdown.suii.cn",
		"access": "4TDqfvaNHKxzx4nFz0YglS_jHlKXECCSSWb1vUr5",
		"secret": "pZ8AnJE5IYgNRUFEB132ohIToJdRe5uxm4ZLLljp",
		"bucket": "tec-test"
	},
	"basic": {
		"synonym": 6,
		"tag": 6
	}
}
`
