package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"movedb/dbconfig"//数据哭连接
	"movedb/exporter"//导出
	"movedb/sqlmaker"//sql生成
	"movedb/importer"//导入
	"sync"
	"math"
	"time"
	"github.com/cihub/seelog"
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
	//fmt.Print(inter.Fieldrule)
	var ImportTableFieldRelation=make(map[string][]string)//导入数据库的字段
	var ExImRelation=make(map[string][]string)//标注我导出的数据导入到那些表里面去
	var tempmap=make(map[string]string)//临时表
	type Fieldrule struct {
		TransferRule string
		ExtraData [][]string
	}
	var NewFeildmap=make(map[string]Fieldrule)
	for Importtarget,exportsource:=range inter.Fieldrule{
			//TODO importtarget not ready
			ImporttableFieldslice:=strings.Split(Importtarget,".")
			ImportTableFieldRelation[ImporttableFieldslice[0]]=append(ImportTableFieldRelation[ImporttableFieldslice[0]],ImporttableFieldslice[1])
			tablefield:=exportsource[0]
			ExportableFieldslice:=strings.Split(tablefield,".")
		    exportField[ExportableFieldslice[0]]=append(exportField[ExportableFieldslice[0]],ExportableFieldslice[1])
		    _,ok:=tempmap[ImporttableFieldslice[0]]
		    if!ok{//如果读取不到那么就说明不存在这导入表
				tempmap[ImporttableFieldslice[0]]=tempmap[ImporttableFieldslice[0]]
				ExImRelation[ExportableFieldslice[0]]=append(ExImRelation[ExportableFieldslice[0]],ImporttableFieldslice[0])
			}
			//如果是onetoone 模式 那么就需要 两个一一对应的参数
			var ExtraData [][]string
			if(exportsource[1]=="OnetoOne"){
				rulemap:=strings.Split(exportsource[2],":")//分隔出来然后
				for _,v:=range rulemap{
					ExtraData=append(ExtraData,strings.Split(v,","))
				}
			}
			NewFeildmap[Importtarget+":"+tablefield]=Fieldrule{
				TransferRule:exportsource[1],
				ExtraData:ExtraData,
			}
	}
	var wg sync.WaitGroup
	//fmt.Println(NewFeildmap)
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
	fmt.Printf("导出所有待用数据")
	//从导出的数据入手 导出的数据多有几个那么就遍历这些数据
	beginTime := time.Now().Unix()
	ExportResult.Range(func(key, value interface{}) bool {
		//每次导入5000数据
		fmt.Printf("处理导出%s导出数据共%d条数据",key.(string),len(value.([]map[string]string)))
		var ig sync.WaitGroup
		ProcessChan := make(chan struct{}, 5)
		goroutineNumber:=5000.00
		GoroutineNumber := int(goroutineNumber)
		goNumber := math.Ceil(float64(len(value.([]map[string]string))) / goroutineNumber)
		for i := 0; i < int(goNumber); i++ {
			start := i * GoroutineNumber
			end := (i + 1) * GoroutineNumber
			var tempSlice []map[string]string
			if i == int(goNumber)-1 {
				tempSlice = value.([]map[string]string)[start:]
			} else {
				tempSlice = value.([]map[string]string)[start:end]
			}
			Importslice:=make(map[string][][]string)
			fmt.Printf("处理导出%s导出数据第%d到%d数据",key.(string),start,end)
			for _,values:=range tempSlice{
				//fmt.Println(values)
				for _,Importable:=range ExImRelation[key.(string)]{
					rowdata:=make([]string,0)
				for _,exportfield:= range exportField[key.(string)]  {
					//fmt.Println(exportfield)
						for _,importfield:=range ImportTableFieldRelation[Importable]{
							//fmt.Println(Importable)
							//fmt.Println(Importable+"."+importfield+":"+key.(string)+"."+exportfield)
							rule,ok:=NewFeildmap[Importable+"."+importfield+":"+key.(string)+"."+exportfield]
							if ok{
								switch (rule.TransferRule){
								case "Default":
									rowdata=append(rowdata,fmt.Sprintf("%q", values[exportfield]))
								case "OnetoOne":
									index :=0
									for nowIndex,content:=range rule.ExtraData[0]{
										if(content==values[exportfield]){
											index=nowIndex
										}
									}
									rowdata=append(rowdata,fmt.Sprintf("%q", rule.ExtraData[1][index]))
								default:
									rowdata=append(rowdata,fmt.Sprintf("%q", values[exportfield]))
								}
							}else{
								continue
							}
						}
					}
					Importslice[Importable]=append(Importslice[Importable],rowdata)
				}
			}

			for tablealias,importtempslice:=range Importslice{
				ig.Add(1)
					var insertsqlmaker=sqlmaker.Insertsqlmaker{
							Sqlmaker:sqlmaker.Sqlmaker{
								Tablename:inter.ImportDb[tablealias]["tablename"],
								Field:ImportTableFieldRelation[tablealias],
							},
							Dataslice:importtempslice,
					}
					var ImportDbconfig=dbconfig.DbConfig{
						Host:inter.ImportDb[tablealias]["host"],
						Username:inter.ImportDb[tablealias]["username"],
						Passwd:inter.ImportDb[tablealias]["passwd"],
						Dbname:inter.ImportDb[tablealias]["dbname"],
					}
					var Importer=importer.Importer{
						Dbconnction:ImportDbconfig.GetDbInstance(),
						Insertsqlmaker:insertsqlmaker,
					}
					ProcessChan <- struct{}{}
					go func(i int) {
						//_ = tempSlice
						seelog.Infof("进程%v开启", i+1)
						Importer.Import(&ig,ProcessChan)
					}(i)
			}
		}
		ig.Wait()
		return true
	})
	finishTime := time.Now().Unix()
	seelog.Infof("实际消耗时间为：%v秒", finishTime-beginTime)





















}
