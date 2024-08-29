package btsite

import (
	"fmt"
)

var SiteHelper = new(helper)

type helper struct {
}

func (sh *helper) GetConfigByCode(code string) (Config, error) {
	// 根据 code 获取站点配置
	for _, siteConfig := range globalConfig.Configs {
		if siteConfig.Id == code {
			return siteConfig, nil
		}
	}
	return Config{}, fmt.Errorf("站点配置不存在: %s", code)
}

func (sh *helper) GetDomain(site Site) (string, error) {
	if len(site.Domain) > 0 {
		return site.Domain, nil
	}
	sc, err := sh.GetConfigByCode(site.Code)
	if err != nil {
		return "", err
	}
	return sc.Domain, nil
}

func (sh *helper) GetApi(site Site) (string, error) {
	if len(site.Api) > 0 {
		return site.Api, nil
	}
	sc, err := sh.GetConfigByCode(site.Code)
	if err != nil {
		return "", err
	}
	return sc.Api, nil
}

func (sh *helper) AllSupportedSites() []Config {
	return globalConfig.Configs
}
