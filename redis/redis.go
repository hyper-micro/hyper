package redis

type Config struct {
	Addr string
}

type Client struct {
	conf Config
}

func NewRedis(conf Config) *Client {
	return &Client{
		conf: conf,
	}
}
