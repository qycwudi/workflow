package rulego

import (
	"regexp"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/datasource"
)

// StartNode A plugin that flow start node ,receiving parameter
type DataSourceMysqlNode struct {
	Config DatabaseNodeConfiguration
}

type DatabaseNodeConfiguration struct {
	DatasourceType        string            `json:"datasourceType"`
	DatasourceId          int64             `json:"datasourceId"`
	DatasourceSql         string            `json:"datasourceSql"`
	DatasourceParamMapper map[string]string `json:"datasourceParamMapper"`
}

func init() {
	_ = rulego.Registry.Register(&DataSourceMysqlNode{})
}

func (n *DataSourceMysqlNode) Type() string {
	return Database
}
func (n *DataSourceMysqlNode) New() types.Node {
	config := DatabaseNodeConfiguration{
		DatasourceType:        "mysql",
		DatasourceId:          1,
		DatasourceSql:         "select * from test",
		DatasourceParamMapper: map[string]string{},
	}
	return &DataSourceMysqlNode{Config: config}
}

func (n *DataSourceMysqlNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	// 读取配置
	marshal, _ := json.Marshal(configuration)
	_ = json.Unmarshal(marshal, &n.Config)
	return nil
}

// OnMsg 处理消息
func (n *DataSourceMysqlNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	// 解析消息中的参数
	msgData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(msg.Data), &msgData); err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	// SQL参数替换
	sql := n.Config.DatasourceSql
	var args []interface{}

	// 使用正则表达式找出所有${xxx}参数
	re := regexp.MustCompile(`\${([^}]+)}`)
	matches := re.FindAllStringSubmatch(sql, -1)

	// 按SQL中参数出现顺序收集参数值
	for _, match := range matches {
		placeholder := match[0] // ${xxx}
		paramName := n.Config.DatasourceParamMapper[placeholder]
		if val, ok := msgData[paramName]; ok {
			// 检查是否是表名参数(在from子句后面)
			// 使用正则表达式检查参数是否在FROM子句后面
			fromRe := regexp.MustCompile(`(?i)FROM\s+` + regexp.QuoteMeta(placeholder))
			if fromRe.MatchString(sql) {
				// 表名不作为预处理参数,直接替换
				sql = regexp.MustCompile(regexp.QuoteMeta(placeholder)).ReplaceAllString(sql, val.(string))
				continue
			}
			args = append(args, val)
		}
	}

	// 替换剩余的${xxx}为?
	sql = re.ReplaceAllString(sql, "?")

	logx.Infof("sql:%s,args:%+v", sql, args)
	// 执行SQL
	rows, err := datasource.DataSourcePool.Query(n.Config.DatasourceId, sql, args...)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		logx.Errorf("get columns failed: %v", err)
		ctx.TellFailure(msg, err)
		return
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
			ctx.TellFailure(msg, err)
			return
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
		ctx.TellFailure(msg, err)
		return
	}

	// 将结果赋值给msg的数据部分
	marshal, _ := json.Marshal(result)
	msg.Data = string(marshal)
	ctx.TellSuccess(msg)
}

func (n *DataSourceMysqlNode) Destroy() {

	// Do some cleanup work
}
