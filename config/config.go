package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	AppURL             string
	GomniauthKey       string
	GoogleClientID     string
	GoogleSecretValue  string
	DbName             string
	SQLDriver          string
	OpenWeatherApiKey  string
	AwsAccessKey       string
	AwsSecretAccessKey string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		AppURL:             cfg.Section("app").Key("URL").String(),
		GomniauthKey:       cfg.Section("gomniauth").Key("security_key").String(),
		GoogleClientID:     cfg.Section("google").Key("clientID").String(),
		GoogleSecretValue:  cfg.Section("google").Key("secret_value").String(),
		DbName:             cfg.Section("db").Key("name").String(),
		SQLDriver:          cfg.Section("db").Key("driver").String(),
		OpenWeatherApiKey:  cfg.Section("openWeather").Key("ApiKey").String(),
		AwsAccessKey:       cfg.Section("aws").Key("ACCESS_KEY").String(),
		AwsSecretAccessKey: cfg.Section("aws").Key("SECRET_ACCESS_KEY").String(),
	}
}
