package kafka

import (
	"database/sql/driver"
	"github.com/helays/utils/dataType"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

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

func (this KafkaConfig) Value() (driver.Value, error) {
	return dataType.DriverValueWithJson(this)
}

func (this *KafkaConfig) Scan(val interface{}) error {
	return dataType.DriverScanWithJson(val, this)
}

func (this KafkaConfig) GormDataType() string {
	return "json"
}

func (KafkaConfig) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
}
