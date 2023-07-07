package kafka

import (
	"github.com/Shopify/sarama"
)

func ConvertHeaderToMap(header sarama.RecordHeader) map[string]string {
	return map[string]string{string(header.Key): string(header.Value)}
}
