package pulsar

import (
	"strings"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	PulsarSchemePrefix    = "pulsar://"
	PulsarSSLSchemePrefix = "pulsar+ssl://"
)

type Config struct {
	URLs  []string
	Token string
	Topic string
}

func NewConfig(urls []string, token, topic string) *Config {
	return &Config{
		URLs:  urls,
		Token: token,
		Topic: topic,
	}
}

func (c *Config) Clone() *Config {
	urls := make([]string, len(c.URLs))
	copy(urls, c.URLs)
	return NewConfig(urls, c.Token, c.Topic)
}

func (c *Config) getURLString() (string, error) {
	var urls []string
	for _, url := range c.URLs {
		url = strings.TrimSpace(url)
		if url == constant.EmptyString {
			continue
		}
		if !strings.HasPrefix(url, PulsarSchemePrefix) && !strings.HasPrefix(url, PulsarSSLSchemePrefix) {
			url = PulsarSchemePrefix + url
		}
		urls = append(urls, url)
	}
	if len(urls) == constant.ZeroInt {
		return constant.EmptyString, errors.Errorf("pulsar urls is empty")
	}
	return strings.Join(urls, constant.CommaString), nil
}