package xmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/doug-martin/goqu/v9/exp"
)

// GetInsertSQLFromObj 获取插入语句
func GetInsertSQLFromObj(tableName string, o interface{}) string {
	return GetAddSQLFromObj(tableName, o, "INSERT INTO")
}

func GetReplaceInsertSQLFromObj(tableName string, o interface{}) string {
	return GetAddSQLFromObj(tableName, o, "REPLACE INTO")
}

func GetAddSQLFromObj(tableName string, o interface{}, addType string) string {
	columns := []string{}
	values := []string{}

	t := reflect.TypeOf(o)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := field.Tag.Get("db")
		if len(value) > 0 {
			columns = append(columns, "`"+value+"`")
			values = append(values, ":"+value)
		}
	}
	return fmt.Sprintf("%s `%s`(%s) VALUES (%s)", addType, tableName, strings.Join(columns, ","), strings.Join(values, ","))
}

// GetColumnsFromObj 获取所有字段
func GetColumnsFromObj(o interface{}) string {
	columns := GetColumnsArrayFromObj(o)
	return strings.Join(columns, ",")
}

// GetColumnsFromObj2 获取所有字段数组
func GetColumnsArrayFromObj(o interface{}) []string {
	columns := []string{}

	t := reflect.TypeOf(o)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := field.Tag.Get("db")
		if len(value) > 0 {
			columns = append(columns, "`"+value+"`")
		}
	}

	return columns
}

// GetColumnsFromObj2 获取所有字段数组
func GetColumnsArrayInterfaceFromObj(o interface{}) []interface{} {
	columns := make([]interface{}, 0)

	t := reflect.TypeOf(o)
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := field.Tag.Get("db")
		if len(value) > 0 {
			columns = append(columns, value)
		}
	}

	return columns
}

// GetColumnListExpressionFromObj 根据对象获取goqu的字段表达式
func GetColumnListExpressionFromObj(o interface{}) exp.ColumnListExpression {
	columns := GetColumnsArrayInterfaceFromObj(o)
	return exp.NewColumnListExpression(columns...)
}

// 判断乐观锁是否更新成功
func HasOptLockUpdateSuccess(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		return errors.New("update failed")
	}
	return nil
}

// 过滤sql.ErrNoRows错误
func ErrNoRowsFilter(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// 获取查询语句
func GetSqlSelect(o interface{}, table string) string {
	return "SELECT " + GetColumnsFromObj(o) + " FROM " + table + " "
}
