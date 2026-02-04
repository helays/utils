package elasticModel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"helay.net/go/utils/v3/close/esClose"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/tools"
	"io"
	"time"
)

// ESScrollCallback 定义ES滚动查询回调函数类型
type ESScrollCallback func(hits []*Hit) error

type EsScroll struct {
	ESClient      *elasticsearch.Client // Elasticsearch客户端
	SearchRequest *esapi.SearchRequest  // 搜索请求

	ScrollTime time.Duration // 滚动查询保持时间
	BatchSize  int           // 每批处理的数据量
	Index      []string      // 查询索引
	QueryBody  io.Reader     // 查询体
	SortFields []string      // 排序字段
}

// DoESSearchWithScroll 执行ES滚动查询
func (this *EsScroll) DoESSearchWithScroll(ctx context.Context, callback ESScrollCallback) error {
	this.ScrollTime = tools.AutoTimeDuration(this.ScrollTime, time.Second, 10*time.Minute)
	// 构建SearchRequest
	if this.SearchRequest == nil {
		// 如果没有提供SearchRequest，则创建一个默认的
		this.SearchRequest = &esapi.SearchRequest{
			Index:   this.Index,
			Body:    this.QueryBody,
			Size:    &this.BatchSize,
			Scroll:  this.ScrollTime,
			Sort:    this.SortFields,
			Pretty:  false,
			Human:   false,
			Timeout: 30 * time.Second,
		}
	} else {
		// 确保SearchRequest有必要的参数
		if this.SearchRequest.Size == nil {
			this.SearchRequest.Size = &this.BatchSize
		}
		if this.SearchRequest.Scroll == 0 {
			this.SearchRequest.Scroll = this.ScrollTime
		}
	}
	// 执行初始搜索请求
	res, err := this.SearchRequest.Do(ctx, this.ESClient)
	if err != nil {
		return fmt.Errorf("执行搜索请求失败: %w", err)
	}
	defer esClose.CloseResp(res)

	if res.IsError() {
		return fmt.Errorf("搜索请求错误: %s", res.String())
	}
	// 解析初始响应
	var initialResponse ESSearchResponse
	if err = json.NewDecoder(res.Body).Decode(&initialResponse); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}
	hits := initialResponse.GetHits()
	hitLen := len(hits)

	// 处理第一批结果
	if hitLen > 0 {
		if err = callback(hits); err != nil {
			return err
		}
	} else {
		return nil // 没有数据直接返回
	}
	// 使用滚动查询获取剩余结果
	scrollID := initialResponse.ScrollID
	for {
		if scrollID == "" {
			break
		}

		scrollReq := esapi.ScrollRequest{
			ScrollID: scrollID,
			Scroll:   this.ScrollTime,
		}

		res, err = scrollReq.Do(ctx, this.ESClient)
		if err != nil {
			return fmt.Errorf("滚动查询失败: %w", err)
		}
		if res.IsError() {
			esClose.CloseResp(res)
			return fmt.Errorf("滚动查询错误: %s", res.String())
		}
		var scrollResponse ESSearchResponse
		if err = json.NewDecoder(res.Body).Decode(&scrollResponse); err != nil {
			esClose.CloseResp(res)
			return fmt.Errorf("解析滚动响应失败: %w", err)
		}
		hits = scrollResponse.GetHits()
		if len(hits) == 0 {
			esClose.CloseResp(res)
			break
		}
		// 处理当前批次结果
		if err = callback(hits); err != nil {
			esClose.CloseResp(res)
			return err
		}
		scrollID = scrollResponse.ScrollID
		esClose.CloseResp(res) // 立即关闭响应体
	}

	// 清除滚动上下文
	clearScrollReq := esapi.ClearScrollRequest{
		ScrollID: []string{scrollID},
	}
	_, err = clearScrollReq.Do(ctx, this.ESClient)
	if err != nil {
		ulogs.Warn("清除滚动上下文失败: %v", err)
	}
	return nil
}
