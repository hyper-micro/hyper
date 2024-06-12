package str

import (
	"github.com/hyper-micro/hyper/internal/json"
)

func JSONMarshal(v any) string {
	s, _ := json.Marshal(v)
	return string(s)
}
