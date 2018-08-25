package model

import (
	"movedb/dbconfig"
	"github.com/cihub/seelog"
)

type Db struct {
	dbconfig.DbConfig
	Tablename string
}

func (this Db)GetDbfield()  {
	dbconnction:=this.DbConfig.GetDbInstance()
	showfieldsql:="SHOW COLUMNS FROM `"+this.Tablename+"`"
	exportRs, err := dbconnction.Query(showfieldsql)
	if err != nil {
		seelog.Errorf("表%v查询出错%v",this.Tablename, err)
	}

	for exportRs.Next() {

	}
}