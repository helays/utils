package httpServerWithDb

import (
	"github.com/helays/utils/db/userDb"
	"github.com/helays/utils/http/httpServer"
	"gorm.io/gorm"
	"net/http"
)

// ModelDelete 执行模型删除操作。
// 该函数接收一个HTTP响应写入器、一个HTTP请求、一个数据库事务指针、一个源模型实例和一个查询配置。
// 它根据查询配置在数据库中删除相关的模型记录。
// 如果删除操作失败，它将返回一个500错误和"删除数据失败"的消息。
// 如果删除成功，它将返回状态码0和"成功"的消息。
func ModelDelete(w http.ResponseWriter, r *http.Request, tx *gorm.DB, src any, c userDb.QueryConfig) {
	if !httpServer.CheckReqPost(w, r) {
		return
	}
	_tx := tx.Model(src)
	if c.Query != nil {
		_tx.Where(c.Query, c.Args...)
	}
	err := _tx.Delete(nil).Error
	if err != nil {
		httpServer.SetReturnError(w, r, err, 500, "删除数据失败")
		return
	}
	httpServer.SetReturnData(w, 0, "成功")
}
