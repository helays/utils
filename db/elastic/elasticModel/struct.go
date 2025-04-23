package elasticModel

// ESSearchResponse 包含成功和错误响应的统一结构
type ESSearchResponse struct {
	// 请求元信息
	Took     int  `json:"took"`      // 查询耗时(毫秒)
	TimedOut bool `json:"timed_out"` // 是否超时

	// 分片信息
	Shards *ShardsInfo `json:"_shards"` // 分片信息

	// 查询结果
	Hits *HitsContainer `json:"hits"` // 命中结果

	// 聚合结果
	Aggregations map[string]interface{} `json:"aggregations"` // 聚合结果

	// 分页/滚动
	ScrollID string `json:"_scroll_id"` // 滚动查询ID
	PitID    string `json:"pit_id"`     // Point in Time ID

	ESError
}

// ShardsInfo 分片信息
type ShardsInfo struct {
	Total      int `json:"total"`      // 总分片数
	Successful int `json:"successful"` // 成功分片数
	Skipped    int `json:"skipped"`    // 跳过分片数
	Failed     int `json:"failed"`     // 失败分片数
}

// HitsContainer 命中结果容器
type HitsContainer struct {
	Total    *TotalHits `json:"total"`     // 总命中数信息
	MaxScore float64    `json:"max_score"` // 最大得分
	Hits     []*Hit     `json:"hits"`      // 命中文档列表
}

// TotalHits 总命中数信息
type TotalHits struct {
	Value    int    `json:"value"`    // 命中总数
	Relation string `json:"relation"` // 关系("eq"精确/"gte"估算)
}

// Hit 单个命中文档
type Hit struct {
	Index       string              `json:"_index"`        // 索引名
	ID          string              `json:"_id"`           // 文档ID
	Score       float64             `json:"_score"`        // 相关性得分
	Source      map[string]any      `json:"_source"`       // 文档原始数据
	Version     int                 `json:"_version"`      // 版本号
	SeqNo       int64               `json:"_seq_no"`       // 序列号
	PrimaryTerm int64               `json:"_primary_term"` // 主任期
	Highlight   map[string][]string `json:"highlight"`     // 高亮结果
	Fields      map[string]any      `json:"fields"`        // 字段数据
	Sort        []any               `json:"sort"`          // 排序值
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Type         string         `json:"type"`          // 错误类型
	Reason       string         `json:"reason"`        // 可读的错误原因
	ResourceType string         `json:"resource.type"` // 相关资源类型
	ResourceID   string         `json:"resource.id"`   // 相关资源ID
	Index        string         `json:"index"`         // 相关索引
	Phase        string         `json:"phase"`         // 错误阶段
	RootCause    []*ErrorDetail `json:"root_cause"`    // 根本原因列表
	CausedBy     map[string]any `json:"caused_by"`     // 导致错误的底层原因
	Metadata     map[string]any `json:"metadata"`      // 错误元数据
}

// HasPartialResults 判断是否有部分结果(部分分片失败)
func (r *ESSearchResponse) HasPartialResults() bool {
	return r.Shards != nil && r.Shards.Failed > 0 && r.Hits != nil
}

// GetErrorMessage 获取错误消息
func (r *ESSearchResponse) GetErrorMessage() string {
	if r.Error == nil {
		return ""
	}

	if len(r.Error.RootCause) > 0 {
		return r.Error.RootCause[0].Reason
	}
	return r.Error.Reason
}

func (r *ESSearchResponse) GetTotal() int {
	if r.Hits == nil || r.Hits.Total == nil {
		return 0
	}
	return r.Hits.Total.Value
}

func (r *ESSearchResponse) GetHits() []*Hit {
	if r.Hits == nil {
		return nil
	}
	return r.Hits.Hits
}
