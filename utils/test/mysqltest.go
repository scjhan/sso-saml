package main

import (
	"chenjunhan/sso-saml/utils/mysql"
	"fmt"
)

func init() {
	mysql.InitMySQL("root", "123456", "saml")
}

func main() {
	o, err := mysql.NewMySQL()
	if err != nil {
		fmt.Println("NewMySQL error, error = %s", err.Error())
		return
	}
	defer o.Close()

	if err := o.Exec("update test set name='two' where id=1"); err != nil {
		fmt.Printf("mysql execute error, err = %s", err.Error())
	} else {
		fmt.Println("execute success")
	}

	ret, num := o.Query("select name from test where id=1")
	for i := 0; i < num; i++ {
		fmt.Printf("name[%d] = %s", i, ret["name"][i].ToString())
	}
}
