package rulego

import (
	"os"
	"testing"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/config"
	"workflow/internal/datasource"
	"workflow/internal/svc"
)

func Test_roleChain_Run_Mysql_Select(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select.json",
		"{\"name\":\"xuetu-1\"}",
	)
}

func Test_roleChain_Run_Mysql_Select_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_NoTableName.json",
		"{\"table\":\"user\",\"name\":\"xuetu-1\"}",
	)
}

func Test_roleChain_Run_Mysql_Select_LeftJoin(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_LeftJoin.json",
		"{\"table1\":\"user\",\"table2\":\"user_info\",\"name\":\"xuetu-1\"}",
	)

	setupAndRunChain(t,
		"./chain/mysql/mysql_select_LeftJoin.json",
		"{\"table1\":\"user\",\"table2\":\"user_info\",\"name\":\"xuetu-2\"}",
	)
}

func Test_roleChain_Run_Mysql_Insert(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_insert.json",
		"{\"table\":\"user\",\"name\":\"xuetu-3\",\"age\":18}",
	)

	setupAndRunChain(t,
		"./chain/mysql/mysql_insert.json",
		"{\"table\":\"user\",\"name\":\"xuetu-4\",\"age\":19}",
	)
}

func Test_roleChain_Run_Mysql_Insert_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_insert_NoTableName.json",
		"{\"table\":\"users\",\"name\":\"xuetu-5\",\"age\":18}",
	)

	setupAndRunChain(t,
		"./chain/mysql/mysql_insert_NoTableName.json",
		"{\"table\":\"users\",\"name\":\"xuetu-6\",\"age\":19}",
	)
}

func Test_roleChain_Run_Mysql_Update(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_update.json",
		"{\"table\":\"users\",\"name\":\"xuetu-7\",\"id\":1}",
	)
}

func Test_roleChain_Run_Mysql_Update_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_update_NoTableName.json",
		"{\"table\":\"users\",\"name\":\"xuetu-8\",\"id\":1}",
	)
}

func Test_roleChain_Run_Mysql_Delete(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_delete.json",
		"{\"table\":\"users\",\"name\":\"xuetu-8\",\"id\":1}",
	)
}

func Test_roleChain_Run_Mysql_DeleteNotableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_delete_NoTableName.json",
		"{\"table\":\"users\",\"name\":\"xuetu-2\",\"id\":1}",
	)
}

// 多表查询 非等值连接
func Test_roleChain_Run_Mysql_Select_Multi_Table(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_multi_table.json",
		"{\"table\":\"users\",\"name\":\"xuetu-2\",\"id\":1}",
	)
}

func Test_roleChain_Run_Mysql_Select_Multi_Table_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_multi_table_NoTableName.json",
		"{\"table1\":\"employees\",\"table2\":\"employees\"}",
	)
}

// inner join
func Test_roleChain_Run_Mysql_Select_Inner_join(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_inner_join.json",
		"{\"table1\":\"employees\",\"table2\":\"employees\"}",
	)
}

func Test_roleChain_Run_Mysql_Select_Inner_join_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_inner_join_NoTableName.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\"}",
	)
}

// left join
func Test_roleChain_Run_Mysql_Select_LeftJoin_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_left_join_NoTableName.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\"}",
	)
}

// right join
func Test_roleChain_Run_Mysql_Select_RightJoin(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_right_join.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\"}",
	)
}

func Test_roleChain_Run_Mysql_Select_RightJoin_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_right_join_NoTableName.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\"}",
	)
}

// union all
func Test_roleChain_Run_Mysql_Select_UnionAll(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_union_all.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\"}",
	)
}

func Test_roleChain_Run_Mysql_Select_UnionAll_NoTableName(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_union_all_NoTableName.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\",\"table3\":\"employees\",\"table4\":\"departments\"}",
	)
}

// 子查询
func Test_roleChain_Run_Mysql_Select_SubSelect(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_sub_select.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\",\"table3\":\"employees\",\"table4\":\"departments\"}",
	)
}

// 聚合函数
func Test_roleChain_Run_Mysql_Select_Function(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_function.json",
		"{\"table1\":\"employees\",\"table2\":\"departments\",\"table3\":\"employees\",\"table4\":\"departments\"}",
	)
}

// 排序分页
func Test_roleChain_Run_Mysql_Select_OrderBy_Limit(t *testing.T) {
	setupAndRunChain(t,
		"./chain/mysql/mysql_select_order_by_limit.json",
		"{\"table\":\"employees\",\"salary\":6000,\"order\":\"salary DESC\",\"limit\":10}",
	)
}

/*
测试数据

CREATE TABLE `user_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `phone` varchar(255) NOT NULL,
  `sex` varchar(255) NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) NOT NULL,
  `age` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

*/

// 测试辅助函数,用于初始化和运行规则链
func setupAndRunChain(t *testing.T, jsonFile string, data string) types.RuleMsg {
	// 初始化数据源
	svcCtx := svc.NewServiceContext(config.Config{
		// 其他测试数据库,MySqlUrn: "xxx",
		MySqlUrn: "root:Root@123@tcp(10.99.43.9:3306)/workflow?charset=utf8mb4&parseTime=True&loc=Local",
	})
	datasource.InitDataSourceManager(svcCtx)

	// 读取配置文件
	file, _ := os.ReadFile(jsonFile)

	// 设置配置
	config := rulego.NewConfig()
	logConf := logx.LogConf{Encoding: "plain"}
	config.RegisterUdf("log", func(msg interface{}) {
		logx.Debugf("log:%+v", msg)
	})
	logx.SetUp(logConf)

	// 创建规则链
	chain, err := rulego.New(
		"ctg1kid3sjti2l614lp0",
		file,
		rulego.WithConfig(config),
		types.WithAspects(&DebugAop{}),
	)
	if err != nil {
		t.Fatalf("load role chain fail: %v", err)
	}

	// 准备消息
	metadata := map[string]string{"env": "jlhalsjdhfoisdbv"}
	msg := types.NewMsg(0, "CANVAS_MSG", types.JSON, metadata, data)

	// 运行规则链
	var result types.RuleMsg
	chain.OnMsgAndWait(msg, types.WithOnEnd(func(ctx types.RuleContext, msg types.RuleMsg, err error, relationType string) {
		result = msg
		if err != nil {
			t.Logf("chain run error: %+v", err)
		}
	}))

	t.Logf("chain run result: %+v", result)
	return result
}
