package sqlClose

import "database/sql"

func CloseMysqlRows(rows *sql.Rows) {
	CloseRows(rows)
}

func CloseRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}
func CloseDb(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

// Deprecated: As of utils v1.1.0, this value is simply [tools.CloseDb].
func CloseMysql(conn *sql.DB) {
	CloseDb(conn)
}

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		_ = stmt.Close()
	}
}
