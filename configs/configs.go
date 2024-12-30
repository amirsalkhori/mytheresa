package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App        App    `yaml:"app" mapstructure:"app"`
	Mysql      Mysql  `yaml:"mysql" mapstructure:"mysql"`
	Redis      Redis  `mapstructure:"redis"`
	HashIDSalt string `mapstructure:"hash_id_salt"`
}

type App struct {
	Name string `yaml:"name" mapstructure:"name"`
}

type Mysql struct {
	Host string `yaml:"host" mapstructure:"host"`
	Port int    `yaml:"port" mapstructure:"port"`
	User string `yaml:"username" mapstructure:"username"`
	Name string `yaml:"database" mapstructure:"database"`
	Pass string `yaml:"password" mapstructure:"password"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

func GetConfig() Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	bindEnv(v)
	_ = v.ReadInConfig()

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		// panic(err)
		fmt.Println("Log error unmarshl")
	}

	c.App = App{
		Name: v.GetString("app.name"),
	}

	return c
}

func bindEnv(v *viper.Viper) {
	envBindings := map[string]string{
		"app.name":       "APP_NAME",
		"mysql.host":     "MYTHERESA_MYSQL_HOST",
		"mysql.port":     "MYTHERESA_MYSQL_PORT",
		"mysql.username": "MYTHERESA_MYSQL_USER",
		"mysql.password": "MYTHERESA_MYSQL_PASSWORD",
		"mysql.database": "MYTHERESA_MYSQL_DB",
		"redis.host":     "MYTHERESA_REDIS_HOST",
		"redis.port":     "MYTHERESA_REDIS_PORT",
		"redis.password": "MYTHERESA_REDIS_PASSWORD",
		"hash_id_salt":   "HASH_ID_SAlT",
	}

	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			fmt.Println("err bind", err)
		}
	}
}
