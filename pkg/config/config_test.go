//go:build unit

package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/laurentpoirierfr/ora-cdc-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func createTestConfigFile(content string) (string, error) {
	file, err := os.CreateTemp("", "bootstrap.yaml")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestGetPropertyString(t *testing.T) {
	content := `
string_property: ${STRING_ENV:default_value}
internal_reference: ${string_property}/suffix
complex_property: ${STRING_ENV:http://default}/path/${internal_reference}/${NON_EXISTENT_ENV}/end
prop1: ${ENV:default}/autre/${prop2}
prop2: ma_valeur
`
	filename, err := createTestConfigFile(content)
	assert.NoError(t, err)
	defer os.Remove(filename)

	os.Setenv("STRING_ENV", "env_value")
	defer os.Unsetenv("STRING_ENV")

	config, err := config.NewConfig(filename)
	assert.NoError(t, err)

	value := config.GetPropertyString("string_property")
	assert.Equal(t, "env_value", value)

	internalValue := config.GetPropertyString("internal_reference")
	assert.Equal(t, "env_value/suffix", internalValue)

	complexValue := config.GetPropertyString("complex_property")
	assert.Equal(t, "env_value/path/env_value/suffix//end", complexValue)

	prop1Value := config.GetPropertyString("prop1")
	assert.Equal(t, "default/autre/ma_valeur", prop1Value)
}

func TestGetPropertyString_SkipEmptyVariables(t *testing.T) {
	content := `
internal_reference: ${string_property}
`
	filename, err := createTestConfigFile(content)
	assert.NoError(t, err)
	defer os.Remove(filename)

	// string_property => default_value
	// internal_reference => skipped

	config, err := config.NewConfig(filename)
	assert.NoError(t, err)

	internalValue := config.GetPropertyString("internal_reference")
	assert.Equal(t, "", internalValue)
}

func TestGetPropertyBool(t *testing.T) {
	content := `
bool_property_true: ${BOOL_TRUE_ENV:true}
bool_property_false: ${BOOL_FALSE_ENV:false}
`
	filename, err := createTestConfigFile(content)
	assert.NoError(t, err)
	defer os.Remove(filename)

	os.Setenv("BOOL_TRUE_ENV", "true")
	os.Setenv("BOOL_FALSE_ENV", "false")
	defer os.Unsetenv("BOOL_TRUE_ENV")
	defer os.Unsetenv("BOOL_FALSE_ENV")

	config, err := config.NewConfig(filename)
	assert.NoError(t, err)

	boolValueTrue := config.GetPropertyBool("bool_property_true")
	assert.True(t, boolValueTrue)

	boolValueFalse := config.GetPropertyBool("bool_property_false")
	assert.False(t, boolValueFalse)
}

func TestGetPropertyStringSlice(t *testing.T) {
	content := `
list:
  - bool_property_true
  - bool_property_false
  - ${env_bool_property_1}
  - ${env_bool_property_2}
`
	filename, err := createTestConfigFile(content)
	assert.NoError(t, err)
	defer os.Remove(filename)

	os.Setenv("env_bool_property_1", "bool_property_1")
	defer os.Unsetenv("env_bool_property_1")

	cfg, err := config.NewConfig(filename)
	assert.NoError(t, err)

	properties := cfg.GetPropertyStringSlice("list")
	assert.ElementsMatch(t, []string{"bool_property_true", "bool_property_false", "bool_property_1"}, properties)
}

func TestGetPropertyMapString(t *testing.T) {
	content := `
map:
  key_true: true
  key_false: false
  key_env_1: ${value_env_1}
  key_env_2: ${value_env_2}
`
	filename, err := createTestConfigFile(content)
	assert.NoError(t, err)
	defer os.Remove(filename)

	os.Setenv("value_env_1", "true")
	defer os.Unsetenv("value_env_1")

	cfg, err := config.NewConfig(filename)
	assert.NoError(t, err)

	properties := cfg.GetPropertyMapString("map")
	assert.Equal(t, map[string]string{
		"key_true":  "true",
		"key_false": "false",
		"key_env_1": "true",
		"key_env_2": "",
	}, properties)
}

func TestGetPropertyTime(t *testing.T) {
	content := `
time_property: ${TIME_ENV:2023-05-25T10:00:00Z}
`
	filename, err := createTestConfigFile(content)
	assert.NoError(t, err)
	defer os.Remove(filename)

	os.Setenv("TIME_ENV", "2024-05-25T10:00:00Z")
	defer os.Unsetenv("TIME_ENV")

	config, err := config.NewConfig(filename)
	assert.NoError(t, err)

	timeValue := config.GetPropertyTime("time_property")
	expectedTime, err := time.Parse(time.RFC3339, "2024-05-25T10:00:00Z")
	assert.NoError(t, err)
	assert.Equal(t, expectedTime, timeValue)
}
