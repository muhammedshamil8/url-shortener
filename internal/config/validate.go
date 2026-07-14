package config

import (
	"fmt"
	"strings"
)

func (c *Config) Validate() error {
	missing := []string{}
	if c.Server.Port == "" {
		missing = append(missing, "APP_PORT")
	}
	if c.Server.BaseURL == "" {
		missing = append(missing, "BASE_URL")
	}
	if c.DB.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DB.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DB.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.DB.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if c.DB.SSLMode == "" {
		missing = append(missing, "DB_SSLMODE")
	}
	if c.Server.AllowedOrigins == nil {
		missing = append(missing, "ALLOWED_ORIGINS")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}
	return nil
}
