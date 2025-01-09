package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/oracledialect"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/enum"
	"workflow/internal/model"
	"workflow/internal/svc"
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
	// 统计加载成功和失败的数量
	successCount := 0
	failCount := 0
	skipCount := 0
	for _, v := range datasource {
		// 跳过fileServer
		if v.Type == enum.FileServerType.String() {
			skipCount++
			logx.Infof("datasource init skip: %d, %s", v.Id, v.Type)
			continue
		}
		// 读取dsn
		//dsn := gjson.Get(v.Config, "dsn").String()
		err := pool.UpdateDataSource(v.Id, v.Config, v.Type, v.Hash)
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

	for _, v := range datasource {
		if _, ok := pool.hash[v.Id]; ok {
			successCount++
		} else {
			failCount++
		}
	}
	logx.Infof("datasource init success: %d, failed: %d, skip: %d", successCount, failCount, skipCount)

	DataSourcePool = pool
}

func (manager *DataSourceManager) addDataSource(id int64, config string, dbType enum.DBType) error {
	// 根据提供的 dbType 创建数据库连接
	dsn := GenDataSourceDSN(dbType, config)
	sqlDB, err := sql.Open(dbType.String(), dsn)
	if err != nil {
		return err
	}
	// 设置连接超时时间为3秒
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 测试连接
	err = sqlDB.PingContext(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("database connection timeout (3 seconds)")
		}
		return err
	}
	// 中型应用
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)
	var bunDB *bun.DB
	switch dbType {
	case enum.MysqlType:
		bunDB = bun.NewDB(sqlDB, mysqldialect.New())
	case enum.SqlServerType:
		bunDB = bun.NewDB(sqlDB, mssqldialect.New())
	case enum.OracleType:
		bunDB = bun.NewDB(sqlDB, oracledialect.New())
	default:
		return errors.New("unsupported database type")
	}

	manager.dbs[id] = bunDB

	return nil
}

// UpdateDataSource 更新数据源连接
func (manager *DataSourceManager) UpdateDataSource(id int64, config, dbType, hash string) error {
	// 检查hash是否变化
	if oldHash, exists := manager.hash[id]; exists && oldHash == hash {
		return nil // hash未变化,无需更新
	}

	// 关闭旧连接
	if oldDB, exists := manager.dbs[id]; exists {
		if err := oldDB.Close(); err != nil {
			logx.Errorf("close old datasource failed: %v", err)
			return err
		}
		delete(manager.dbs, id)
	}

	// 创建新连接
	if err := manager.addDataSource(id, config, enum.DBType(dbType)); err != nil {
		logx.Errorf("create new datasource failed: %v", err)
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

func (manager *DataSourceManager) Insert(id int64, sql string, args ...interface{}) (sql.Result, error) {
	db, ok := manager.dbs[id]
	if !ok {
		return nil, errors.New("data source not found")
	}
	return db.ExecContext(context.Background(), sql, args...)
}

func (manager *DataSourceManager) Update(id int64, sql string, args ...interface{}) (sql.Result, error) {
	db, ok := manager.dbs[id]
	if !ok {
		return nil, errors.New("data source not found")
	}
	return db.ExecContext(context.Background(), sql, args...)
}

func (manager *DataSourceManager) Delete(id int64, sql string, args ...interface{}) (sql.Result, error) {
	db, ok := manager.dbs[id]
	if !ok {
		return nil, errors.New("data source not found")
	}
	return db.ExecContext(context.Background(), sql, args...)
}
