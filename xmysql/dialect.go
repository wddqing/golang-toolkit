package xmysql

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var mysqlDialect = goqu.Dialect("mysql")

func MysqlDialect() *goqu.DialectWrapper {
	return &mysqlDialect
}
