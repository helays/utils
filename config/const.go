package config

// interval 定义
const (
	IntervalSecond = "SECOND"
	IntervalMinute = "MINUTE"
	IntervalHour   = "HOUR"
	IntervalDay    = "DAY"
	IntervalWeek   = "WEEK"
	IntervalMonth  = "MONTH"
	IntervalYear   = "YEAR"

	IntervalSecondLabel = "秒"
	IntervalMinuteLabel = "分"
	IntervalHourLabel   = "时"
	IntervalDayLabel    = "天"
	IntervalWeekLabel   = "周"
	IntervalMonthLabel  = "月"
	IntervalYearLabel   = "年"
)

type Interval struct {
	Key   int    `json:"key" yaml:"key"`     // key
	Val   string `json:"val" yaml:"val"`     // val
	Label string `json:"Label" yaml:"Label"` //单位
}

var IntervalLists = []Interval{
	{
		Key:   1,
		Val:   IntervalSecond,
		Label: IntervalSecondLabel,
	},
	{
		Key:   60,
		Val:   IntervalMinute,
		Label: IntervalMinuteLabel,
	},
	{
		Key:   3600,
		Val:   IntervalHour,
		Label: IntervalHourLabel,
	},
	{
		Key:   86400,
		Val:   IntervalDay,
		Label: IntervalDayLabel,
	},
	{
		Key:   604800,
		Val:   IntervalWeek,
		Label: IntervalWeekLabel,
	},
	{
		Key:   2592000,
		Val:   IntervalMonth,
		Label: IntervalMonthLabel,
	},
	{
		Key:   31536000,
		Val:   IntervalYear,
		Label: IntervalYearLabel,
	},
}

var IntervalKeyMap = make(map[int]Interval)
var IntervalKeyValMap = make(map[string]Interval)
var IntervalKeyLabelMap = make(map[string]Interval)

// 自动初始化
func init() {
	for _, v := range IntervalLists {
		IntervalKeyMap[v.Key] = v
		IntervalKeyValMap[v.Val] = v
		IntervalKeyLabelMap[v.Label] = v
	}
}

const (
	DbClientGorm = "gorm"
	DbClientEs   = "elastic"
)

// 关系数据库
const (
	DbTypeMysql      = "mysql"
	DbTypePostgres   = "postgres"
	DbTypePostgresql = "postgresql"
	DbTypePg         = "pg"
	DbTypeSqlite     = "sqlite"
	DbTypeMssql      = "mssql"
	DbTypeOracle     = "oracle"
	DbTypeSqlserver  = "sqlserver"
	DbTypeTiDB       = "tidb"
)

// kv数据库
const (
	DbTypeMongo = "mongo"
	DbTypeRedis = "redis"
)

// 消息队列
const (
	QueueTypeName     = "消息队列"
	QueueTypeKafka    = "kafka" // kafka消息队列
	QueueTypeRabbit   = "rabbit"
	QueueTypeRocketmq = "rocketmq"
	QueueTypeRabbitmq = "rabbitmq"
)

const (
	QueueKafkaProducer = "kafka_producer" // kafka生产者
	QueueKafkaConsumer = "kafka_consumer" // kafka消费者
)

const (
	QueueRoleAsync = "async" // 异步
	QueueRoleSync  = "sync"  // 同步
)

// 文件存储
const (
	FileTypeName  = "文件"
	FileTypeFtp   = "ftp"
	FileTypeSftp  = "sftp"
	FileTypeLocal = "local"
	FileTypeOss   = "oss"
	FileTypeMinio = "minio"
	FileTypeHdfs  = "hdfs"
)

// 集群类型
const (
	ClusterEtcd      = "etcd"
	ClusterNacos     = "nacos"
	ClusterZookeeper = "zookeeper"
)

const (
	ClientInfoHost   = "host"
	ClientInfoUser   = "user"
	ClientInfoPasswd = "passwd"
)

const (
	ExportFileTypeExcel = "excel"
	ExportFileTypeCsv   = "csv"
)

const (
	ProgramLangJAVA   = "java"
	ProgramLangPHP    = "php"
	ProgramLangPYTHON = "python"
	ProgramLangGOLANG = "golang"
	ProgramLangLUA    = "lua"
)

const (
	ProtocolTCP = "tcp"
	ProtocolUDP = "udp"
)

type SortType int

func (s SortType) String() string {
	switch s {
	case SortAsc:
		return "asc"
	case SortDesc:
		return "desc"
	}
	return ""
}

const (
	SortAsc  SortType = iota // 升序
	SortDesc                 // 降序
)
