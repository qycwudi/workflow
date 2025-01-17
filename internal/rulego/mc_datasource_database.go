package rulego

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
	str "github.com/rulego/rulego/utils/str"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/datasource"
)

// DataSourceDatabaseNode A plugin that flow start node ,receiving parameter
type DataSourceDatabaseNode struct {
	Config DataSourceDatabaseNodeConfiguration
}

type DataSourceDatabaseNodeConfiguration struct {
	DatasourceType        string            `json:"datasourceType"`
	DatasourceId          int64             `json:"datasourceId"`
	DatasourceSql         string            `json:"datasourceSql"`
	DatasourceParamMapper map[string]string `json:"datasourceParamMapper"`
}

func init() {
	_ = rulego.Registry.Register(&DataSourceDatabaseNode{})
}

func (n *DataSourceDatabaseNode) Type() string {
	return Database
}
func (n *DataSourceDatabaseNode) New() types.Node {
	config := DataSourceDatabaseNodeConfiguration{
		DatasourceType:        n.Config.DatasourceType,
		DatasourceId:          n.Config.DatasourceId,
		DatasourceSql:         n.Config.DatasourceSql,
		DatasourceParamMapper: n.Config.DatasourceParamMapper,
	}
	return &DataSourceDatabaseNode{Config: config}
}

func (n *DataSourceDatabaseNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	// 读取配置
	marshal, _ := json.Marshal(configuration)
	_ = json.Unmarshal(marshal, &n.Config)
	return nil
}

// OnMsg 处理消息
func (n *DataSourceDatabaseNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 解析消息参数
	msgData, err := parseMessageData(msg.Data, msg.Metadata)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	// 处理SQL语句和参数
	sql, args, err := n.processSQLAndParams(msgData)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	logx.Infof("sql:%s,args:%+v", sql, args)

	// 执行SQL
	if err := n.executeSQL(ctx, msg, sql, args); err != nil {
		ctx.TellFailure(msg, err)
		return
	}
}

func (n *DataSourceDatabaseNode) Destroy() {}

// 解析消息数据为map
func parseMessageData(data string, metadata map[string]string) (map[string]interface{}, error) {
	msgData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(data), &msgData); err != nil {
		return nil, err
	}
	env := make(map[string]interface{})
	// 添加metadata到msgData
	env["metadata"] = metadata
	// 添加msgData到msgData
	env["msg"] = msgData
	return env, nil
}

// 处理SQL语句和参数
func (n *DataSourceDatabaseNode) processSQLAndParams(msgData map[string]interface{}) (string, []interface{}, error) {
	sql := n.Config.DatasourceSql
	var args []interface{}

	// 使用正则表达式找出所有${xxx}参数
	re := regexp.MustCompile(`\${([^}]+)}`)
	matches := re.FindAllStringSubmatch(sql, -1)

	// 处理参数
	for _, match := range matches {
		placeholder := match[0]
		val := str.ExecuteTemplate(placeholder, msgData)
		// 检查是否是表名参数
		if isTableNameParam(sql, placeholder) {
			sql = replaceTableName(sql, placeholder, val)
			continue
		}
		args = append(args, val)
	}

	// 替换剩余的${xxx}为?
	sql = re.ReplaceAllString(sql, "?")
	return sql, args, nil
}

// 检查是否是表名参数
func isTableNameParam(sql string, placeholder string) bool {
	// 预编译正则表达式以提高性能
	quotedPlaceholder := regexp.QuoteMeta(placeholder)
	// (?i) 表示不区分大小写
	pattern := fmt.Sprintf(`(?i)(FROM|INTO|UPDATE|JOIN)\s+%s|(?i)(CREATE|DROP|ALTER|TRUNCATE)\s+TABLE\s+%s|(?i)(FROM|JOIN)\s+([^\s,]+)(?:\s+AS\s+|\s)([^\s,]+).*,\s+%s`,
		quotedPlaceholder,
		quotedPlaceholder,
		quotedPlaceholder)

	re := regexp.MustCompile(pattern)
	return re.MatchString(sql)
}

// 替换表名
func replaceTableName(sql string, placeholder string, tableName string) string {
	return regexp.MustCompile(regexp.QuoteMeta(placeholder)).ReplaceAllString(sql, tableName)
}

// 执行SQL语句
func (n *DataSourceDatabaseNode) executeSQL(ctx types.RuleContext, msg types.RuleMsg, sql string, args []interface{}) error {
	if n.Config.DatasourceType == "Oracle" {
		sql = strings.ReplaceAll(sql, ";", "")
	}

	sqlType := strings.ToUpper(strings.TrimSpace(sql))
	switch {
	case strings.HasPrefix(sqlType, "SELECT"):
		return n.executeQuery(ctx, msg, sql, args)
	case strings.HasPrefix(sqlType, "INSERT"):
		return n.executeInsert(ctx, msg, sql, args)
	case strings.HasPrefix(sqlType, "UPDATE"):
		return n.executeUpdate(ctx, msg, sql, args)
	case strings.HasPrefix(sqlType, "DELETE"):
		return n.executeDelete(ctx, msg, sql, args)
	default:
		return fmt.Errorf("unsupported SQL type: %s", sqlType)
	}
}

// 执行查询
func (n *DataSourceDatabaseNode) executeQuery(ctx types.RuleContext, msg types.RuleMsg, sql string, args []interface{}) error {
	rows, err := datasource.DataSourcePool.Query(n.Config.DatasourceId, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	result, err := queryResult(rows)
	if err != nil {
		return err
	}

	msg.Data = string(result)
	ctx.TellSuccess(msg)
	return nil
}

func queryResult(rows *sql.Rows) ([]byte, error) {

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		logx.Errorf("get columns failed: %v", err)
		return nil, err
	}

	// 存储所有行的结果
	var result []map[string]interface{}

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
			logx.Errorf("scan failed: %v", err)
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
	marshal, _ := json.Marshal(result)
	return marshal, nil
}

// 执行插入
func (n *DataSourceDatabaseNode) executeInsert(ctx types.RuleContext, msg types.RuleMsg, sql string, args []interface{}) error {
	result, err := datasource.DataSourcePool.Insert(n.Config.DatasourceId, sql, args...)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	msg.Data = fmt.Sprintf(`{"affected":%d,"lastInsertId":%d}`, 1, id)
	ctx.TellSuccess(msg)
	return nil
}

// 执行更新
func (n *DataSourceDatabaseNode) executeUpdate(ctx types.RuleContext, msg types.RuleMsg, sql string, args []interface{}) error {
	result, err := datasource.DataSourcePool.Update(n.Config.DatasourceId, sql, args...)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	msg.Data = fmt.Sprintf(`{"affected":%d}`, affected)
	ctx.TellSuccess(msg)
	return nil
}

// 执行删除
func (n *DataSourceDatabaseNode) executeDelete(ctx types.RuleContext, msg types.RuleMsg, sql string, args []interface{}) error {
	result, err := datasource.DataSourcePool.Delete(n.Config.DatasourceId, sql, args...)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	msg.Data = fmt.Sprintf(`{"affected":%d}`, affected)
	ctx.TellSuccess(msg)
	return nil
}
