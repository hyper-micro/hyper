package codec

import "gopkg.in/yaml.v3"

type YamlCodec struct{}

func (YamlCodec) Decode(b []byte, v map[string]interface{}) error {
	return yaml.Unmarshal(b, &v)
}
