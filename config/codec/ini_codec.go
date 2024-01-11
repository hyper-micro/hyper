package codec

import (
	"strings"

	"gopkg.in/ini.v1"
)

type IniCodec struct {
	KeyDelimiter string
}

func (c IniCodec) Decode(b []byte, v map[string]interface{}) error {
	cfg, err := ini.Load(b)
	if err != nil {
		return err
	}

	sections := cfg.Sections()
	for i := 0; i < len(sections); i++ {
		section := sections[i]
		keys := section.Keys()
		secNamePaths := strings.Split(section.Name(), c.KeyDelimiter)
		m := v
		for _, path := range secNamePaths {
			if _, ok := m[path]; !ok {
				m[path] = make(map[string]interface{})
				m = m[path].(map[string]interface{})
				continue
			}
			m1, ok := m[path].(map[string]interface{})
			if !ok {
				m1 = make(map[string]interface{})
				m[path] = m1
			}
			m = m1
		}

		for j := 0; j < len(keys); j++ {
			key := keys[j]
			m[key.Name()] = key.Value()
		}
	}
	return nil
}
