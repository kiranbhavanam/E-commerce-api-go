package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig struct{
	Host string
	Port int
	User string
	Password string
	DatabaseName string
	SSLMode string
	MaxOpenconns int
	MaxIdleconns int
	MaxLifeTime time.Duration
}

func LoadDBConfig() *DatabaseConfig{
	port,_:=strconv.Atoi(getEnv("DB_PORT","5433"))
	maxOpen,_:=strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS","25"))
	maxIdle,_:=strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS","5"))
	maxLifetime,_:=time.ParseDuration(getEnv("DB_MAX_LIFETIME","5m"))
	
	return &DatabaseConfig{
		Host: getEnv("DB_HOST","localhost"),
		Port:port,
		User:getEnv("DB_USER","postgres"),
		Password: getEnv("DB_PASSWORD","password"),
		DatabaseName: getEnv("DB_NAME","e-commerce"),
		SSLMode: getEnv("SSL_MODE","disable"),
		MaxOpenconns: maxOpen,
		MaxIdleconns: maxIdle,
		MaxLifeTime: maxLifetime,
	}
}
func (cfg *DatabaseConfig)ConnectionString()string{
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
	cfg.Host,cfg.Port,cfg.User,cfg.Password,cfg.DatabaseName,cfg.SSLMode)
	
}
func getEnv(s string,d string) string{
	if value:=os.Getenv(s);value!=""{
		return value
	}
	return d
}