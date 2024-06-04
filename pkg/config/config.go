package config

import (
	echoserver "example.com/fxdemo/pkg/http/echo/server"
	"example.com/fxdemo/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	defaultHTTPPort               = ":8000"
	defaultHTTPRWTimeout          = 60 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10000
	defaultLimiterBurst           = 20000
	defaultLimiterTTL             = 1 * time.Minute
	defaultVerificationCodeLength = 8

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment  string
		Mongo        MongoConfig
		HTTP         HTTPConfig
		Echo         echoserver.EchoConfig
		Limiter      LimiterConfig
		CacheTTL     time.Duration `mapstructure:"ttl"`
		ServiceUrl   string
		LoggerConfig logger.LoggerConfig
	}
	EchoConfig struct {
		Port                string   `mapstructure:"port" validate:"required"`
		Development         bool     `mapstructure:"development"`
		BasePath            string   `mapstructure:"basePath" validate:"required"`
		DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
		IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
		Timeout             int      `mapstructure:"timeout"`
		Host                string   `mapstructure:"host"`
	}
	MongoConfig struct {
		MongoConnection string `json:"mongo_connection"`
		MongoDbName     string `json:"mongo_db_name"`
	}
	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init() (*Config, *MongoConfig,
	*HTTPConfig,
	*echoserver.EchoConfig,
	*LimiterConfig,
	logger.LoggerConfig, error) {
	configsDir := "pkg/config/configs"
	populateDefaults()
	err := godotenv.Load()
	if err != nil {
		logrus.Info("Error  from load env. this mean the application load on the cloud not from a file.")
	}
	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, nil, nil, nil, nil, logger.LoggerConfig{}, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, nil, nil, nil, nil, logger.LoggerConfig{}, err
	}

	setFromEnv(&cfg)

	return &cfg, &cfg.Mongo, &cfg.HTTP, &cfg.Echo, &cfg.Limiter, cfg.LoggerConfig, nil
}

func unmarshal(cfg *Config) error {

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) {
	// TODO use envconfig https://github.com/kelseyhightower/envconfig
	cfg.Mongo.MongoConnection = os.Getenv("MONGO_CONNECTION")
	cfg.Mongo.MongoDbName = os.Getenv("MONGO_DB_NAME")

	cfg.ServiceUrl = os.Getenv("SERVICE_URL")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.Environment = os.Getenv("APP_ENV")

	cfg.Echo.Port = defaultHTTPPort
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}
