package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigIni(t *testing.T) {
	conf, err := New(PathTypeFile, "./testdata/config.ini")
	require.NoError(t, err)

	testBaseConfigData(t, conf)
	testIniConfigData(t, conf)
}

func TestConfigJson(t *testing.T) {
	conf, err := New(PathTypeFile, "./testdata/config.json")
	require.NoError(t, err)

	testBaseConfigData(t, conf)
	testNotIniConfigData(t, conf)
}

func TestConfigToml(t *testing.T) {
	conf, err := New(PathTypeFile, "./testdata/config.toml")
	require.NoError(t, err)

	testBaseConfigData(t, conf)
	testNotIniConfigData(t, conf)
}

func TestConfigYaml(t *testing.T) {
	conf, err := New(PathTypeFile, "./testdata/config.yaml")
	require.NoError(t, err)

	testBaseConfigData(t, conf)
	testNotIniConfigData(t, conf)
}

func testNotIniConfigData(t *testing.T, conf Config) {
	var (
		intSlice     = []int{1, 3, 5, 7}
		intSlice2    = []int{1, 3, 5}
		stringSlice  = []string{"a", "b", "c"}
		stringSlice2 = []string{"a", "b", "c", "d"}
	)

	assert.Equal(t, conf.Get("config.apiVersion"), "apps/v1")
	assert.Equal(t, conf.GetIntSlice("config.testData.numbers"), intSlice)
	assert.Equal(t, conf.GetStringSlice("config.testData.strings"), stringSlice)
	assert.Equal(t, conf.GetIntSliceOrDefault("config.testData.numbers", intSlice2), intSlice)
	assert.Equal(t, conf.GetIntSliceOrDefault("config.testData.notExist", intSlice), intSlice)
	assert.Equal(t, conf.GetStringSliceOrDefault("config.testData.strings", stringSlice2), stringSlice)
	assert.Equal(t, conf.GetStringSliceOrDefault("config.testData.notExist", stringSlice), stringSlice)
}

func testIniConfigData(t *testing.T, conf Config) {
	assert.Equal(t, conf.Get("config.default.apiVersion"), "apps/v1")
}

func testBaseConfigData(t *testing.T, conf Config) {
	var (
		stringMapString = map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		defaultStringMapString = map[string]string{
			"key11": "value11",
			"key22": "value22",
		}
		stringMap = map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		defaultStringMap = map[string]interface{}{
			"key11": "value11",
			"key22": "value22",
		}
		emptyStringMap       = map[string]interface{}{}
		emptyStringMapString = map[string]string{}
		emptyTime            = time.Time{}
		testTime, _          = time.Parse(time.RFC3339, "2022-04-19T13:15:58Z")
	)

	assert.Equal(t, conf.GetBool("config.testData.switch"), true)
	assert.Equal(t, conf.GetString("config.testData.name"), "test")
	assert.Equal(t, conf.GetTime("config.testData.time"), testTime)
	assert.Equal(t, conf.GetDuration("config.testData.duration"), time.Duration(100))
	assert.Equal(t, conf.GetStringMap("config.stringMap"), stringMap)
	assert.Equal(t, conf.GetStringMapString("config.stringMap"), stringMapString)
	assert.Equal(t, conf.GetInt("config.testData.number"), 102400)
	assert.Equal(t, conf.GetFloat64("config.testData.pi"), 3.1415926)

	assert.Equal(t, conf.GetFloat64OrDefault("config.testData.pi", 3.14), 3.1415926)
	assert.Equal(t, conf.GetFloat64OrDefault("config.testData.notExist", 3.14), 3.14)
	assert.Equal(t, conf.GetFloat64OrDefault("config.testData.zero", 3.14), 0.0)

	assert.Equal(t, conf.GetBoolOrDefault("config.testData.switch", false), true)
	assert.Equal(t, conf.GetBoolOrDefault("config.testData.notExist", false), false)

	assert.Equal(t, conf.GetStringOrDefault("config.testData.name", "test1"), "test")
	assert.Equal(t, conf.GetStringOrDefault("config.testData.notExist", "test1"), "test1")
	assert.Equal(t, conf.GetStringOrDefault("config.testData.emptyString", "test1"), "")

	assert.Equal(t, conf.GetIntOrDefault("config.testData.number", 1024), 102400)
	assert.Equal(t, conf.GetIntOrDefault("config.testData.notExist", 1024), 1024)
	assert.Equal(t, conf.GetIntOrDefault("config.testData.zero", 1024), 0)

	assert.Equal(t, conf.GetStringMapOrDefault("config.stringMap", defaultStringMap), stringMap)
	assert.Equal(t, conf.GetStringMapOrDefault("config.notExist", defaultStringMap), defaultStringMap)
	assert.Equal(t, conf.GetStringMapOrDefault("config.emptyStringMap", defaultStringMap), emptyStringMap)

	assert.Equal(t, conf.GetStringMapStringOrDefault("config.stringMap", defaultStringMapString), stringMapString)
	assert.Equal(t, conf.GetStringMapStringOrDefault("config.notExist", defaultStringMapString), defaultStringMapString)
	assert.Equal(t, conf.GetStringMapStringOrDefault("config.emptyStringMap", defaultStringMapString), emptyStringMapString)

	assert.Equal(t, conf.GetTimeOrDefault("config.testData.time", emptyTime), testTime)
	assert.Equal(t, conf.GetTimeOrDefault("config.testData.notExist", testTime), testTime)
	assert.Equal(t, conf.GetTimeOrDefault("config.testData.zeroTime", testTime), emptyTime)

	assert.Equal(t, conf.GetDurationOrDefault("config.testData.duration", time.Duration(10)), time.Duration(100))
	assert.Equal(t, conf.GetDurationOrDefault("config.testData.notExist", time.Duration(10)), time.Duration(10))
	assert.Equal(t, conf.GetDurationOrDefault("config.testData.zero", time.Duration(10)), time.Duration(0))
}
