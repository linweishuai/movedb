package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"movedb/dbconfig"
	"movedb/exporter"
	"movedb/sqlmaker"
	"sync"
)

func main() {
	config_path := flag.String("c","/home/manu/sample/","config_file path")
	flag.Parse()
	fmt.Printf("配置文件位置 : %s\n",*config_path)
	content,err:=ioutil.ReadFile(*config_path)
	 if err!=nil{
		fmt.Println("读取配置文件出错")
	}
	type Conifg struct {
		ExportDb map[string]map[string]string
		ImportDb map[string]map[string]string
		Fieldrule map[string][]string
	}
	var inter Conifg
	jsonerr := json.Unmarshal(content, &inter)
	if jsonerr != nil {
		fmt.Println("error in translating,", err.Error())
		return
	}
	//解析出要导出的所有字段
	var exportField=make(map[string][]string)
	var ExportResult sync.Map
	for _,exportsource:=range inter.Fieldrule{
			//TODO importtarget not ready
			tablefield:=exportsource[0]
			tableFieldslice:=make([]string,0)
			tableFieldslice=strings.Split(tablefield,".")
		    exportField[tableFieldslice[0]]=append(exportField[tableFieldslice[0]],tableFieldslice[1])
	}
	var wg sync.WaitGroup
	fmt.Println(exportField)
	for tableName,FieldSlice:=range exportField{
		wg.Add(1)
		exportDB:=dbconfig.DbConfig{
			Host:inter.ExportDb[tableName]["host"],
			Username:inter.ExportDb[tableName]["username"],
			Passwd:inter.ExportDb[tableName]["passwd"],
			Dbname:inter.ExportDb[tableName]["dbname"],
		}
		var Dbconnection=exportDB.GetDbInstance();
		Exporter:=exporter.Exporter{
			Dbconnction:Dbconnection,
			Selectsqlmaker:sqlmaker.Selectsqlmaker{
				Sqlmaker:sqlmaker.Sqlmaker{
					Tablename:inter.ExportDb[tableName]["tablename"],
					Field:FieldSlice,
				},
				From:0,
				End:0,
			},
		}
		go Exporter.Export(&wg,&ExportResult,tableName)
	}
	wg.Wait()
	fmt.Println(ExportResult.Load("a"))
}
