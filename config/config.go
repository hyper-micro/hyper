package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type Config interface {
	GetFloat64(key string) float64
	GetFloat64OrDefault(key string, def float64) float64
	GetBool(key string) bool
	GetBoolOrDefault(key string, def bool) bool
	GetString(key string) string
	GetStringOrDefault(key string, def string) string
	GetInt(key string) int
	GetIntOrDefault(key string, def int) int
	GetIntSlice(key string) []int
	GetIntSliceOrDefault(key string, def []int) []int
	GetStringMap(key string) map[string]interface{}
	GetStringMapOrDefault(key string, def map[string]interface{}) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringOrDefault(key string, def map[string]string) map[string]string
	GetStringSlice(key string) []string
	GetStringSliceOrDefault(key string, def []string) []string
	GetTime(key string) time.Time
	GetTimeOrDefault(key string, def time.Time) time.Time
	GetDuration(key string) time.Duration
	GetDurationOrDefault(key string, def time.Duration) time.Duration
	Get(key string) interface{}
}

type FileType uint8

const (
	UnknownFileType FileType = iota
	YamlFileType
	TomlFileType
	IniFileType
	JsonFileType
)

const (
	PathTypeFile = 1
	PathTypePath = 2
)

const (
	defaultKeyDelim = "."
)

type config struct {
	kv      map[string]map[string]interface{}
	kvCache *sync.Map
}

func New(pathType int, paths ...string) (Config, error) {
	conf := &config{
		kv:      make(map[string]map[string]interface{}),
		kvCache: new(sync.Map),
	}
	var err error
	if pathType == PathTypeFile {
		err = conf.loadConfigs(paths...)
	} else {
		err = conf.loadPaths(paths...)
	}
	return conf, err
}

func (c *config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.getValue(key))
}

func (c *config) GetFloat64OrDefault(key string, def float64) float64 {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToFloat64(val)
}

func (c *config) GetBool(key string) bool {
	return cast.ToBool(c.getValue(key))
}

func (c *config) GetBoolOrDefault(key string, def bool) bool {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToBool(val)
}

func (c *config) GetString(key string) string {
	return cast.ToString(c.getValue(key))
}

func (c *config) GetStringOrDefault(key string, def string) string {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToString(val)
}

func (c *config) GetInt(key string) int {
	return cast.ToInt(c.getValue(key))
}

func (c *config) GetIntOrDefault(key string, def int) int {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToInt(val)
}

func (c *config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.getValue(key))
}

func (c *config) GetIntSliceOrDefault(key string, def []int) []int {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToIntSlice(val)
}

func (c *config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.getValue(key))
}

func (c *config) GetStringMapOrDefault(key string, def map[string]interface{}) map[string]interface{} {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToStringMap(val)
}

func (c *config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.getValue(key))
}

func (c *config) GetStringMapStringOrDefault(key string, def map[string]string) map[string]string {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToStringMapString(val)
}

func (c *config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.getValue(key))
}

func (c *config) GetStringSliceOrDefault(key string, def []string) []string {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToStringSlice(val)
}

func (c *config) GetTime(key string) time.Time {
	return cast.ToTime(c.getValue(key))
}

func (c *config) GetTimeOrDefault(key string, def time.Time) time.Time {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToTime(val)
}

func (c *config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.getValue(key))
}

func (c *config) GetDurationOrDefault(key string, def time.Duration) time.Duration {
	val := c.getValue(key)
	if val == nil {
		return def
	}
	return cast.ToDuration(val)
}

func (c *config) Get(key string) interface{} {
	return c.getValue(key)
}

func (c *config) getValue(key string) interface{} {
	lk := strings.ToLower(key)
	if cacheVal, ok := c.kvCache.Load(lk); ok {
		return cacheVal
	}
	keys := strings.Split(lk, defaultKeyDelim)
	val := c.getValueFromMaps(keys)
	c.kvCache.Store(lk, val)
	return val
}

func (c *config) loadPaths(paths ...string) error {
	if len(paths) == 0 {
		return errors.New("configuration path not set")
	}
	for _, p := range paths {
		p = formatPathSeparator(p)
		p = strings.TrimRight(p, string(os.PathSeparator))

		files, _ := ioutil.ReadDir(p)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			configFile := strings.Join([]string{p, file.Name()}, string(os.PathSeparator))
			_ = c.loadConfig(configFile)
		}
	}
	return nil
}

func (c *config) loadConfigs(configFiles ...string) error {
	if len(configFiles) == 0 {
		return errors.New("configuration file not set")
	}

	for _, file := range configFiles {
		file = formatPathSeparator(file)
		if err := c.loadConfig(file); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) loadConfig(configFile string) error {
	var (
		fileFullName = path.Base(configFile)
		fileExt      = path.Ext(configFile)
		fileName     = fileFullName[0 : len(fileFullName)-len(fileExt)]
	)

	fileExt = strings.Trim(fileExt, ".")

	if fileName == "" {
		return fmt.Errorf("loadConfig %s: file name cannot be empty", configFile)
	}

	fileType := c.getFileTypeByExtension(fileExt)
	if fileType == UnknownFileType {
		return fmt.Errorf("loadConfig %s: file type is not supported", configFile)
	}

	b, err := c.readLocalFile(configFile)
	if err != nil {
		return err
	}

	kv := make(map[string]interface{})
	if err := c.decodeReader(b, kv, fileType); err != nil {
		return err
	}

	mapsKey2Lower(kv)
	c.kv[fileName] = kv

	return nil
}

func (c *config) decodeReader(b []byte, cfg map[string]interface{}, fileType FileType) error {
	dc, ok := decoders[fileType]
	if !ok {
		panic(fmt.Sprintf("fileType %v no decoder", fileType))
	}
	return dc.Decode(b, cfg)
}

func (c *config) getFileTypeByExtension(ext string) FileType {
	switch strings.ToLower(ext) {
	case "yaml":
		return YamlFileType
	case "yml":
		return YamlFileType
	case "toml":
		return TomlFileType
	case "ini":
		return IniFileType
	case "json":
		return JsonFileType
	default:
		return UnknownFileType
	}
}

func (c *config) readLocalFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *config) getValueFromMaps(keys []string) interface{} {
	if len(keys) == 0 {
		return nil
	}
	fileKey := keys[0]
	mapKey := keys[1:]
	var (
		val  interface{} = c.kv[fileKey]
		nval map[string]interface{}
		ok   bool
	)
	for _, k := range mapKey {
		nval, ok = val.(map[string]interface{})
		if !ok {
			return nil
		}
		val, ok = nval[k]
		if !ok {
			return nil
		}
	}
	return val
}

func mapsKey2Lower(kv map[string]interface{}) {
	for k, v := range kv {
		switch v.(type) {
		case map[interface{}]interface{}:
			v = cast.ToStringMap(v)
			mapsKey2Lower(v.(map[string]interface{}))
		case map[string]interface{}:
			mapsKey2Lower(v.(map[string]interface{}))
		}
		lk := strings.ToLower(k)
		if lk != k {
			delete(kv, k)
			kv[lk] = v
		}
	}
}

func formatPathSeparator(p string) string {
	p = strings.ReplaceAll(p, "/", string(os.PathSeparator))
	p = strings.ReplaceAll(p, "\\", string(os.PathSeparator))
	return p
}
