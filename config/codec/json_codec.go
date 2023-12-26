package codec

import "encoding/json"

type JsonCodec struct{}

func (c JsonCodec) Decode(b []byte, v map[string]interface{}) error {
	return json.Unmarshal(b, &v)
}
