package config

type NacosConfig struct {
	Host      string `json:"host"`      //nacos ip
	Port      uint64 `json:"port"`      //端口
	Namespace string `json:"namespace"` //命名空间
	User      string `json:"user"`      //用户名
	Password  string `json:"password"`  //密码
	DataId    string `json:"dataid"`    //数据集
	Group     string `json:"group"`     //分组名称
}

type SrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type AliSmsConfig struct {
	ApiKey     string `mapstructure:"key" json:"key"`
	ApiSecrect string `mapstructure:"secrect" json:"secrect"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ServerConfig struct {
	Name         string        `mapstructure:"name" json:"name"`
	Host         string        `mapstructure:"host" json:"host"`
	Tags         []string      `mapstructure:"tags" json:"tags"`
	Port         int           `mapstructure:"port" json:"port"`
	UserSrvInfo  SrvConfig    `mapstructure:"user_srv" json:"user_srv"`
	GoodsSrvInfo  SrvConfig   `mapstructure:"goods_srv" json:"user_srv"`
	InventorySrvInfo  SrvConfig `mapstructure:"inventory_srv" json:"user_srv"`
	JWTInfo      JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo   AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo    RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo   ConsulConfig  `mapstructure:"consul" json:"consul"`
	JaegerConfig JaegerConfig  `mapstructure:"jaeger" json:"jaeger"`
	AliPayInfo   AlipayConfig  `mapstructure:"alipay" json:"alipay"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}
