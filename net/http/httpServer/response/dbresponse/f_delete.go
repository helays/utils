package dbresponse

import (
	"net/http"

	"github.com/helays/utils/v2/db/userDb"
	"github.com/helays/utils/v2/net/http/httpServer/response"
	"gorm.io/gorm"
)

// ModelDelete 执行模型删除操作。
// 该函数接收一个HTTP响应写入器、一个HTTP请求、一个数据库事务指针、一个源模型实例和一个查询配置。
// 它根据查询配置在数据库中删除相关的模型记录。
// 如果删除操作失败，它将返回一个500错误和"删除数据失败"的消息。
// 如果删除成功，它将返回状态码0和"成功"的消息。
func ModelDelete[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, c userDb.Curd) {
	if !response.CheckReqPost(w, r) {
		return
	}
	_tx := tx.Model(new(T))
	if c.Where.Query != "" {
		_tx.Where(c.Where.Query, c.Where.Args...)
	}
	err := _tx.Delete(nil).Error
	if err != nil {
		response.SetReturnError(w, r, err, 500, "销毁数据失败")
		return
	}
	response.SetReturnData(w, 0, "成功")
}

// PostModelUpdate 执行模型删除操作。是通过更新状态字段的方式实现的
func PostModelUpdate[T any](w http.ResponseWriter, r *http.Request, tx *gorm.DB, c userDb.Curd, errMsg string) {
	if !response.CheckReqPost(w, r) {
		return
	}
	if err := userDb.UpdateWithoutValid[T](tx, c); err != nil {
		response.SetReturnError(w, r, err, 500, errMsg)
		return
	}
	response.SetReturnCode(w, r, 0, "成功")
}
