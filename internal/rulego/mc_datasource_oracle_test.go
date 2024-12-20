package rulego

import (
	"testing"
)

func Test_roleChain_Run_Oracle_Insert(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_insert.json",
		"{\"id\":11,\"name\":\"wuyuli\",\"age\":18}",
	)
}

func Test_roleChain_Run_Oracle_Select(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select.json",
		"{\"id\":2,\"name\":\"wuyuli\",\"age\":18}",
	)
}

func Test_roleChain_Run_Oracle_Update(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_update.json",
		"{\"table\":\"users\",\"id\":11,\"name\":\"wuyuli\",\"age\":18}",
	)
}

func Test_roleChain_Run_Oracle_Delete(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_delete.json",
		"{\"table\":\"users\",\"id\":11}",
	)
}

func Test_roleChain_Run_Oracle_Select_Function(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select_function.json",
		"{\"table\":\"T_ACCOUNT\",\"id\":11}",
	)
}

func Test_roleChain_Run_Oracle_Select_Inner_Join(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select_inner_join.json",
		"{\"table1\":\"T_OWNERS\",\"table2\":\"T_OWNERTYPE\",\"table3\":\"T_ADDRESS\",\"table4\":\"T_AREA\"}",
	)
}

func Test_roleChain_Run_Oracle_Select_Left_Join(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select_left_join.json",
		"{\"table1\":\"T_OWNERS\",\"table2\":\"T_ACCOUNT\"}",
	)
}

func Test_roleChain_Run_Oracle_Select_Right_Join(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select_right_join.json",
		"{\"table1\":\"T_OWNERS\",\"table2\":\"T_ACCOUNT\"}",
	)
}

func Test_roleChain_Run_Oracle_Select_Sub_Select(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select_sub_select.json",
		"{\"table1\":\"T_OWNERS\",\"table2\":\"T_ACCOUNT\"}",
	)
}

func Test_roleChain_Run_Oracle_Select_OrderBy_Limit(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_select_order_by_limit.json",
		"{\"table1\":\"T_OWNERS\",\"table2\":\"T_ACCOUNT\"}",
	)
}
