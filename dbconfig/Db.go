package dbconfig

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Host string
	Username string
	Passwd string
	Dbname string
}

func (this *DbConfig)GetDbInstance() (*sql.DB,error)  {
	db, err := sql.Open("mysql", this.Username+":"+this.Passwd+"@tcp("+this.Host+":3306)/"+this.Dbname)
	if err != nil {
		return nil,err
	}
	return db,nil
}

