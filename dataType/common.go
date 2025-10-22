package dataType

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/helays/utils/v2/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func JsonDbDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case config.DbTypeSqlite:
		return "JSON"
	case config.DbTypeMysql:
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") && CheckVersionSupportsJSON(v.ServerVersion) {
			return "JSON"
		}
		return "LONGTEXT"
	case config.DbTypePostgres:
		return "JSONB"
	case config.DbTypeSqlserver:
		return "NVARCHAR(MAX)"
	}
	return "TEXT"
}

func DriverValueWithJson(val any) (driver.Value, error) {
	if val == nil {
		return nil, nil
	}

	b, err := json.Marshal(val)
	return string(b), err
}

// DriverScanWithJson 解析json
func DriverScanWithJson[T any](val any, dst *T) error {
	if val == nil {
		*dst = *new(T)
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	case json.RawMessage:
		ba = v
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", val))
	}
	if len(ba) < 1 {
		*dst = *new(T)
		return nil
	}

	// 如果需要用 UseNumber，就在dst 的实体类型上 实现 json.UnmarshalJSON
	if err := json.Unmarshal(ba, dst); err != nil {
		return fmt.Errorf("failed to unmarshal JSON value: %w", err)
	}
	return nil
}

// CheckVersionSupportsJSON 检查版本是否支持JSON
// mysql版本高于 5.7.8 ，才支持json
func CheckVersionSupportsJSON(versionStr string) bool {
	versionParts := strings.Split(strings.TrimSpace(strings.Split("versionStr", "-")[0]), ".")
	if len(versionParts) < 3 {
		return false
	}
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return false
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return false
	}

	patch, err := strconv.Atoi(versionParts[2])
	if err != nil {
		return false
	}
	return major > 5 || (major == 5 && minor > 7) || (major == 5 && minor == 7 && patch >= 8)
}
func marshalSlice(v any) ([]byte, error) {
	if v == nil || reflect.ValueOf(v).Len() < 1 {
		return []byte("[]"), nil
	}
	return json.Marshal(v)
}
func arrayValue(m any) (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := marshalSlice(m)

	return string(ba), err
}

func arrayScan(m any, val any) error {
	if val == nil {
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", val)
	}
	rd := bytes.NewReader(ba)
	decoder := json.NewDecoder(rd)
	decoder.UseNumber()
	return decoder.Decode(m)
}

func arrayGormValue(jm any, db *gorm.DB) clause.Expr {
	data, _ := marshalSlice(jm)
	switch db.Dialector.Name() {
	case config.DbTypeMysql:
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") && CheckVersionSupportsJSON(v.ServerVersion) {
			fmt.Println(v.ServerVersion)
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
