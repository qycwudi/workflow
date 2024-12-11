package datasource

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type DataSourceManager struct {
	dbs  map[string]*bun.DB
	hash map[string]string // 存储每个数据源的hash值
}

func NewDataSourceManager() *DataSourceManager {
	return &DataSourceManager{
		dbs:  make(map[string]*bun.DB),
		hash: make(map[string]string),
	}
}

func (manager *DataSourceManager) AddDataSource(name, dsn, dbType string) error {
	// 根据提供的 dbType 创建数据库连接
	sqlDB, err := sql.Open(dbType, dsn)
	if err != nil {
		return err
	}

	var bunDB *bun.DB
	switch dbType {
	case "mysql":
		bunDB = bun.NewDB(sqlDB, mysqldialect.New())
	case "sqlserver":
		bunDB = bun.NewDB(sqlDB, mssqldialect.New())
	default:
		return errors.New("unsupported database type")
	}

	manager.dbs[name] = bunDB

	return nil
}

// UpdateDataSource 更新数据源连接
func (manager *DataSourceManager) UpdateDataSource(name, dsn, dbType, hash string) error {
	// 检查hash是否变化
	if oldHash, exists := manager.hash[name]; exists && oldHash == hash {
		return nil // hash未变化,无需更新
	}

	// 关闭旧连接
	if oldDB, exists := manager.dbs[name]; exists {
		if err := oldDB.Close(); err != nil {
			return err
		}
		delete(manager.dbs, name)
	}

	// 创建新连接
	if err := manager.AddDataSource(name, dsn, dbType); err != nil {
		return err
	}

	// 更新hash
	manager.hash[name] = hash
	return nil
}

func (manager *DataSourceManager) Query(name string, sql string, args ...interface{}) (*sql.Rows, error) {
	db, ok := manager.dbs[name]
	if !ok {
		return nil, errors.New("data source not found")
	}

	return db.QueryContext(context.Background(), sql, args...)
}

// func main() {
// 	manager := NewDataSourceManager()

// 	// 添加 MySQL 数据源
// 	err := manager.AddDataSource("mysql_db", "root:root@tcp(192.168.49.2:31426)/wk?charset=utf8mb4&parseTime=True&loc=Local", "mysql")
// 	if err != nil {
// 		log.Fatalf("failed to add data source: %v", err)
// 	}

// 	// // 添加 SQL Server 数据源
// 	// err = manager.AddDataSource("sqlserver_db", "sqlserver://username:password@localhost:1433?database=dbname")
// 	// if err != nil {
// 	// 	log.Fatalf("failed to add data source: %v", err)
// 	// }

// 	// 查询
// 	rows, err := manager.Query("mysql_db", "SELECT module_name as mn,module_type as mt FROM module")
// 	if err != nil {
// 		log.Fatalf("query failed: %v", err)
// 	}
// 	defer rows.Close()

// 	// 获取列信息
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		log.Fatalf("get columns failed: %v", err)
// 	}

// 	// 处理查询结果
// 	for rows.Next() {
// 		// 创建一个切片来存储所有列的值
// 		values := make([]interface{}, len(columns))
// 		// 创建一个切片来存储每列值的指针
// 		scanArgs := make([]interface{}, len(columns))
// 		for i := range values {
// 			scanArgs[i] = &values[i]
// 		}

// 		// 扫描当前行的数据到values切片中
// 		err := rows.Scan(scanArgs...)
// 		if err != nil {
// 			log.Fatalf("scan failed: %v", err)
// 		}

// 		// 打印每一列的值
// 		for i, col := range columns {
// 			val := values[i]
// 			// 处理null值
// 			if val == nil {
// 				fmt.Printf("%s: NULL\n", col)
// 			} else {
// 				// 将字节数组转换为字符串
// 				if b, ok := val.([]byte); ok {
// 					fmt.Printf("%s: %s\n", col, string(b))
// 				} else {
// 					fmt.Printf("%s: %v\n", col, val)
// 				}
// 			}
// 		}
// 		fmt.Println("-------------------")
// 	}
// }
