package exporter

import (
	"database/sql"
	"movedb/sqlmaker"
	"sync"
	"github.com/cihub/seelog"
)

type Exporter struct {
	Dbconnction *sql.DB
	Selectsqlmaker 	sqlmaker.Selectsqlmaker
}

func (this Exporter)Export(wg sync.WaitGroup) []map[string]string {
		result:=make([]map[string]string,0)
		selectsql := this.Selectsqlmaker.SelectSqlmaker()
		exportRs, err := this.Dbconnction.Query(selectsql)
		if err != nil {
			seelog.Errorf("查询出错%v", err)
		}
		var ExportFieldSlice = make([]interface{}, 0, len(this.Selectsqlmaker.Sqlmaker.Field))
		copyExportField:=Copy(this.Selectsqlmaker.Sqlmaker.Field)
		for _, key := range this.Selectsqlmaker.Sqlmaker.Field {
			ExportFieldSlice = append(ExportFieldSlice, copyExportField[key])
		}
		for exportRs.Next() {
			if err := exportRs.Scan(ExportFieldSlice...); err != nil {
				seelog.Errorf("读取数据出错%v", err)
			}
			var rowdata=make(map[string]string)
			for key,value:=range this.Selectsqlmaker.Sqlmaker.Field{
				rowdata[value]=string(*(ExportFieldSlice[key].(*[]byte)))
			}
			result = append(result, rowdata)
		}
		wg.Done()
		return result
}
	func Copy(region []string) map[string]*[]byte {
		copy := make(map[string]*[]byte)

		for _, value := range region {
			copy[value] = new([]byte)
		}
		return copy
	}