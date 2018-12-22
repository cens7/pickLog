package ot

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	driverName = "mysql"
	user       = ""
	pass       = ""
	protocol   = ""
	ip         = ""
	port       = ""
	dbName     = ""
)

func QryServe(info *AppInfo) (se *ServeInfo) {

	if info.hostIp == "" {
		fmt.Println("登录获取对象错误")
		return &ServeInfo{}
	}
	
	serve := &ServeInfo{}

	path := strings.Join([]string{user, ":", pass, "@", protocol, "(",
		ip, ":", port, ")/", dbName, "?charset=utf8"}, "")

	db, e := sql.Open(driverName, path)
	if e != nil {
		panic(e.Error())
	}
	var sqlStr = "select * from i_server_key WHERE ip = '" + info.hostIp + "'"
	rows, err := db.Query(sqlStr)
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	serve.hostIp = info.hostIp
	serve.appName = info.name
	serve.username = string(values[2])
	serve.passoword = string(values[3])

	fmt.Println("\n---------->当前应用所在IP：",serve.hostIp,", 应用名：",serve.appName)

	return serve

}

//服务器信息
type ServeInfo struct {
	appName   string
	hostIp    string
	username  string
	passoword string
}
