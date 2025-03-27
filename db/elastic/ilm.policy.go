package elastic

import (
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/helays/utils/close/esClose"
	"github.com/helays/utils/config"
	"github.com/helays/utils/tools"
)

type IlmPolicy struct {
	IlmPolicyName       string `json:"ilm_policy_name" yaml:"ilm_policy_name" ini:"ilm_policy_name"`                      // ilm策略名称
	MaxAge              string `json:"max_age" yaml:"max_age" ini:"max_age"`                                              // 最大存在时间
	MaxPrimaryShardDocs int    `json:"max_primary_shard_docs" yaml:"max_primary_shard_docs" ini:"max_primary_shard_docs"` // 最大文档数
	MaxPrimaryShardSize string `json:"max_primary_shard_size" yaml:"max_primary_shard_size" ini:"max_primary_shard_size"` // 最大主分片大小
	MaxDocs             int    `json:"max_docs" yaml:"max_docs" ini:"max_docs"`                                           // 最大文档数
	DeleteMinAge        string `json:"delete_min_age" yaml:"delete_min_age" ini:"delete_min_age"`                         // 处于删除阶段的数据留存时间 0d
}

// Create 创建生命周期策略
func (this IlmPolicy) Create(client *elasticsearch.Client) error {
	if err := this.exists(client); err == nil || !errors.Is(err, config.ErrNotFound) {
		return err
	}
	policy := map[string]any{
		"policy": map[string]any{
			"phases": map[string]any{
				"hot": map[string]any{
					"actions": map[string]any{
						"rollover": map[string]any{
							"max_age":                this.MaxAge,
							"max_primary_shard_size": this.MaxPrimaryShardSize,
							"max_primary_shard_docs": this.MaxPrimaryShardDocs,
							"max_docs":               this.MaxDocs,
						},
					},
					"min_age": "0ms",
				},
				"delete": map[string]any{
					"min_age": this.DeleteMinAge,
					"actions": map[string]any{
						"delete": map[string]any{"delete_searchable_snapshot": true},
					},
				},
			},
		},
	}
	// 需要创建策略
	req := esapi.ILMPutLifecycleRequest{
		Body:   tools.Any2Reader(policy),
		Policy: this.IlmPolicyName,
	}
	resp, err := req.Do(context.Background(), client)
	defer esClose.CloseResp(resp)
	if err != nil {
		return err
	}
	if !resp.IsError() {
		return nil
	}
	return fmt.Errorf("创建生命周期策略失败：%s")
}

func (this IlmPolicy) exists(client *elasticsearch.Client) error {
	req := esapi.ILMGetLifecycleRequest{Policy: this.IlmPolicyName}
	resp, err := req.Do(context.Background(), client)
	defer esClose.CloseResp(resp)
	if err != nil {
		return err
	}
	if !resp.IsError() {
		return nil // 策略存在
	}
	if resp.StatusCode == 404 {
		return config.ErrNotFound
	}
	return fmt.Errorf("判断策略存在失败：%s", resp.String())
}

func GetIlmPolicy(client *elasticsearch.Client) {

}
