package postgresconnect

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"helay.net/go/utils/v3/db"
)

func InitDB(c *db.Dbbase) (*gorm.DB, error) {
	//postgres://user:password@host1:port1/database?target_session_attrs=read-write&TimeZone=Asia/Shanghai
	//dsn = "postgres://" + c.User + ":" + c.Pwd + "@" + strings.Join(c.Host, ",") + "/" + c.Dbname + "?TimeZone=Asia/Shanghai"
	dialector := postgres.New(postgres.Config{
		DSN:                  c.Dsn(),
		PreferSimpleProtocol: true,
	})

	return c.Connect(&dialector)
}
