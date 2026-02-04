package kafka

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/dataType"
)

// KafkaMessageTypeEnum
// noinspection all
var KafkaMessageTypeEnum = map[string]string{
	sarama.SASLTypePlaintext:   sarama.SASLTypePlaintext,
	sarama.SASLTypeSCRAMSHA256: sarama.SASLTypeSCRAMSHA256,
	sarama.SASLTypeSCRAMSHA512: sarama.SASLTypeSCRAMSHA512,
	sarama.SASLTypeGSSAPI:      sarama.SASLTypeGSSAPI,
	sarama.SASLTypeOAuth:       sarama.SASLTypeOAuth,
}

// noinspection all
type KafkaConfig struct {
	Addrs       []string      `yaml:"addrs" json:"addrs" ini:"addrs,omitempty"`
	Version     string        `yaml:"version" json:"version" ini:"version"` // kafka版本
	Sasl        bool          `yaml:"sasl" json:"sasl" ini:"sasl"`
	User        string        `yaml:"user" json:"user" ini:"user"`
	Password    string        `yaml:"password" json:"password" ini:"password"`
	Mechanism   string        `yaml:"mechanism" json:"mechanism" ini:"mechanism"`
	Offset      int64         `yaml:"offset" json:"offset" ini:"offset"`                // 默认从最新开始消费 -1 -2从最后
	MaxRetry    int           `yaml:"max_retry" json:"max_retry" ini:"max_retry"`       // 生产消息失败，默认重试3次
	Timeout     time.Duration `json:"timeout" yaml:"timeout" ini:"timeout"`             // 超时时间
	Compression bool          `json:"compression" yaml:"compression" ini:"compression"` // 发送消息是否开启压缩
	// 这里的kafka无复杂业务，可以用下方的相关配置
	ProducerMessage ProducerMessage `json:"producer_message" yaml:"producer_message" ini:"producer_message"`
	// 生产者配置
	//消费者配置
	GroupName string `json:"group_name" yaml:"group_name" ini:"group_name"`
}

type ProducerMessage struct {
	Topic  string            `json:"topic" yaml:"topic" ini:"topic"`
	Key    string            `json:"key" yaml:"key" ini:"key"`
	Header map[string]string `json:"header" yaml:"header" ini:"header"`
	Role   string            `json:"role" yaml:"role" ini:"role"` // 同步生产者 或者异步生产者
}

// noinspection all
func (kc KafkaConfig) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(kc)
}

// noinspection all
func (kc *KafkaConfig) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, kc)
}

// noinspection all
func (kc KafkaConfig) GormDataType() string {
	return "json"
}

// noinspection all
func (KafkaConfig) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}

// noinspection all
func (kc *KafkaConfig) Valid() error {
	if len(kc.Addrs) < 1 {
		return errors.New("kafka地址不能为空")
	}
	if kc.Sasl {
		if _, ok := KafkaMessageTypeEnum[kc.Mechanism]; !ok {
			return errors.New("sasl机制错误")
		}
		if kc.User == "" {
			return errors.New("sasl用户名不能为空")
		}
		if kc.Password == "" {
			return errors.New("sasl密码不能为空")
		}
	}
	return nil
}

// noinspection all
func (kc *KafkaConfig) RemovePasswd() {
	kc.Password = ""
}

// noinspection all
func (kc *KafkaConfig) SetInfo(args ...any) {
	if len(args) != 2 {
		return
	}
	switch args[0].(string) {
	case config.ClientInfoHost:
		kc.Addrs = args[1].([]string)
	case config.ClientInfoUser:
		kc.User = args[1].(string)
	case config.ClientInfoPasswd:
		kc.Password = args[1].(string)
	}
}
