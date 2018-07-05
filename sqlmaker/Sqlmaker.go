package sqlmaker

import (
	"strconv"
	"bytes"
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
	if this.From!=0&&this.End!=0{
		sqlQuery += (" where id >=" + strconv.Itoa(this.From)+" and id <="+strconv.Itoa(this.End))
	}
	return sqlQuery
}
type Insertsqlmaker struct {
	Sqlmaker Sqlmaker
	Dataslice [][]string
}
func (this Insertsqlmaker) InsertSqlmaker() string {
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO " + this.Sqlmaker.Tablename + " (")

	for _, field := range this.Sqlmaker.Field {
		buffer.WriteString(" " + field + ",")
	}
	buffer.Truncate(buffer.Len() - 1)

	buffer.WriteString(")VALUES ")

	for _, item := range  this.Dataslice{
		buffer.WriteString("(")
		for _, value := range item {
			buffer.WriteString(value)
			buffer.WriteString(",")
		}
		buffer.Truncate(buffer.Len() - 1)

		buffer.WriteString("),")
	}
	buffer.Truncate(buffer.Len() - 1)
	return  buffer.String()
}

