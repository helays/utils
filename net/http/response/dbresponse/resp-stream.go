package dbresponse

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	"helay.net/go/utils/v3/close/sqlClose"
	"helay.net/go/utils/v3/logger/ulogs"
	"helay.net/go/utils/v3/net/http/response"
)

type Stream[T any] struct {
	Tx           *gorm.DB
	Rows         *sql.Rows
	EnableTotals bool
	Totals       int64
}

// RespLists
//
//	{
//	   "code":0,
//	   "msg":"成功",
//	   "total":0,
//	   "data":[]
//	}
func (s *Stream[T]) RespLists(w http.ResponseWriter) {
	defer sqlClose.CloseRows(s.Rows)
	response.RespJson(w)
	_, _ = w.Write([]byte(`{"code":0,"msg":"成功",`))
	if s.EnableTotals {
		_, _ = w.Write([]byte(`"total":0,`))
	}
	_, _ = w.Write([]byte(`"data":[`))
	isFirst := true
	for s.Rows.Next() {
		var result T
		if err := s.Tx.ScanRows(s.Rows, &result); err != nil {
			ulogs.Debugf("[流式响应] Scan 错误 %v", err)
			continue
		}
		if !isFirst {
			_, _ = w.Write([]byte(","))
		}
		isFirst = false
		enc := json.NewEncoder(w)
		if err := enc.Encode(result); err != nil {
			ulogs.Debug("[流式响应] Json Encode 失败 %v", err)
		}
	}
	_, _ = w.Write([]byte(`]}`))
}
