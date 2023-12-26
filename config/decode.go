package config

import "github.com/hyper-micro/hyper/config/codec"

type Decoder interface {
	Decode([]byte, map[string]interface{}) error
}

var decoders = make(map[FileType]Decoder)

func RegisterDecoder(fileType FileType, decoder Decoder) {
	decoders[fileType] = decoder
}

func init() {
	RegisterDecoder(YamlFileType, codec.YamlCodec{})
	RegisterDecoder(IniFileType, codec.IniCodec{KeyDelimiter: defaultKeyDelim})
	RegisterDecoder(JsonFileType, codec.JsonCodec{})
	RegisterDecoder(TomlFileType, codec.TomlCodec{})
}
