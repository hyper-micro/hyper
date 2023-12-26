package codec

import "github.com/BurntSushi/toml"

type TomlCodec struct{}

func (TomlCodec) Decode(b []byte, v map[string]interface{}) error {
	return toml.Unmarshal(b, &v)
}
