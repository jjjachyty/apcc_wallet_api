package dimSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/utils"
)

//All 查询所有的配置
func AllConfig() (configs []dimMod.DimConfig, err error) {
	configs = make([]dimMod.DimConfig, 0)
	err = models.GetBeans(&configs, dimMod.DimConfig{})
	return
}

func AddOrUpdateConfig(config dimMod.DimConfig) (err error) {
	if err = utils.HMSet(config.Key, map[string]interface{}{config.SubKey: config.Value}); err == nil {
		err = models.Create(config)
	}
	return
}

func InitDimConfig() {
	if configs, err := AllConfig(); err == nil {
		for _, config := range configs {
			if err := utils.HSet(config.Key, config.SubKey, config.Value); err != nil {
				utils.SysLog.Panic("Redis缓存设置错误")
			}
		}
	} else {
		utils.SysLog.Panic("加载配置表dim_config出错")
	}
}
