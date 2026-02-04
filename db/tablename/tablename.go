package tablename

// TableName 表名结构需要实现的方法
type TableName interface {
	String() string
	MigrateError() string
	MigrateComment() (string, string)
	Comment() string
}
