package models

type Config struct {
	Http struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	}

	Service struct {
		Name string `json:"name"`
	}

	JWTConfig struct {
		SigningKey string
	}

	AliSmsConfig struct {
		ApiKey     string
		ApiSecrect string
	}

	RedisConfig struct {
		Host   string
		Port   int
		Expire int
	}

}

type Http struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type Service struct {
	Name string `json:"name"`
}

type ConfigList struct {
	Develop Config `json:"develop"`
	Test    Config `json:"test"`
	Release Config `json:"release"`
}

type Config1 struct {
	MainMysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		DataBase string `yaml:"database"`
		Charset  string `yaml:"charset"`
	}

	Redis struct {
		Host    string `yaml:"host"`
		Library string `yaml:"library"`
	}

	RabbitMQ struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
	}
}
