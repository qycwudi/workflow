package components

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/rotisserie/eris"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/datasource"
	"workflow/internal/enum"
	"workflow/pkg/core"
)

type DatabaseComponent struct {
	config DatabaseConfig
}
type DatabaseConfig struct {
	DatasourceId      int64           `json:"datasourceId"`
	SQL               string          `json:"sql"`
	ErrorHandlingMode string          `json:"errorHandlingMode"`
	Retry             int64           `json:"retry"`
	Timeout           int64           `json:"timeout"`
	DatasourceType    string          `json:"datasourceType"`
	ExceptionConfig   ExceptionConfig `json:"exceptionConfig"`
}

var databaseComponentPool = sync.Pool{
	New: func() interface{} {
		return &DatabaseComponent{}
	},
}

func NewDatabaseComponent(config any) (*DatabaseComponent, error) {
	jsonConfig, err := sonic.Marshal(config)
	if err != nil {
		return nil, eris.Wrap(err, "failed to parse database component config")
	}
	component := databaseComponentPool.Get().(*DatabaseComponent)
	var databaseConfig DatabaseConfig
	if err := sonic.Unmarshal(jsonConfig, &databaseConfig); err != nil {
		return nil, eris.Wrap(err, "failed to parse database component config")
	}
	databaseConfig.ExceptionConfig = ExceptionConfig{
		Timeout:    databaseConfig.Timeout,
		RetryTimes: int(databaseConfig.Retry),
		OutputOnError: map[string]any{
			"outputList": []map[string]any{},
			"rowNum":     0,
		},
	}
	component.config = databaseConfig
	return component, nil
}

func (d *DatabaseComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	inputMap, ok := input.(map[string]any)
	if !ok {
		return nil, eris.New("input type mismatch")
	}

	// 直接替换所有变量为 ?，并收集参数
	sql, args, err := d.replaceExprs(d.config.SQL, inputMap, func(expr, old string, value any) string {
		return strings.ReplaceAll(expr, old, "?")
	})
	if err != nil {
		return nil, err
	}

	// sql 替换完成，开始执行
	param := fmt.Sprintf("sql:%s,args:%+v", sql, args)
	logx.Debugf("[DATABASE] %s", param)
	output, err := d.executeSQL(ctx, sql, args)
	if err != nil {
		logx.Errorf("[DATABASE] execute failed param:%s,error:%s", param, err.Error())
		return &core.Result{
			Route:  []string{Failed},
			Output: d.config.ExceptionConfig.OutputOnError,
		}, eris.Wrap(err, "failed to replace exprs param:"+param)
	}
	return &core.Result{
		Route:  []string{Success},
		Output: output,
	}, nil
}

func (d *DatabaseComponent) Validate() []core.ValidationError {
	return nil
}

func (d *DatabaseComponent) AnalyzeInputs(ctx context.Context) (any, error) {
	return nil, nil
}

func (d *DatabaseComponent) Exception() ExceptionConfig {
	return d.config.ExceptionConfig
}

func (d *DatabaseComponent) Clear() {
	endItemComponentPool.Put(d)
}

func (d *DatabaseComponent) replaceExprs(expr string, inputMap map[string]any, replaceFun func(expr, old string, value any) string) (string, []any, error) {
	re := regexp.MustCompile(`{{.*?}}`)
	matches := re.FindAllString(expr, -1)

	var ok bool
	var value any
	var args []any
	for _, field := range matches {
		cutFiled := strings.Trim(field, "{{}}")
		if value, ok = inputMap[cutFiled]; !ok {
			// 本节点 input 里没有该变量
			return "", nil, eris.Errorf("node input does not exist variable [%s]", cutFiled)
		}
		if value == nil {
			logx.Errorf("[DATABASE] value is nil,cutFiled:%s", cutFiled)
			value = reflect.Zero(reflect.TypeOf(value)).Interface()
		}
		// 如果是数组,要转成 f1,f2,f3 这种格式
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			slice := reflect.ValueOf(value)
			for i := 0; i < slice.Len(); i++ {
				args = append(args, slice.Index(i).Interface())
			}
			expr = replaceFun(expr, field, args)
		} else {
			// 替换
			args = append(args, value)
			expr = replaceFun(expr, field, value)
		}
	}
	return expr, args, nil
}

