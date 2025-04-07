package elastic

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/helays/utils/tools"
	"log"
	"sync"
	"time"
)

// BulkIndexerConfig 批量索引器配置
type BulkIndexerConfig struct {
	FlushBytes    int           `json:"flush_bytes" yaml:"flush_bytes" ini:"flush_bytes"`          // 刷新字节阈值
	FlushInterval time.Duration `json:"flush_interval" yaml:"flush_interval" ini:"flush_interval"` // 刷新时间间隔
	NumWorkers    int           `json:"num_workers" yaml:"num_workers" ini:"num_workers"`          // 工作协程数
	Silent        bool          `json:"silent" yaml:"silent" ini:"silent"`                         // 是否静默模式
	Client        *elasticsearch.Client
	Ctx           context.Context
	ErrLogr       func(err error, msg ...any)
}

// DefaultBulkIndexerConfig 返回默认的批量索引器配置
func DefaultBulkIndexerConfig(client *elasticsearch.Client, ctx context.Context) BulkIndexerConfig {
	return BulkIndexerConfig{
		FlushBytes:    5e+6, // 5MB
		FlushInterval: 30 * time.Second,
		NumWorkers:    4,
		Client:        client,
		Ctx:           ctx,
	}
}

// BulkError 表示索引错误信息
type BulkError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// BulkInsertResult 批量插入结果
type BulkInsertResult struct {
	SuccessCount int                              `json:"success_count"`
	FailedCount  int                              `json:"failed_count"`
	Responses    []esutil.BulkIndexerResponseItem `json:"responses"`
}

// BulkInsert 批量插入文档到Elasticsearch
func BulkInsert(index string, cfg BulkIndexerConfig, documents []map[string]any) (*BulkInsertResult, error) {
	if len(documents) == 0 {
		return &BulkInsertResult{}, nil
	}

	// 创建批量索引器
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index,
		Client:        cfg.Client,
		FlushBytes:    cfg.FlushBytes,
		FlushInterval: cfg.FlushInterval,
		NumWorkers:    cfg.NumWorkers,
	})
	if err != nil {
		return nil, fmt.Errorf("创建批量索引器失败: %w", err)
	}
	// 用于收集响应和同步的锁
	var (
		mu        sync.Mutex
		responses []esutil.BulkIndexerResponseItem
	)
	// 添加文档到批量索引器
	for _, doc := range documents {
		err = bi.Add(cfg.Ctx, esutil.BulkIndexerItem{
			Action: "index",
			Body:   tools.Any2Reader(doc),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				if !cfg.Silent {
					mu.Lock()
					defer mu.Unlock()
					responses = append(responses, res)
				}
			},
			OnFailure: func(_ context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				mu.Lock()
				defer mu.Unlock()
				responses = append(responses, res)
				if err != nil {
					cfg.ErrLogr(err, "文档索引失败", "索引:", index)
				} else {
					cfg.ErrLogr(fmt.Errorf("文档索引失败: %s: %s", res.Error.Type, res.Error.Reason), "索引:", index)
				}
			},
		})
		if err != nil {
			cfg.ErrLogr(err, "添加文档到批量索引器失败", "索引:", index)
			continue
		}
	}

	// 等待所有文档处理完成
	if err = bi.Close(cfg.Ctx); err != nil {
		return nil, fmt.Errorf("批量索引器关闭时出错: %w", err)
	}
	stats := bi.Stats()
	result := &BulkInsertResult{
		SuccessCount: int(stats.NumAdded),
		FailedCount:  int(stats.NumFailed),
		Responses:    responses,
	}

	return result, nil
}

// BulkInsertWithRetry 带重试的批量插入
func BulkInsertWithRetry(index string, cfg BulkIndexerConfig, documents []map[string]any, maxRetries int, retryInterval time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = BulkInsert(index, cfg, documents)
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			log.Printf("批量插入失败，准备重试 (%d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
		}
	}
	return fmt.Errorf("经过 %d 次重试后仍然失败: %w", maxRetries, err)
}
