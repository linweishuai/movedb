package importer

import (
	"database/sql"
	"movedb/sqlmaker"
	"github.com/cihub/seelog"
	"os"
)

type Importer struct {
	Dbconnction *sql.DB
	Insertsqlmaker 	sqlmaker.Insertsqlmaker
}

func (this Importer)Import(processchan chan struct{}) {
		selectsql := this.Insertsqlmaker.InsertSqlmaker()
		_, err := this.Dbconnction.Exec(selectsql)
		if err != nil {
			seelog.Errorf("插入出错%v", err)
			os.Exit(2)
			}
		seelog.Infof("%s导入数据成功",this.Insertsqlmaker.Sqlmaker.Tablename)
		<-processchan
}