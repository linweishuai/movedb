package controller

import (
	"fmt"
	"io/ioutil"
	"github.com/cihub/seelog"
	"movedb/config"
	"encoding/json"
	"movedb/exporter"
	"sync"
	"time"
	"math"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
	"movedb/tranfer"
	"movedb/sqlmaker"
	"movedb/importer"
	"movedb/ws"
)

func Doexport()  {
	//config_path := flag.String("c","./config.json","config_file path")
	//goruntimeNumber:=flag.Float64("n",5000.00,"每个导入携程处理数量")
	//flag.Parse()
	//fmt.Printf("配置文件位置 : %s\n",*config_path)
	config_path:="./config.json"
	goruntimeNumber:=5000.00
	content,err:=ioutil.ReadFile(config_path)
	if err!=nil{
		seelog.Infof("读取配置文件出错")
		ws.ReceiveApplication("读取配置文件出错")
	}

	var Allconfig config.Allconfig;
	var config config.Conifg
	jsonerr := json.Unmarshal(content, &config)
	if jsonerr != nil {
		seelog.Infof("error in translating,", err.Error())
		ws.ReceiveApplication("error in translating,"+err.Error())
		return
	}
	Allconfig.Initrule(config)
	Allconfig.InitexportdbInstance(config)
	var Transferchan=make(chan exporter.Tranferdata,len(config.ExportDb))
	go func() {
		var wg sync.WaitGroup
		for Tname,Export:=range Allconfig.Exporter{
			wg.Add(1)
			go Export.Export(Transferchan,Tname,&wg,config.ExportDb[Tname]["tablename"])
		}
		wg.Wait()
		close(Transferchan)
		seelog.Infof("所有表数据到处完成")
		ws.ReceiveApplication("所有表数据到处完成")
	}()
	//从导出的数据入手 导出的数据多有几个那么就遍历这些数据
	beginTime := time.Now().Unix()
	for transferdata:=range Transferchan {
		//每次导入5000数据
		key:=transferdata.TableName
		value:=transferdata.Data
		ws.ReceiveApplication(fmt.Sprintf("处理导出%s导出数据共%d条数据",config.ExportDb[key]["tablename"],len(value)))
		seelog.Infof("处理导出%s导出数据共%d条数据",config.ExportDb[key]["tablename"],len(value))
		var ig sync.WaitGroup
		ProcessChan := make(chan struct{}, Allconfig.Cpunumber*2)//导入进程是cpu*2
		goroutineNumber:=goruntimeNumber
		GoroutineNumber := int(goroutineNumber)
		goNumber := math.Ceil(float64(len(value)) / goroutineNumber)
		for i := 0; i < int(goNumber); i++ {
			start := i * GoroutineNumber
			end := (i + 1) * GoroutineNumber
			var tempSlice []map[string]string
			if i == int(goNumber)-1 {
				tempSlice = value[start:]
			} else {
				tempSlice = value[start:end]
			}
			//fmt.Println(tempSlice)
			//os.Exit(1)
			Importslice:=make(map[string][]map[string]string)
			ws.ReceiveApplication(fmt.Sprintf("处理导出%s导出数 据第%d到%d数据",config.ExportDb[key]["tablename"],start, len(tempSlice)))
			seelog.Infof("处理导出%s导出数 据第%d到%d数据",config.ExportDb[key]["tablename"],start, len(tempSlice))
			vm := otto.New()
			for _,values:=range tempSlice{
				//seelog.Infof(values)
				for _,Importable:=range Allconfig.ExImRelation[key]{
					//如果有条件限制 就取出条件限制
					//example $account_type eq 'money' and id egt 5000
					//transfer $account_type == 'money' && id >= 5000
					//replace 'coupon'=='money' && '20000'>=5000 `javascript`
					if conditon,ok:=Allconfig.Tablerule[key+Importable];ok{
						var reg=regexp.MustCompile("\\$\\w+")
						javascript:=string(reg.ReplaceAllFunc([]byte(conditon), func(i []byte) []byte {
							return []byte(fmt.Sprintf("%q",values[strings.Replace(string(i),"$","",-1)]))
						}))
						jsres,_:=vm.Run(javascript)
						res,_:=jsres.ToBoolean()
						if !res{
							continue;
						}
					}
					var rowdata=make(map[string]string)
					//seelog.Infof(key.(string)+Importable)
					for _,exportfield:= range Allconfig.ExportTableField[key]  {
						//seelog.Infof(exportfield)
						for _,importfield:=range Allconfig.ImportTableField[Importable]{
							//seelog.Infof(importfield)
							//seelog.Infof(Importable+"."+importfield+":"+key.(string)+"."+exportfield)
							rule,ok:=Allconfig.Fieldrule[Importable+"."+importfield+":"+key+"."+exportfield]
							if ok{
								rowdata[importfield]=tranfer.DoTransfer(rule.TransferRule,exportfield,values,rule)
							}else{
								continue
							}
						}
					}
					Importslice[Importable]=append(Importslice[Importable],rowdata)
				}
			}
			//fmt.Println(Importslice)
			//fmt.Println(ImportTableFieldRelation)
			//os.Exit(1)
			for tablealias,importtempslice:=range Importslice{
				ig.Add(1)
				var insertsqlmaker=sqlmaker.Insertsqlmaker{
					Sqlmaker:sqlmaker.Sqlmaker{
						Tablename:config.ImportDb[tablealias]["tablename"],
						Field:Allconfig.ImportTableField[tablealias],
					},
					Dataslice:importtempslice,
				}
				var Importer=importer.Importer{
					Dbconnction:Allconfig.Importer[tablealias],
					Insertsqlmaker:insertsqlmaker,
				}
				ProcessChan <- struct{}{}
				go func(i int) {
					ws.ReceiveApplication(fmt.Sprintf("导入表%s进程%v开启",config.ImportDb[tablealias]["tablename"], i+1))
					seelog.Infof("导入表%s进程%v开启",config.ImportDb[tablealias]["tablename"], i+1)
					Importer.Import(ProcessChan)
					defer ig.Done()
				}(i)
			}
			ws.ReceiveApplication(fmt.Sprintf("处理导出%s导出数据第%d到%d数据完成",config.ExportDb[key]["tablename"],start,end))
			seelog.Infof("处理导出%s导出数据第%d到%d数据完成",config.ExportDb[key]["tablename"],start,end)
		}
		ig.Wait()
	}
	finishTime := time.Now().Unix()
	defer func() {
		ws.ReceiveApplication(fmt.Sprintf("导入完成实际消耗时间为：%v秒", finishTime-beginTime))
		seelog.Infof("导入完成实际消耗时间为：%v秒", finishTime-beginTime)
		seelog.Flush()
	}()
}
