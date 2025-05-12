package enum

type DBType string

const (
	MysqlType      DBType = "mysql"
	SqlServerType  DBType = "sqlserver"
	OracleType     DBType = "oracle"
	FileServerType DBType = "fileServer"
	ModelType      DBType = "model"
)

func (t DBType) String() string {
	return string(t)
}
