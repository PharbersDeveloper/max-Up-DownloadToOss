package BmFactory

import (
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmDataStorage"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmHandler"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmResource"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmMiddleware"
)

type BmTable struct{}

var BLACKMIRROR_MODEL_FACTORY = map[string]interface{}{
	"BmFiles":          BmModel.Files{},
}

var BLACKMIRROR_RESOURCE_FACTORY = map[string]interface{}{
	"BmFilesResource":            BmResource.BmFilesResource{},
}

var BLACKMIRROR_STORAGE_FACTORY = map[string]interface{}{
	"BmFilesStorage":            BmDataStorage.BmFilesStorage{},
}

var BLACKMIRROR_MIDDLEWARE_FACTORY = map[string]interface{}{
	"BmCheckTokenMiddleware": BmMiddleware.BmCheckTokenMiddleware{},
}

var BLACKMIRROR_DAEMON_FACTORY = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
}

var BLACKMIRROR_FUNCTION_FACTORY = map[string]interface{}{
	"BmUploadToOssHandler":     	   BmHandler.UploadToOssHandler{},
	"BmAccountHandler":            	   BmHandler.AccountHandler{},
	"BmUserAgentHandler":       	   BmHandler.UserAgentHandler{},
	"BmGenerateAccessTokenHandler":    BmHandler.BmGenerateAccessTokenHandler{},
	"BmRefreshAccessTokenHandler":     BmHandler.RefreshAccessTokenHandler{},
}


func (t BmTable) GetModelByName(name string) interface{} {
	return BLACKMIRROR_MODEL_FACTORY[name]
}

func (t BmTable) GetResourceByName(name string) interface{} {
	return BLACKMIRROR_RESOURCE_FACTORY[name]
}

func (t BmTable) GetStorageByName(name string) interface{} {
	return BLACKMIRROR_STORAGE_FACTORY[name]
}

func (t BmTable) GetDaemonByName(name string) interface{} {
	return BLACKMIRROR_DAEMON_FACTORY[name]
}

func (t BmTable) GetFunctionByName(name string) interface{} {
	return BLACKMIRROR_FUNCTION_FACTORY[name]
}

func (t BmTable) GetMiddlewareByName(name string) interface{} {
	return BLACKMIRROR_MIDDLEWARE_FACTORY[name]
}