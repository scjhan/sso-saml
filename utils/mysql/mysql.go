package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var gMySQLUser string
var gMySQLPasswd string
var gMySQLName string

// InitMySQL init the mysql configuration
func InitMySQL(user string, passwd string, name string) {
	gMySQLName = name
	gMySQLPasswd = passwd
	gMySQLUser = user
}

// DbNode the database node value wrapper
type DbNode struct {
	value sql.RawBytes
}

// ToInt conver DbNode struct to int
func (thiz *DbNode) ToInt() int {
	iv, err := strconv.Atoi(string(thiz.value))
	if err != nil {
		panic("can not conver DbNode.value (type sql.RawBytes) to int")
	}
	return iv
}

// ToString Conver DbNode struct to string
func (thiz *DbNode) ToString() string {
	return string(thiz.value)
}

// MySQL the mysql struct
type MySQL struct {
	db *sql.DB
}

// Close close mysql
func (thiz *MySQL) Close() {
	thiz.db.Close()
}

// Query execute mysql query return value splic and it's length
func (thiz *MySQL) Query(query string) (map[string][]DbNode, int) {
	if thiz.db == nil {
		panic("nil MySQL instance")
	}

	rows, err := thiz.db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	cols, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	retval := make(map[string][]DbNode, len(cols))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	length := 0
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			panic(err.Error())
		}
		for i, val := range values {
			retval[cols[i]] = append(retval[cols[i]], DbNode{val})
		}
		length++
	}

	return retval, length
}

// Exec exectue mysql update insert delete .etc
func (thiz *MySQL) Exec(query string) error {
	if thiz.db == nil {
		return errors.New("nil MySQL instance")
	}

	if _, err := thiz.db.Exec(query); err != nil {
		return err
	}

	return nil
}

// NewMySQL return a MySQL instance
func NewMySQL() (MySQL, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", gMySQLUser, gMySQLPasswd, gMySQLName))
	if err != nil {
		return MySQL{nil}, err
	}

	db.Exec(fmt.Sprintf("use %s", gMySQLName))

	return MySQL{db}, err
}
