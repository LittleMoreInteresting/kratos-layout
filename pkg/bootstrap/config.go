package bootstrap

import (
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	// consul config
	consulKratos "github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/hashicorp/consul/api"
)

// getConfigKey 获取合法的配置名
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewConsulConfigSource 创建一个远程配置源 - Consul
func NewConsulConfigSource(configHost, configKey string) config.Source {
	consulClient, err := api.NewClient(&api.Config{
		Address: configHost,
	})
	if err != nil {
		panic(err)
	}
	consulSource, err := consulKratos.New(consulClient,
		consulKratos.WithPath(getConfigKey(configKey, true)),
	)
	if err != nil {
		panic(err)
	}

	return consulSource
}

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return file.NewSource(filePath)
}

// NewConfigProvider 创建一个配置
func NewConfigProvider(conf, consul, configKey string) config.Config {
	s := []config.Source{}
	if consul == "" {
		s = append(s, NewFileConfigSource(conf))
	} else {
		s = append(s, NewConsulConfigSource(consul, configKey))
	}
	return config.New(
		config.WithSource(
			s...,
		),
	)
}
