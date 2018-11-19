package importer

import (
	"database/sql"
	"movedb/sqlmaker"
	"github.com/cihub/seelog"
	"fmt"
	"movedb/ws"
)

type Importer struct {
	Dbconnction *sql.DB
	Insertsqlmaker 	sqlmaker.Insertsqlmaker
}

func (this Importer)Import(processchan chan struct{}) {
		selectsql := this.Insertsqlmaker.InsertSqlmaker()
		//fmt.Println(selectsql)
		_, err := this.Dbconnction.Exec(selectsql)
		if err != nil {
			ws.ReceiveApplication(fmt.Sprintf("%s插入出错%v",this.Insertsqlmaker.Sqlmaker.Tablename,err))
				return
			}
		seelog.Infof("%s导入数据成功",this.Insertsqlmaker.Sqlmaker.Tablename)
		<-processchan
}