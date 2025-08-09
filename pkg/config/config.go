package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	echoserver "kaptan/pkg/http/echo/server"
	"os"
	"time"
)

const (
	defaultHTTPPort               = ":8005"
	defaultHTTPRWTimeout          = 60 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10000
	defaultLimiterBurst           = 20000
	defaultLimiterTTL             = 1 * time.Minute
	defaultVerificationCodeLength = 8
)

type (
	Config struct {
		Environment    string
		Mongo          MongoConfig
		RedisConfig    RedisConfig
		AwsConfig      AwsConfig
		HTTP           HTTPConfig
		Echo           echoserver.EchoConfig
		Limiter        LimiterConfig
		CacheTTL       time.Duration `mapstructure:"ttl"`
		ServiceUrl     string
		JWTConfig      JWTConfig
		FirebaseConfig FirebaseConfig
		Mysql          MysqlConfig
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

	MysqlConfig struct {
		HOST     string
		PORT     string
		USERNAME string
		PASSWORD string
		DATABASE string
	}

	FirebaseConfig struct {
		DatabaseURL string `json:"database_url"`
		FcmFilePath string `json:"fcm_file_path"`
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

	AwsConfig struct {
		AccessKey  string
		SecretKey  string
		Region     string
		EndPoint   string
		BucketName string
	}

	JWTConfig struct {
		AdminSigningKey    string
		AdminExpires       time.Duration `mapstructure:"admin_expires"`
		PortalSigningKey   string
		KitchenSigningKey  string
		PortalExpires      time.Duration `mapstructure:"portal_expires"`
		UserSigningKey     string
		UserExpires        time.Duration `mapstructure:"user_expires"`
		UserTempSigningKey string
		UserTempExpires    time.Duration `mapstructure:"user_temp_expires"`
	}

	RedisConfig struct {
		RedisHost           string
		RedisPort           string
		RedisDbUserUsername string
		RedisDbUserPassword string
		RedisDbKey          string
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init() (*Config, error) {
	populateDefaults()
	err := godotenv.Load()
	if err != nil {
		logrus.Info("Error  from load env. this mean the application load on the cloud not from a file.")
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)

	return &cfg, nil
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

	if err := viper.UnmarshalKey("jwt", &cfg.JWTConfig); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) {
	cfg.Mongo.MongoConnection = os.Getenv("MONGO_CONNECTION")
	cfg.Mongo.MongoDbName = os.Getenv("MONGO_DB_NAME")

	cfg.RedisConfig.RedisHost = os.Getenv("REDIS_HOST")
	cfg.RedisConfig.RedisPort = os.Getenv("REDIS_PORT")
	cfg.RedisConfig.RedisDbUserUsername = os.Getenv("REDIS_USERNAME")
	cfg.RedisConfig.RedisDbUserPassword = os.Getenv("REDIS_PASSWORD")
	cfg.RedisConfig.RedisDbKey = os.Getenv("REDIS_DB_KEY")

	cfg.ServiceUrl = os.Getenv("SERVICE_URL")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.Environment = os.Getenv("APP_ENV")

	cfg.AwsConfig.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	cfg.AwsConfig.SecretKey = os.Getenv("AWS_SECRET_ID")
	cfg.AwsConfig.Region = os.Getenv("AWS_REGION")
	cfg.AwsConfig.BucketName = os.Getenv("AWS_BUCKET_NAME")
	cfg.AwsConfig.EndPoint = os.Getenv("AWS_END_POINT")

	cfg.JWTConfig.AdminSigningKey = os.Getenv("JWT_SECRET_ADMIN")
	cfg.JWTConfig.PortalSigningKey = os.Getenv("JWT_SECRET_PORTAL")
	cfg.JWTConfig.PortalExpires = time.Hour * 24  //24 hours
	cfg.JWTConfig.UserExpires = time.Hour * 24    //24 hours
	cfg.JWTConfig.UserTempExpires = time.Hour * 1 //1 hours
	cfg.JWTConfig.KitchenSigningKey = os.Getenv("JWT_SECRET_KITCHEN")
	cfg.JWTConfig.UserSigningKey = os.Getenv("JWT_SECRET_USER")
	cfg.JWTConfig.UserTempSigningKey = os.Getenv("JWT_SECRET_USER_TEMP")

	cfg.FirebaseConfig.DatabaseURL = os.Getenv("REAlTIME_DATABASE_URL")
	cfg.FirebaseConfig.FcmFilePath = os.Getenv("FCM_FILE_PATH")

	cfg.Mysql.HOST = os.Getenv("DB_HOST")
	cfg.Mysql.PORT = os.Getenv("DB_PORT")
	cfg.Mysql.USERNAME = os.Getenv("DB_USERNAME")
	cfg.Mysql.PASSWORD = os.Getenv("DB_PASSWORD")
	cfg.Mysql.DATABASE = os.Getenv("DB_DATABASE")

	var port = defaultHTTPPort
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	cfg.Echo.Port = port
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
