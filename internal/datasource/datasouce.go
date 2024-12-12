package datasource

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/model"
	"workflow/internal/svc"
)

const (
	MysqlType     = "mysql"
	SqlServerType = "sqlserver"
)

type DataSourceManager struct {
	dbs  map[int64]*bun.DB
	hash map[int64]string // 存储每个数据源的hash值
}

var DataSourcePool *DataSourceManager

func (manager *DataSourceManager) GetHash() map[int64]string {
	return manager.hash
}

func InitDataSourceManager(svcCtx *svc.ServiceContext) {
	pool := &DataSourceManager{
		dbs:  make(map[int64]*bun.DB),
		hash: make(map[int64]string),
	}
	// 加载 datasource
	datasource, err := svcCtx.DatasourceModel.FindBySwitch(context.Background(), model.DatasourceSwitchOn)
	if err != nil {
		panic(err)
	}
	for _, v := range datasource {
		// 读取dsn
		dsn := gjson.Get(v.Config, "dsn").String()
		err := pool.UpdateDataSource(v.Id, dsn, v.Type, v.Hash)
		logx.Infof("datasource init: %+v", v)
		if err != nil {
			logx.Errorf("datasource init failed: %s", err.Error())
			continue
		}
		// 更新数据源状态
		err = svcCtx.DatasourceModel.UpdateStatus(context.Background(), v.Id, model.DatasourceStatusConnected)
		if err != nil {
			logx.Errorf("datasource update status failed: %s", err.Error())
			continue
		}
	}
	// 统计加载成功和失败的数量
	successCount := 0
	failCount := 0
	for _, v := range datasource {
		if _, ok := pool.hash[v.Id]; ok {
			successCount++
		} else {
			failCount++
		}
	}
	logx.Infof("datasource init success: %d, failed: %d", successCount, failCount)

	DataSourcePool = pool
}

func (manager *DataSourceManager) addDataSource(id int64, dsn, dbType string) error {
	// 根据提供的 dbType 创建数据库连接
	sqlDB, err := sql.Open(dbType, dsn)
	if err != nil {
		return err
	}
	// 测试连接
	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	var bunDB *bun.DB
	switch dbType {
	case MysqlType:
		bunDB = bun.NewDB(sqlDB, mysqldialect.New())
	case SqlServerType:
		bunDB = bun.NewDB(sqlDB, mssqldialect.New())
	default:
		return errors.New("unsupported database type")
	}

	manager.dbs[id] = bunDB

	return nil
}

// UpdateDataSource 更新数据源连接
func (manager *DataSourceManager) UpdateDataSource(id int64, dsn, dbType, hash string) error {
	// 检查hash是否变化
	if oldHash, exists := manager.hash[id]; exists && oldHash == hash {
		return nil // hash未变化,无需更新
	}

	// 关闭旧连接
	if oldDB, exists := manager.dbs[id]; exists {
		if err := oldDB.Close(); err != nil {
			return err
		}
		delete(manager.dbs, id)
	}

	// 创建新连接
	if err := manager.addDataSource(id, dsn, dbType); err != nil {
		return err
	}

	// 更新hash
	manager.hash[id] = hash
	return nil
}

// 清理链接
func (manager *DataSourceManager) ClearDataSource(id int64) error {
	if oldDB, exists := manager.dbs[id]; exists {
		err := oldDB.Close()
		logx.Infof("clear datasource: %d, err: %v", id, err)
		delete(manager.dbs, id)
		delete(manager.hash, id)
	}
	return nil
}

func (manager *DataSourceManager) Query(id int64, sql string, args ...interface{}) (*sql.Rows, error) {
	db, ok := manager.dbs[id]
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
