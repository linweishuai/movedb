package tablerule

import (
	"database/sql"
)

type Tablerule interface {
	SetExportdb() []*sql.DB
	SetImportdb() []*sql.DB
}