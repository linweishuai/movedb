package importer

import (
	"database/sql"
	"movedb/sqlmaker"
	"sync"
	"github.com/cihub/seelog"
)

type Importer struct {
	Dbconnction *sql.DB
	Insertsqlmaker 	sqlmaker.Insertsqlmaker
}

func (this Importer)Import(wg *sync.WaitGroup,processchan chan struct{}) {
		selectsql := this.Insertsqlmaker.InsertSqlmaker()
		_, err := this.Dbconnction.Exec(selectsql)
		if err != nil {
			seelog.Errorf("插入出错%v", err)
		}
		wg.Done()
		<-processchan
}