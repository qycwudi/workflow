package components

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"workflow/internal/datasource"
	"workflow/pkg/core"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
)

type DatabaseComponent struct {
	config DatabaseConfig
}
type DatabaseConfig struct {
	DatasourceType  string      `json:"datasourceType"`
	DatasourceId    int64       `json:"datasourceId"`
	SQL             string      `json:"sql"`
	ExceptionConfig interface{} `json:"exceptionConfig"`
}

var databaseComponentPool = sync.Pool{
	New: func() interface{} {
		return &DatabaseComponent{}
	},
}

func NewDatabaseComponent(config json.RawMessage) (*DatabaseComponent, error) {
	// 使用pool
	component := databaseComponentPool.Get().(*DatabaseComponent)
	var databaseConfig DatabaseConfig
	if err := sonic.Unmarshal(config, &databaseConfig); err != nil {
		return nil, errors.New("解析迭代结束组件配置失败: " + err.Error())
	}
	component.config = databaseConfig
	return component, nil
}

func (d *DatabaseComponent) Execute(ctx context.Context, input any) (*core.Result, error) {
	inputMap, ok := input.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("input 类型不匹配")
	}

	var err error
	var args []any
	var condition string
	if splitSql := strings.Split(strings.ToLower(d.config.SQL), "where"); len(splitSql) == 2 {
		d.config.SQL = splitSql[0]
		condition = "where" + splitSql[1]
	}
	if condition != "" {
		if condition, args, err = d.replaceExprs(condition, inputMap, func(expr, old string, value any) string {
			return strings.ReplaceAll(expr, old, "?")
		}); err != nil {
			return nil, err
		}
		d.config.SQL = d.config.SQL + condition
	}

	if d.config.SQL, _, err = d.replaceExprs(d.config.SQL, inputMap, func(expr, old string, value any) string {
		return strings.ReplaceAll(expr, old, reflect.ValueOf(value).String())
	}); err != nil {
		return nil, err
	}

	// sql 替换完成，开始执行
	output, err := d.executeSQL(ctx, d.config.SQL, args)
	if err != nil {
		logx.Infow("[DATABASE组件] 执行失败",
			logx.Field("SQL", d.config.SQL))
		return &core.Result{
			Route:  []string{Failed},
			Output: nil,
		}, err
	}

	logx.Infow("[DATABASE组件] 执行成功",
		logx.Field("SQL", d.config.SQL))
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

func (d *DatabaseComponent) Clear() {
	endItemComponentPool.Put(d)
}

func (d *DatabaseComponent) replaceExprs(expr string, inputMap map[string]any, replaceFun func(expr, old string, value any) string) (string, []any, error) {
	re := regexp.MustCompile(`{{.*?}}`)
	matches := re.FindAllString(expr, -1)

	var ok bool
	var value any
	var args []any
	for _, filed := range matches {
		cutFiled, _ := strings.CutPrefix(filed, "{{")
		cutFiled, _ = strings.CutSuffix(cutFiled, "}}")
		if value, ok = inputMap[cutFiled]; !ok {
			return "", nil, fmt.Errorf("input 中没有 %s 变量", cutFiled)
		}
		if value == nil {
			value = reflect.Zero(reflect.TypeOf(value)).Interface()
		}
		// 替换
		args = append(args, value)
		expr = replaceFun(expr, filed, value)
	}
	return expr, args, nil
}

// 执行SQL语句
func (d *DatabaseComponent) executeSQL(ctx context.Context, sql string, args []interface{}) (any, error) {
	if d.config.DatasourceType == "Oracle" {
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
	rows, err := datasource.DataSourcePool.Query(d.config.DatasourceId, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := queryResult(rows)
	if err != nil {
		logx.Errorw("[DATABASE组件] Query执行失败",
			logx.Field("SQL", sql),
			logx.Field("错误", err))
		return nil, err
	}
	return map[string]any{
		"outputList": result,
		"rowNum":     len(result),
	}, nil
}

// 执行插入
func (d *DatabaseComponent) executeInsert(ctx context.Context, sql string, args []interface{}) (any, error) {
	result, err := datasource.DataSourcePool.Insert(d.config.DatasourceId, sql, args...)
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
	result, err := datasource.DataSourcePool.Update(d.config.DatasourceId, sql, args...)
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
	result, err := datasource.DataSourcePool.Delete(d.config.DatasourceId, sql, args...)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	return map[string]any{
		"outputList": []int64{},
		"rowNum":     affected,
	}, nil
}
