package sonyflakekit

import (
	"sync"
	"time"

	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/sony/sonyflake/v2"
)

type SonyFlake struct {
	BitsSequence  int           `json:"bits_sequence" yaml:"bits_sequence" ini:"bits_sequence"`       // 序列号位数
	BitsMachineID int           `json:"bits_machine_id" yaml:"bits_machine_id" ini:"bits_machine_id"` // 机器ID位数
	TimeUnit      time.Duration `json:"time_unit" yaml:"time_unit" ini:"time_unit"`                   // 时间单位
	StartTime     time.Time     `json:"start_time" yaml:"start_time" ini:"start_time"`                // 起始时间
}
type IDGenerator struct {
	sf *sonyflake.Sonyflake
}

var (
	idInstance *IDGenerator
	idOnce     sync.Once
	cfg        = SonyFlake{
		BitsSequence:  12,
		BitsMachineID: 10,
		TimeUnit:      time.Millisecond,
		StartTime:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
)

func Init(c SonyFlake) {
	cfg = c
}

// NewIDGenerator 创建一个ID生成器
func NewIDGenerator() *IDGenerator {
	idOnce.Do(func() {
		settings := sonyflake.Settings{
			BitsSequence:  cfg.BitsSequence,
			BitsMachineID: cfg.BitsMachineID,
			TimeUnit:      cfg.TimeUnit,
			StartTime:     cfg.StartTime,
		}
		if sf, err := sonyflake.New(settings); err != nil {
			ulogs.DieCheckerr(err, "sonyflake初始化失败")
		} else {
			idInstance = &IDGenerator{
				sf: sf,
			}
		}
	})
	return idInstance
}

// GenerateID 生成ID
func (g *IDGenerator) GenerateID() (int64, error) {
	return g.sf.NextID()
}

// BatchGenerateID 批量生成ID
func (g *IDGenerator) BatchGenerateID(count int) ([]int64, error) {
	ids := make([]int64, 0, count)
	for i := 0; i < count; i++ {
		id, err := g.sf.NextID() // 安全并发调用
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
