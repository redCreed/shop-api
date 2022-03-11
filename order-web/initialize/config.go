package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"mxshop-api/order-web/global"
)

//使用nacos配置
func InitConfig() {
	debug := true
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("goods-web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("读取配置出现错误:", err)
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		fmt.Println("反序列化配置出现错误:", err)
		panic(err)
	}
	fmt.Sprintf("配置信息:%+v\n", global.NacosConfig)

	//从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	fmt.Println("content:", content)
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(&global.ServerConfig)
}

//初始化配置
//func InitConfig(mode string) {
//	v := viper.New()
//	v.AddConfigPath("config")
//	v.SetConfigName("config")
//	v.SetConfigType("yaml")
//	if err := v.ReadInConfig(); err != nil {
//		fmt.Println(runtime.Caller(0))
//	}
//
//	cfg := new(models.Config)
//	if err := v.UnmarshalKey(mode, cfg); err != nil {
//		panic(err)
//	}
//
//	//var cfg models.ConfigList
//	//if err := v.Unmarshal(&cfg); err != nil {
//	//	panic(err)
//	//}
//
//	global.Config =cfg
//	fmt.Println("cfg:", cfg)
//}
