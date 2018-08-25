package config

import (
	"runtime"
	"strings"
	"movedb/exporter"
	"movedb/dbconfig"
	"movedb/sqlmaker"
	"database/sql"
)


type Fieldrule struct {
	TransferRule string
	ExtraData [][]string
}

type Allconfig struct {
	ExportTableField map[string][]string //导出表字段
	ImportTableField map[string][]string //导入表字段
	ExImRelation map[string][]string //导出表和导入表的联系标注导出表的数据 导入到那个表
	Fieldrule map[string]Fieldrule //字段和字段之规则
	Tablerule map[string]string		//表和表之间的条件
	Exporter map[string]exporter.Exporter//导出数据库实例
	Importer map[string]*sql.DB//导出数据库实例
	GoruntimeNumber float64
	Cpunumber int
}
func (this *Allconfig) Initrule(config Conifg){
	this.Cpunumber=runtime.NumCPU()//导入进程做好是cpu的两倍
	var tabletempmap=make(map[string]struct{})//临时表
	var exportFieldtempmap=make(map[string]struct{})//唯一
	this.ExportTableField=make(map[string][]string)
	this.ImportTableField=make(map[string][]string)
	this.ExImRelation=make(map[string][]string)
	this.Fieldrule=make(map[string]Fieldrule)
	this.Tablerule=make(map[string]string)
	for Importtarget,exportsource:=range config.Fieldrule{
		//TODO importtarget not ready
		ImporttableFieldslice:=strings.Split(Importtarget,".")
		this.ImportTableField[ImporttableFieldslice[0]]=append(this.ImportTableField[ImporttableFieldslice[0]],ImporttableFieldslice[1])
		tablefield:=exportsource[0]
		ExportableFieldslice:=strings.Split(tablefield,".")
		if _,ok:=exportFieldtempmap[ExportableFieldslice[0]+ExportableFieldslice[1]];!ok{
			this.ExportTableField[ExportableFieldslice[0]]=append(this.ExportTableField[ExportableFieldslice[0]],ExportableFieldslice[1])
			exportFieldtempmap[ExportableFieldslice[0]+ExportableFieldslice[1]]= struct{}{}
		}
		if _,ok:=tabletempmap[ImporttableFieldslice[0]];!ok{
			this.ExImRelation[ExportableFieldslice[0]]=append(this.ExImRelation[ExportableFieldslice[0]],ImporttableFieldslice[0])
			tabletempmap[ImporttableFieldslice[0]]=struct{}{}
		}
		//如果是onetoone 模式 那么就需要 两个一一对应的参数
		var ExtraData [][]string
		if(exportsource[1]=="OnetoOne"){
			rulemap:=strings.Split(exportsource[2],":")//分隔出来然后
			for _,v:=range rulemap{
				ExtraData=append(ExtraData,strings.Split(v,","))
			}
		}
		this.Fieldrule[Importtarget+":"+tablefield]=Fieldrule{
			TransferRule:exportsource[1],
			ExtraData:ExtraData,
		}
	}
	//fmt.Println(ExImRelation)
	//os.Exit(1)
	//这是时候加上rulefield 里面的追加字段 并且去重
	for ruletablename,rulefieldslice:=range config.RuleField{
		this.ExportTableField[ruletablename]=append(this.ExportTableField[ruletablename],rulefieldslice...)
		this.ExportTableField[ruletablename]=RemoveRepByLoop(this.ExportTableField[ruletablename]);
	}
	//
	ConditionRegexp :=strings.NewReplacer("neq","!=","eq","==","egt",">=","gt",">","elt","<= ","lt","<")
	relationRegexp :=strings.NewReplacer("and","&&","or","||")
	for eitable,condtion:=range config.Tablerule{
		javascript:=relationRegexp.Replace(ConditionRegexp.Replace(condtion));
		this.Tablerule[eitable]=javascript
	}

}
func (this *Allconfig) InitexportdbInstance(config Conifg){
	this.Exporter=make(map[string]exporter.Exporter)
	this.Importer=make(map[string]*sql.DB)
	for tableName,FieldSlice:=range this.ExportTableField {
		exportDB := dbconfig.DbConfig{
			Host:     config.ExportDb[tableName]["host"],
			Username: config.ExportDb[tableName]["username"],
			Passwd:   config.ExportDb[tableName]["passwd"],
			Dbname:   config.ExportDb[tableName]["dbname"],
		}
		var Dbconnection,_ = exportDB.GetDbInstance();
		Exporter := exporter.Exporter{
			Dbconnction: Dbconnection,
			Selectsqlmaker: sqlmaker.Selectsqlmaker{
				Sqlmaker: sqlmaker.Sqlmaker{
					Tablename: config.ExportDb[tableName]["tablename"],
					Field:     FieldSlice,
				},
			},
		}
		this.Exporter[tableName]=Exporter
	}
	for tableName,_:=range this.ImportTableField {
		ImportDb := dbconfig.DbConfig{
			Host:     config.ImportDb[tableName]["host"],
			Username: config.ImportDb[tableName]["username"],
			Passwd:   config.ImportDb[tableName]["passwd"],
			Dbname:   config.ImportDb[tableName]["dbname"],
		}
		var Dbconnection,_ = ImportDb.GetDbInstance();
		this.Importer[tableName]=Dbconnection
	}
}
func RemoveRepByLoop(slc []string) []string {
	result := []string{}  // 存放结果
	for i := range slc{
		flag := true
		for j := range result{
			if slc[i] == result[j] {
				flag = false  // 存在重复元素，标识为false
				break
			}
		}
		if flag {  // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}