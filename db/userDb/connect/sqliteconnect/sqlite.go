package sqliteconnect

import (
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	"helay.net/go/utils/v3/db"
)

// InitDB 使用纯 Go 实现的 SQLite 驱动（ncruces/go-sqlite3）初始化 GORM 连接。
// 相比 mattn/go-sqlite3，它通过 WebAssembly 运行 SQLite，不依赖 cgo，
// 因此可以直接交叉编译和静态编译。
// 注意：其 DSN 参数与 mattn 版本不同，PRAGMA 需要通过 "_pragma" 传入，详见 db.sqliteDSN。
func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	return c.Connect(new(gormlite.Open(c.Dsn())))
}
