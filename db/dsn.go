package db

import (
	"fmt"
	"net/url"
	"strings"

	"helay.net/go/utils/v3/config"
	"helay.net/go/utils/v3/tools"
)

func (this *Dbbase) Dsn() string {
	//dsn := url.URL{
	//	User: url.UserPassword(this.User, this.Pwd),
	//	Host: strings.Join(this.Host, ","),
	//	Path: this.Dbname,
	//}
	switch this.DbType {
	case config.DbTypePostgres, config.DbTypePg:
		return this.postgresqlDSN()
		//dsn.Scheme = "postgres"
		//// 如果下面这里 设置成TimeZone ，有几率会出现时间异常
		//dsn.RawQuery = fmt.Sprintf("search_path=%s&timezone=%s", this.Schema, "Asia/Shanghai")
		//return dsn.String()
	case config.DbTypeMysql:
		//dsn.Scheme = "mysql" // mysql 不需要这个
		// mysql 密码里面的特殊字符 不用序列化
		return this.mysqlDSN()
	case config.DbTypeSqlite:
		return this.sqliteDSN()
	}
	return ""
}

func (d *Dbbase) postgresqlDSN() string {
	hosts := make([]string, 0, len(d.Host))
	ports := make([]string, 0, len(d.Host))
	for _, addr := range d.Host {
		tmp := strings.Split(addr, ":")
		hosts = append(hosts, tmp[0])
		ports = append(ports, tmp[1])
	}
	var builds []string
	builds = append(builds, "host="+strings.Join(hosts, ","))
	builds = append(builds, "port="+strings.Join(ports, ","))
	builds = append(builds, "user="+d.User)
	builds = append(builds, "password="+d.Pwd)
	builds = append(builds, "dbname="+d.Dbname)
	builds = append(builds, "search_path="+d.Schema)
	builds = append(builds, "TimeZone=Asia/Shanghai")
	if d.PostgresOpt != nil {
		builds = append(builds, d.PostgresOpt.dsn()...)
	}
	return strings.Join(builds, " ")
}

func (d *Dbbase) mysqlDSN() string {
	dsn := url.URL{
		User: url.UserPassword(d.User, d.Pwd),
		Host: strings.Join(d.Host, ","),
		Path: d.Dbname,
	}
	query := dsn.Query()
	query.Set("charset", "utf8mb4")
	query.Set("parseTime", "True")
	query.Set("loc", "Local")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", d.User, d.Pwd, dsn.Host, d.Dbname, query.Encode())
}

// sqliteDSN 生成 ncruces/go-sqlite3 可识别的 DSN。
// 该驱动的 DSN 规则与 mattn/go-sqlite3 不同：
//   - 标准 SQLite URI 参数（mode、cache 等）原样保留；
//   - PRAGMA 必须通过重复的 "_pragma" 参数传入，不能直接写 "journal_mode=WAL"；
//   - 时间格式由 "_timefmt" 指定，没有 loc 参数（因此这里使用与 mattn
//     的 "_time_format=sqlite" 完全相同的布局，保证已有时区偏移的时间值读写一致；
//     旧数据里不含时区的 "2006-01-02 15:04:05" 字符串无法匹配该布局，
//     驱动会原样返回字符串，交由 dataType.CustomTime 按本地时区解析）。
func (d *Dbbase) sqliteDSN() string {
	filePath := ":memory:"
	if len(d.Host) > 0 && d.Host[0] != "" {
		filePath = tools.Fileabs(d.Host[0])
	}
	query := url.Values{}
	query.Set("mode", "rwc")
	query.Set("cache", "shared")
	// busy_timeout 需要放在其他 PRAGMA 之前生效
	if d.Timeout > 0 {
		query.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", d.Timeout))
	}
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "synchronous(NORMAL)")
	// mattn/go-sqlite3 的 SQLiteTimestampFormats[0]，保持时间写入格式不变
	query.Set("_timefmt", "2006-01-02 15:04:05.999999999-07:00")

	return fmt.Sprintf("file:%s?%s", filePath, query.Encode())
}
