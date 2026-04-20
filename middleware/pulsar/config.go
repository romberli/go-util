package pulsar

type Config struct {
	URL   string
	Token string
	Topic string
}

func NewConfig(url, token, topic string) *Config {
	return &Config{
		URL:   url,
		Token: token,
		Topic: topic,
	}
}

func (c *Config) Clone() *Config {
	return &Config{
		URL:   c.URL,
		Token: c.Token,
		Topic: c.Topic,
	}
}
