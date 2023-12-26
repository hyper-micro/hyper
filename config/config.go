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

type FileType uint8

const (
	UnknownFileType FileType = iota
	YamlFileType
	TomlFileType
	IniFileType
	JsonFileType
)

const (
	defaultKeyDelim = "."
)

type Config struct {
	kv      map[string]map[string]interface{}
	kvCache *sync.Map
}

func New() *Config {
	return &Config{
		kv:      make(map[string]map[string]interface{}),
		kvCache: new(sync.Map),
	}
}

func (c *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(c.getValue(key))
}

func (c *Config) GetBool(key string) bool {
	return cast.ToBool(c.getValue(key))
}

func (c *Config) GetString(key string) string {
	return cast.ToString(c.getValue(key))
}

func (c *Config) GetInt(key string) int {
	return cast.ToInt(c.getValue(key))
}

func (c *Config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(c.getValue(key))
}

func (c *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(c.getValue(key))
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(c.getValue(key))
}

func (c *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(c.getValue(key))
}

func (c *Config) GetTime(key string) time.Time {
	return cast.ToTime(c.getValue(key))
}

func (c *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(c.getValue(key))
}

func (c *Config) Get(key string) interface{} {
	return c.getValue(key)
}

func (c *Config) getValue(key string) interface{} {
	lk := strings.ToLower(key)
	if cacheVal, ok := c.kvCache.Load(lk); ok {
		return cacheVal
	}
	keys := strings.Split(lk, defaultKeyDelim)
	val := c.getValueFromMaps(keys)
	c.kvCache.Store(lk, val)
	return val
}

func (c *Config) LoadPaths(paths ...string) error {
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

func (c *Config) LoadConfigs(configFiles ...string) error {
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

func (c *Config) loadConfig(configFile string) error {
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

func (c *Config) decodeReader(b []byte, cfg map[string]interface{}, fileType FileType) error {
	dc, ok := decoders[fileType]
	if !ok {
		panic(fmt.Sprintf("fileType %v no decoder", fileType))
	}
	return dc.Decode(b, cfg)
}

func (c *Config) getFileTypeByExtension(ext string) FileType {
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

func (c *Config) readLocalFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Config) getValueFromMaps(keys []string) interface{} {
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
