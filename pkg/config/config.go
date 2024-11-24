package config

import (
	"os"
	"regexp"
	"strings"
	"time"

	"log/slog"

	"github.com/spf13/viper"
)

type Config interface {
	// GetPropertyString function to get a string property
	GetPropertyString(name string) string

	// GetPropertyBool function to get a boolean property
	GetPropertyBool(name string) bool

	// GetPropertyTime function to get a time property
	GetPropertyTime(name string) time.Time

	// GetPropertyStringSlice function to get a string slice property
	GetPropertyStringSlice(name string) []string

	// GetPropertyMapString function to get a map of strings property
	GetPropertyMapString(name string) map[string]string

	// AddConfigFile add a new config file yaml
	AddConfigFile(name string)
}

// Config struct to hold the configuration data and methods
type config struct {
	v *viper.Viper
}

// NewConfig function to initialize the configuration
func NewConfig(filename string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return &config{v: v}, nil
}

// NewConfig function to initialize the configuration
func NewConfigFromString(yamlContent string) (Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(yamlContent))
	if err != nil {
		return nil, err
	}
	return &config{v: v}, nil
}

// resolvePlaceholders function to resolve environment variables and references
func (c *config) resolvePlaceholders(value string) (string, error) {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	for {
		matches := re.FindAllStringSubmatch(value, -1)
		if len(matches) == 0 {
			break
		}

		for _, match := range matches {
			fullMatch := match[0]
			key := match[1]
			parts := strings.SplitN(key, ":", 2)
			envVar := parts[0]
			defaultValue := ""
			if len(parts) == 2 {
				defaultValue = parts[1]
			}

			var resolved string
			if val, exists := os.LookupEnv(envVar); exists {
				resolved = val
			} else if internalVal := c.v.GetString(envVar); internalVal != "" {
				resolved, _ = c.resolvePlaceholders(internalVal) // resolve internal reference
			} else {
				resolved = defaultValue
			}

			value = strings.Replace(value, fullMatch, resolved, -1)
		}
	}
	return value, nil
}

// GetPropertyString function to get a string property
func (c *config) GetPropertyString(name string) string {
	value := c.v.GetString(name)
	resolved, err := c.resolvePlaceholders(value)
	if err != nil {
		slog.Error("Error resolving placeholders property:", name, err.Error())
		return ""
	}
	return resolved
}

// GetPropertyBool function to get a boolean property
func (c *config) GetPropertyBool(name string) bool {
	value := c.GetPropertyString(name)
	return strings.ToLower(value) == "true"
}

// GetPropertyTime function to get a time property
func (c *config) GetPropertyTime(name string) time.Time {
	value := c.GetPropertyString(name)
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		slog.Error("Error parsing time property:", name, err.Error())
		return time.Time{}
	}
	return parsedTime
}

// GetPropertyStringSlice function to get a string slice property
func (c *config) GetPropertyStringSlice(name string) []string {
	value := c.v.GetStringSlice(name)
	var result []string
	for _, v := range value {
		resolved, err := c.resolvePlaceholders(v)
		if err != nil {
			slog.Error("Error resolving placeholders property:", name, err.Error())
			return nil
		}
		if resolved != "" {
			result = append(result, resolved)
		}
	}
	return result
}

// GetPropertyMapString function to get a map of strings property
func (c *config) GetPropertyMapString(name string) map[string]string {
	value := c.v.GetStringMapString(name)
	result := make(map[string]string)
	for k, v := range value {
		resolved, err := c.resolvePlaceholders(v)
		if err != nil {
			slog.Error("Error resolving placeholders property:", name, err.Error())
			return nil
		}
		result[k] = resolved
	}
	return result
}

// AddConfigFile add a new config file yaml
func (c *config) AddConfigFile(name string) {
	c.v.SetConfigFile(name)
}
