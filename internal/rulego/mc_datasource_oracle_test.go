package rulego

import (
	"testing"
)

/*
oracle版本：oracle_11g

-- 价格区间
create table t_pricetable
(
    id 			number primary key,
    price 		number(10,2),
    ownertypeid number,
    minnum 		number,
    maxnum 		number
);

--业主类型
create table t_ownertype
(
    id 		number primary key,
    name 	varchar2(30)
);

--业主表
create table t_owners
(
    id 			number primary key,
    name 		varchar2(30),
    addressid 	number,
    housenumber varchar2(30),
    watermeter 	varchar2(30),
    adddate 	date,
    ownertypeid number
);

--区域表
create table t_area
(
    id 		number,
    name 	varchar2(30)
);

--收费员表
create table t_operator
(
    id 		number,
    name 	varchar2(30)
);

--地址表
create table t_address
(
    id 			number primary key,
    name 		varchar2(100),
    areaid 		number,
    operatorid 	number
);

--账务表--
create table t_account
(
    id 			number primary key,
    owneruuid 	number,
    ownertype 	number,
    areaid 		number,
    year 		char(4),
    month 		char(2),
    num0 		number,
    num1 		number,
    usenum 		number,
    meteruser 	number,
    meterdate 	date,
    money 		number(10,2),
    isfee 		char(1),
    feedate 	date,
    feeuser 	number
);

*/

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
