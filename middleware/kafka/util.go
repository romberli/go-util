package kafka

import (
	"github.com/IBM/sarama"
)

func ConvertHeaderToMap(header sarama.RecordHeader) map[string]string {
	return map[string]string{string(header.Key): string(header.Value)}
}
