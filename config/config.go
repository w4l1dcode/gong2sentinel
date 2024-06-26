package config

import (
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLogLevel      = "DEBUG"
	defaultRetentionDays = 90
)

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" valid:"optional"`
	} `json:"log"`

	Microsoft struct {
		AppID          string `yaml:"app_id" env:"MS_APP_ID" valid:"minstringlength(3)"`
		SecretKey      string `yaml:"secret_key" env:"MS_SECRET_KEY" valid:"minstringlength(3)"`
		TenantID       string `yaml:"tenant_id" env:"MS_TENANT_ID" valid:"minstringlength(3)"`
		SubscriptionID string `yaml:"subscription_id" env:"MS_SUB_ID" valid:"minstringlength(3)"`

		DataCollection struct {
			Endpoint                 string `yaml:"endpoint" env:"MS_DCR_ENDPOINT" valid:"minstringlength(3)"`
			RuleID                   string `yaml:"rule_id" env:"MS_DCR_RULE" valid:"minstringlength(3)"`
			StreamNameAuditing       string `yaml:"stream_name_auditing" env:"MS_DCR_STREAM_AUDITING" valid:"minstringlength(3)"`
			StreamNameCallUserAccess string `yaml:"stream_name_user_access" env:"MS_DCR_STREAM_CALL_USER_ACCESS" valid:"minstringlength(3)"`
		} `yaml:"dcr"`

		ResourceGroup string `yaml:"resource_group" env:"MS_RSG_ID" valid:"minstringlength(3)"`
		WorkspaceName string `yaml:"workspace_name" env:"MS_WS_NAME" valid:"minstringlength(3)"`

		RetentionDays uint32 `yaml:"retention_days" env:"MS_RETENTION_DAYS" valid:"optional"`
	} `yaml:"microsoft"`

	Gong struct {
		AccessKey    string `yaml:"access_key" env:"GONG_ACCESS_KEY" valid:"minstringlength(3)"`
		AccessSecret string `yaml:"access_secret" env:"GONG_ACCESS_SECRET" valid:"minstringlength(3)"`
		LookupHours  int64  `yaml:"lookup_hours" env:"GONG_LOOKUP_HOURS" valid:"numeric"`
	} `yaml:"gong"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if c.Microsoft.RetentionDays == 0 {
		c.Microsoft.RetentionDays = defaultRetentionDays
	}

	if valid, err := validator.ValidateStruct(c); !valid || err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	if c.Gong.LookupHours <= 0 {
		return fmt.Errorf("invalid lookup hours, should be positive number: %d", c.Gong.LookupHours)
	}

	return nil
}

func (c *Config) Load(path string) error {
	if path != "" {
		configBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load configuration file at '%s': %v", path, err)
		}

		if err = yaml.Unmarshal(configBytes, c); err != nil {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("could not load environment: %v", err)
	}

	return nil
}
