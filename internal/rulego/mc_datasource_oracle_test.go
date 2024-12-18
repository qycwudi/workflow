package rulego

import "testing"

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

func Test_roleChain_Run_Oracle_Inner_Join(t *testing.T) {
	setupAndRunChain(t,
		"./chain/oracle/oracle_inner_join.json",
		"{\"table\":\"users\",\"id\":11}",
	)
}
