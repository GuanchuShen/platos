package config

import (
	"time"

	"github.com/spf13/viper"
)

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// GetEndpointsForDiscovery 获取服务发现的地址
func GetEndpointsForDiscovery() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

// GetTimeoutForDiscovery 获取连接服务发现集群的超时时间，单位：秒
func GetTimeoutForDiscovery() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

func GetServicePathForIPConf() string {
	return viper.GetString("ip_conf.service_path")
}

func GetCachedRedisEndpointList() []string {
	return viper.GetStringSlice("cache.redis.endpoints")
}

// IsDebug 判断是否为 debug 环境
func IsDebug() bool {
	env := viper.GetString("global.env")
	return env == "debug"
}
