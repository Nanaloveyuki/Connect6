package config

import (
	"errors"
	"strings"

	"connect6/backend/internal/constant"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	GitHub GitHubConfig `mapstructure:"github"`
	HTTP   HTTPConfig   `mapstructure:"http"`
	Graph  GraphConfig  `mapstructure:"graph"`
}

type AppConfig struct {
	Env                 string `mapstructure:"env"`
	Port                string `mapstructure:"port"`
	ReadTimeoutSeconds  int    `mapstructure:"readTimeoutSeconds"`
	WriteTimeoutSeconds int    `mapstructure:"writeTimeoutSeconds"`
}

type GitHubConfig struct {
	BaseURL               string `mapstructure:"baseURL"`
	Token                 string `mapstructure:"token"`
	RequestTimeoutSeconds int    `mapstructure:"requestTimeoutSeconds"`
	FollowingLimit        int    `mapstructure:"followingLimit"`
	OrgLimit              int    `mapstructure:"orgLimit"`
}

type HTTPConfig struct {
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
}

type GraphConfig struct {
	MaxDepth int `mapstructure:"maxDepth"`
}

func Load(configDir string) (Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)
	v.SetEnvPrefix("CONNECT6")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFound) {
			return Config{}, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.env", constant.DefaultEnvironment)
	v.SetDefault("app.port", constant.DefaultHTTPPort)
	v.SetDefault("app.readTimeoutSeconds", constant.DefaultHTTPReadTimeoutSeconds)
	v.SetDefault("app.writeTimeoutSeconds", constant.DefaultHTTPWriteTimeoutSeconds)
	v.SetDefault("github.baseURL", constant.DefaultGitHubBaseURL)
	v.SetDefault("github.requestTimeoutSeconds", constant.DefaultGitHubTimeoutSeconds)
	v.SetDefault("github.followingLimit", constant.DefaultGitHubFollowingLimit)
	v.SetDefault("github.orgLimit", constant.DefaultGitHubOrgLimit)
	v.SetDefault("http.allowedOrigins", constant.DefaultAllowedOrigins)
	v.SetDefault("graph.maxDepth", constant.DefaultGraphMaxDepth)
}
