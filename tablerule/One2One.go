package tablerule

import (
	"movedb/dbconfig"
	"database/sql"
)

type One2One struct {
	Exportdb dbconfig.DbConfig
	Importdb dbconfig.DbConfig
	ExportField []string
	ImportField []string
}

func (this One2One)SetExportdb() []*sql.DB {
	exportdb:=make([]*sql.DB,0)
	exportdbpointer:=this.Exportdb.GetDbInstance()
	exportdb=append(exportdb,exportdbpointer)
	return exportdb
}
func (this One2One)SetImportdb() []*sql.DB {
	importdb:=make([]*sql.DB,0)
	importdbpointer:=this.Importdb.GetDbInstance()
	importdb=append(importdb,importdbpointer)
	return importdb
}