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
	URLs      []string
	Token     string
	Topic     string
	NoLogging bool
}

func NewConfig(urls []string, token, topic string, noLogging bool) *Config {
	return &Config{
		URLs:      urls,
		Token:     token,
		Topic:     topic,
		NoLogging: noLogging,
	}
}

func (c *Config) Clone() *Config {
	urls := make([]string, len(c.URLs))
	copy(urls, c.URLs)
	return NewConfig(urls, c.Token, c.Topic, c.NoLogging)
}

func (c *Config) getURLString() (string, error) {
	var urlList []string
	for _, url := range c.URLs {
		url = strings.TrimSpace(url)
		if url != constant.EmptyString {
			urlList = append(urlList, url)
		}
	}
	if len(urlList) == constant.ZeroInt {
		return constant.EmptyString, errors.Errorf("pulsar url list is empty")
	}
	urls := strings.Join(urlList, constant.CommaString)
	if !strings.HasPrefix(urls, PulsarSchemePrefix) && !strings.HasPrefix(urls, PulsarSSLSchemePrefix) {
		urls = PulsarSchemePrefix + urls
	}

	return urls, nil
}
