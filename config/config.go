package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB DB
}

type DB struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}

func GetAppConfig(filename, path string) *Config {
	conf := loadConfig(filename, path)
	return conf
}

func InitRouters() *chi.Mux {
	r := chi.NewRouter()
	// setup cors here ...
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)
	return r
}

func loadConfig(filename, path string) *Config {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("fatal: error reading config file", err.Error())
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("fatal: error reading config variable", err.Error())
	}
	return &conf
}

// GetDBConn ...
func GetDBConn(log *zap.SugaredLogger, conf DB) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user='%s' password=%s dbname=%s port=%s sslmode=%s", conf.HOST, conf.USER, conf.PASSWORD, conf.DBNAME, conf.PORT, "disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// GetConsoleLogger ...
func GetConsoleLogger() *zap.SugaredLogger {
	encoder := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
