package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL        string
	LogLevel           string
	HTTPPort           string
	JWTSecret          string
	AccessTTL          time.Duration
	RefreshTTL         time.Duration
	RedisAddress       string
	RedisPassword      string
	RedisDB            int
	StorageEndpoint    string
	StorageAccessKey   string
	StorageSecretKey   string
	StorageBucketName  string
	StorageUseSSL      bool
	SiteBaseURL        string
	ExternalAPIURL     string
	ExternalAPITimeout time.Duration
	ExternalAPIRetry   int
	ExternalAPIRate    float64
}

func InitConfig(path string) Config {
	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}
	AccessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TTL"))
	if err != nil {
		panic(err)
	}
	RefreshTTl, err := time.ParseDuration(os.Getenv("REFRESH_TTL"))
	if err != nil {
		panic(err)
	}
	redisdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}
	useSSL, err := strconv.ParseBool(os.Getenv("STORAGE_USE_SSL"))
	if err != nil {
		useSSL = false // Default to false for local MinIO
	}
	externalTimeout, err := time.ParseDuration(os.Getenv("EXTERNAL_API_TIMEOUT"))
	if err != nil || externalTimeout <= 0 {
		externalTimeout = 10 * time.Second
	}
	externalRetry, err := strconv.Atoi(os.Getenv("EXTERNAL_API_RETRY"))
	if err != nil || externalRetry < 0 {
		externalRetry = 2
	}
	externalRate, err := strconv.ParseFloat(os.Getenv("EXTERNAL_API_RATE"), 64)
	if err != nil || externalRate <= 0 {
		externalRate = 5.0
	}
	return Config{
		PostgresURL:        os.Getenv("POSTGRES_URL"),
		LogLevel:           os.Getenv("LOG_LEVEL"),
		HTTPPort:           os.Getenv("HTTP_PORT"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AccessTTL:          AccessTTL,
		RefreshTTL:         RefreshTTl,
		RedisAddress:       os.Getenv("REDIS_ADDRESS"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            redisdb,
		StorageEndpoint:    os.Getenv("STORAGE_ENDPOINT"),
		StorageAccessKey:   os.Getenv("STORAGE_ACCESS_KEY"),
		StorageSecretKey:   os.Getenv("STORAGE_SECRET_KEY"),
		StorageBucketName:  os.Getenv("STORAGE_BUCKET_NAME"),
		StorageUseSSL:      useSSL,
		SiteBaseURL:        os.Getenv("SITE_BASE_URL"),
		ExternalAPIURL:     os.Getenv("EXTERNAL_API_URL"),
		ExternalAPITimeout: externalTimeout,
		ExternalAPIRetry:   externalRetry,
		ExternalAPIRate:    externalRate,
	}
}
