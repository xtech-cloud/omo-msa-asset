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
		"ip": "192.168.1.10",
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
		"source": "http://rdpup1.suii.cn",
		"access": "4TDqfvaNHKxzx4nFz0YglS_jHlKXECCSSWb1vUr5",
		"secret": "pZ8AnJE5IYgNRUFEB132ohIToJdRe5uxm4ZLLljp",
		"bucket": "tec-test"
	},
	"examine":{
		"type":"baidu",
		"app":"40017362",
		"access":"tNEOCTXw0j6wQCGgIEOj54f6",
		"secret":"1w5pDR2bqDf0vgmHUCYqd8YPoACw6d25"
	},
	"detection":{
		"type":"baidu",
		"app":"24196251",
		"address":"https://aip.baidubce.com/rest/2.0/face/v3/detect",
		"access":"Z5NUW1i1gzu4VBgOcMUDg8IK",
		"secret":"X3EaQselwNLCvAcWH6FdRvUgRkTZo0T2",
		"user":{
			"add":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add",
			"delete":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/delete",
			"update":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/update",
			"list":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/group/getusers",
			"get":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/get"
		},
		"face":{
			"add":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add",
			"delete":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/face/delete",
			"list":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/face/getlist"
		},
		"group":{
			"add":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/group/add",
			"delete":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/group/delete",
			"list":"https://aip.baidubce.com/rest/2.0/face/v3/faceset/group/getlist"
		}
	},
	"basic": {
		"synonym": 6,
		"tag": 6
	}
}
`
