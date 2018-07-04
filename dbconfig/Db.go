package dbconfig

import (
	"database/sql"
	"github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Host string
	Username string
	Passwd string
	Dbname string
}

func (this *DbConfig)GetDbInstance() *sql.DB  {
	db, err := sql.Open("mysql", this.Username+":"+this.Passwd+"@tcp("+this.Host+":3306)/"+this.Dbname)
	if err != nil {
		seelog.Errorf("打开数据库出错%v", err)
	}
	return db
}

