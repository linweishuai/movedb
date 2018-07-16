package sqlmaker

import (
	"bytes"
)

type Sqlmaker struct {
	Tablename	string
	Field	[]string
}

type Selectsqlmaker struct {
	Sqlmaker Sqlmaker
}

func (this Selectsqlmaker) SelectSqlmaker() string {
	sqlQuery := "select"
	for _, value := range this.Sqlmaker.Field {
		sqlQuery += (" `" + value + "` ,")
	}
	sqlQueryByte := []byte(sqlQuery)
	sqlQuery = string(sqlQueryByte[:len(sqlQueryByte)-1])
	sqlQuery+=(" from "+this.Sqlmaker.Tablename)
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