// 执行SQL语句
func (d *DatabaseComponent) executeSQL(ctx context.Context, sql string, args []interface{}) (any, error) {
	if d.config.DatasourceType == enum.OracleType.String() {
		sql = strings.ReplaceAll(sql, ";", "")
	}

	sqlType := strings.ToUpper(strings.TrimSpace(sql))
	switch {
	case strings.HasPrefix(sqlType, "SELECT"):
		return d.executeQuery(ctx, sql, args)
	case strings.HasPrefix(sqlType, "INSERT"):
		return d.executeInsert(ctx, sql, args)
	case strings.HasPrefix(sqlType, "UPDATE"):
		return d.executeUpdate(ctx, sql, args)
	case strings.HasPrefix(sqlType, "DELETE"):
		return d.executeDelete(ctx, sql, args)
	default:
		return nil, fmt.Errorf("unsupported SQL type: %s", sqlType)
	}
}

func queryResult(rows *sql.Rows) ([]any, error) {
	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		logx.Errorf("get columns failed: %v", err)
		return nil, err
	}

	// 存储所有行的结果
	var result []any

	// 处理查询结果
	for rows.Next() {
		// 创建一个切片来存储所有列的值
		values := make([]interface{}, len(columns))
		// 创建一个切片来存储每列值的指针
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// 扫描当前行的数据到values切片中
		err := rows.Scan(scanArgs...)
		if err != nil {
			logx.Errorf("row scan failed: %v", err)
			return nil, err
		}

		// 将当前行数据转换为map
		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				rowData[col] = nil
			} else {
				// 将字节数组转换为字符串
				if b, ok := val.([]byte); ok {
					rowData[col] = string(b)
				} else {
					rowData[col] = val
				}
			}
		}
		result = append(result, rowData)
	}

	// 添加错误检查
	if err = rows.Err(); err != nil {
		logx.Errorf("row iteration failed: %v", err)
		return nil, err
	}
	// 将结果赋值给msg的数据部分
	// marshal, _ := json.Marshal(result)
	// return marshal, nil
	return result, nil
}

// 执行查询
func (d *DatabaseComponent) executeQuery(ctx context.Context, sql string, args []interface{}) (any, error) {
	rows, err := datasource.DataSourcePool.Query(ctx, d.config.DatasourceId, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := queryResult(rows)
	if err != nil {
		logx.Errorf("[DATABASE] query failed,SQL:%s,error:%s", sql, err.Error())
		return nil, err
	}
	return map[string]any{
		"outputList": result,
		"rowNum":     len(result),
	}, nil
}

// 执行插入
func (d *DatabaseComponent) executeInsert(ctx context.Context, sql string, args []interface{}) (any, error) {
	result, err := datasource.DataSourcePool.Insert(ctx, d.config.DatasourceId, sql, args...)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return map[string]any{
		"outputList": []int64{id},
		"rowNum":     1,
	}, nil
}

// 执行更新
func (d *DatabaseComponent) executeUpdate(ctx context.Context, sql string, args []interface{}) (any, error) {
	result, err := datasource.DataSourcePool.Update(ctx, d.config.DatasourceId, sql, args...)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	return map[string]any{
		"outputList": []int64{},
		"rowNum":     affected,
	}, nil
}

// 执行删除
func (d *DatabaseComponent) executeDelete(ctx context.Context, sql string, args []interface{}) (any, error) {
	result, err := datasource.DataSourcePool.Delete(ctx, d.config.DatasourceId, sql, args...)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	return map[string]any{
		"outputList": []int64{},
		"rowNum":     affected,
	}, nil
}
