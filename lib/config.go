package lib

import (
	"os"

	"github.com/pocketbase/pocketbase/core"
	yaml "gopkg.in/yaml.v2"
)

type OAuthClient struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Config struct {
	AppName       string      `yaml:"app_name"`
	AppUrl        string      `yaml:"app_url"`
	GoogleClient  OAuthClient `yaml:"google_client"`
	BackupCron    string      `yaml:"backup_cron"`
	BackupMaxKeep int         `yaml:"backup_max_keep"`
}

func NewConfigFromFile(file string) (*Config, error) {
	c := Config{}
	configData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configData, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Config) InitSettings(app core.App) {
	app.Settings().Meta.AppName = c.AppName
	app.Settings().Meta.AppUrl = c.AppUrl

	app.Settings().GoogleAuth.Enabled = true
	app.Settings().GoogleAuth.ClientId = c.GoogleClient.ClientId
	app.Settings().GoogleAuth.ClientSecret = c.GoogleClient.ClientSecret
	app.Settings().Backups.Cron = c.BackupCron
	app.Settings().Backups.CronMaxKeep = c.BackupMaxKeep
}
