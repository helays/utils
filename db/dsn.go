package db

import (
	"fmt"
	"github.com/helays/utils/config"
	"github.com/helays/utils/tools"
	"net/url"
	"strings"
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

func (d *Dbbase) sqliteDSN() string {
	if len(d.Host) < 1 {
		return ""
	}
	args := []string{
		"cache=shared",
		"mode=rwc",
		"_pragma=journal_mode(WAL)",
		"_pragma=synchronous(NORMAL)",
	}
	if d.Timeout > 0 {
		args = append(args, fmt.Sprintf("_pragma=busy_timeout(%d)", d.Timeout))
	}
	filePath := tools.Fileabs(d.Host[0])
	return fmt.Sprintf("file:%s?%s", filePath, strings.Join(args, "&"))
}
