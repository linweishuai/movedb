package sqlmaker

import (
	"strconv"
	"bytes"
	"github.com/cihub/seelog"
)

type Sqlmaker struct {
	Tablename	string
	Field	[]string
}

type Selectsqlmaker struct {
	Sqlmaker Sqlmaker
	From int
	End int
}

func (this Selectsqlmaker) SelectSqlmaker() string {
	sqlQuery := "select"
	for _, value := range this.Sqlmaker.Field {
		sqlQuery += (" `" + value + "` ,")
	}
	sqlQueryByte := []byte(sqlQuery)
	sqlQuery = string(sqlQueryByte[:len(sqlQueryByte)-1])
	sqlQuery+=(" from "+this.Sqlmaker.Tablename)
	if this.End!=0{
		sqlQuery += (" limit " + strconv.Itoa(this.From)+" ,"+strconv.Itoa(this.End))
		seelog.Infof("准备导出%d到%d数据",this.From,this.End)
	}
	return sqlQuery
}
type Insertsqlmaker struct {
	Sqlmaker Sqlmaker
	Dataslice []map[string]string
}
func (this Insertsqlmaker) InsertSqlmaker() string {
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO " + this.Sqlmaker.Tablename + " (")

	for _, field := range this.Sqlmaker.Field {
		buffer.WriteString("`" + field + "`,")
	}
	buffer.Truncate(buffer.Len() - 1)

	buffer.WriteString(")VALUES ")
	//fmt.Println(this.Dataslice)
	//fmt.Println(this.Sqlmaker.Field)
	//os.Exit(1)
	for _, item := range  this.Dataslice{
		buffer.WriteString("(")
		for _, fieldname := range this.Sqlmaker.Field {
			buffer.WriteString(item[fieldname])
			buffer.WriteString(",")
		}
		buffer.Truncate(buffer.Len() - 1)

		buffer.WriteString("),")
	}
	buffer.Truncate(buffer.Len() - 1)
	return  buffer.String()
}

