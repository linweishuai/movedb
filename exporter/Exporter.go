package exporter

import (
	"database/sql"
	"movedb/sqlmaker"
	"github.com/cihub/seelog"
)

type Exporter struct {
	Dbconnction *sql.DB
	Selectsqlmaker 	sqlmaker.Selectsqlmaker
}
type Tranferdata struct {
	TableName string
	Data []map[string]string
}

func (this Exporter)Export(exportResult chan Tranferdata,tableName string) {
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

		exportcount:=0
		for exportRs.Next() {
			if err := exportRs.Scan(ExportFieldSlice...); err != nil {
				seelog.Errorf("读取数据出错%v", err)
			}
			var rowdata=make(map[string]string)
			for key,value:=range this.Selectsqlmaker.Sqlmaker.Field{
				rowdata[value]=string(*(ExportFieldSlice[key].(*[]byte)))
			}
			//fmt.Println(rowdata)
			exportcount++
			result = append(result, rowdata)
		}
		exportRs.Close()
		//os.Exit(2)
		seelog.Infof("%s导出%d条数据",tableName,exportcount)
		TranferData:=Tranferdata{
			TableName:tableName,
			Data:result,
		}
		exportResult<-TranferData
}
	func Copy(region []string) map[string]*[]byte {
		copy := make(map[string]*[]byte)

		for _, value := range region {
			copy[value] = new([]byte)
		}
		return copy
	}