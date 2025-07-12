package elastic

import (
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/helays/utils/v2/close/esClose"
	"github.com/helays/utils/v2/config"
	"github.com/helays/utils/v2/tools"
)

type Index struct {
	IndexName     string            `json:"index_name" yaml:"index_name" ini:"index_name"`                // 索引名
	ExtName       string            `json:"ext_name" yaml:"ext_name" ini:"ext_name"`                      // 扩展名
	IlmPolicyName string            `json:"ilm_policy_name" yaml:"ilm_policy_name" ini:"ilm_policy_name"` // ilm策略名
	Mappings      map[string]string `json:"mappings" yaml:"mappings" ini:"mappings"`                      // 索引映射
}

// Create 创建索引
func (this Index) Create(client *elasticsearch.Client) error {
	if this.IndexName == "" {
		return fmt.Errorf("索引名不能为空")
	}
	idxName := this.IndexName
	if idxName != "" {
		idxName += "-" + this.ExtName
	}
	if err := this.exists(client, idxName); err == nil || errors.Is(err, config.ErrNotFound) {
		return err
	}
	body := map[string]any{
		"aliases": map[string]any{
			idxName: map[string]any{"is_write_index": true},
		},
	}
	// 配置ilm策略
	if this.IlmPolicyName != "" {
		body["settings"] = map[string]any{
			"index": map[string]any{
				"lifecycle": map[string]any{
					"name":           this.IlmPolicyName,
					"rollover_alias": idxName,
				},
			},
		}
	}
	// 配置索引映射
	if this.Mappings != nil && len(this.Mappings) > 0 {
		var mapping = make(map[string]any)
		for k, v := range this.Mappings {
			mapping[k] = map[string]any{"type": v}
		}
		body["mappings"] = map[string]any{
			"properties": mapping,
		}
	}
	req := esapi.IndicesCreateRequest{
		Index: idxName,
		Body:  tools.Any2Reader(body),
	}
	resp, err := req.Do(context.Background(), client)
	defer esClose.CloseResp(resp)
	if err != nil {
		return err
	}
	if !resp.IsError() {
		return nil
	}
	return fmt.Errorf("创建索引失败：%s", resp.String())
}

// 判断索引是否存在
func (this Index) exists(client *elasticsearch.Client, idxName string) error {
	req := esapi.IndicesExistsRequest{
		Index: []string{idxName},
	}
	resp, err := req.Do(context.Background(), client)
	defer esClose.CloseResp(resp)
	if err != nil {
		return err
	}
	if !resp.IsError() {
		return nil
	}
	if resp.StatusCode == 404 {
		return config.ErrNotFound
	}
	return fmt.Errorf("判断索引存在失败：%s", resp.String())
}
