package exporter

import (
	"database/sql"
	"movedb/sqlmaker"
	"github.com/cihub/seelog"
	"math"
	"strconv"
	"sync"
)

type Exporter struct {
	Dbconnction *sql.DB
	Selectsqlmaker 	sqlmaker.Selectsqlmaker
	Totalnum float64
}
type Tranferdata struct {
	TableName string
	Data []map[string]string
}

func (this Exporter)Export(exportResult chan <- Tranferdata,tableName string,wg *sync.WaitGroup,realName string) {
		this.getSum()//调用一下获取总数 然后开始循环查询 每次每个表 默认查1w条数据
		goexportnum := int(math.Ceil(float64(this.Totalnum / 50000)))
		selectsql := this.Selectsqlmaker.SelectSqlmaker()
		for i := 0; i <=goexportnum; i++ {
			var doSql=selectsql;
			start:=i*50000
			var result = make([]map[string]string, 0)
			doSql+=" limit "+strconv.Itoa(start)+", 50000"
			//fmt.Println(doSql)
			exportRs, err := this.Dbconnction.Query(doSql)
			if err != nil {
				seelog.Errorf("表%v查询出错%v",realName, err)
			}

			var ExportFieldSlice= make([]interface{}, 0, len(this.Selectsqlmaker.Sqlmaker.Field))
			copyExportField := Copy(this.Selectsqlmaker.Sqlmaker.Field)
			for _, key := range this.Selectsqlmaker.Sqlmaker.Field {
				ExportFieldSlice = append(ExportFieldSlice, copyExportField[key])
			}
			var exportcount = 0
			for exportRs.Next() {
				if err := exportRs.Scan(ExportFieldSlice...); err != nil {
					seelog.Errorf("读取数据出错%v", err)
				}
				var rowdata = make(map[string]string)
				for key, value := range this.Selectsqlmaker.Sqlmaker.Field {
					rowdata[value] = string(*(ExportFieldSlice[key].(*[]byte)))
				}
				//fmt.Println(rowdata)
				exportcount++
				result = append(result, rowdata)
			}
			//seelog.Info(result)
			exportRs.Close()
			//os.Exit(2)
			seelog.Infof("%s导出%d条数据 from %d", realName, exportcount,start)
			TranferData := Tranferdata{
				TableName: tableName,
				Data:      result,
			}
			if(len(result)!=0){
				exportResult<-TranferData//如果不为空就把数据放进去 为空的话 就开始表示这个表已经到处结束了
			}else{
				wg.Done()
				return
			}
		}
}
	func Copy(region []string) map[string]*[]byte {
		copy := make(map[string]*[]byte)

		for _, value := range region {
			copy[value] = new([]byte)
		}
		return copy
	}
	func (this *Exporter)getSum()  {
		var querycount="select count(*) from "+this.Selectsqlmaker.Sqlmaker.Tablename
		exportRs:=this.Dbconnction.QueryRow(querycount)
		err:=exportRs.Scan(&this.Totalnum)
		if err!=nil{
			seelog.Errorf("%v查询总数出错", this.Selectsqlmaker.Sqlmaker.Tablename)
		}
	}