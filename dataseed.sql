INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('0def1c72-a83f-43a6-b6b1-c4a4c589d16b','代码执行💻','process','{\"type\": \"jsTransform\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"jsScript\", \"type\": \"code-input\", \"label\": \"处理脚本\", \"config\": {\"theme\": \"vs-dark\", \"height\": 200, \"options\": {\"minimap\": {\"enabled\": false}, \"fontSize\": 14, \"lineNumbers\": true}, \"language\": \"javascript\", \"defaultValue\": \"function Filter(msg, metadata, msgType) {  \\n  return { msg: msg, metadata: metadata, msgType: msgType };\\n}\"}}], \"runnable\": true}',4);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('0e36fd17-1f91-44d7-b124-346194e7f031','开始😄','input','{\"type\": \"start\", \"point\": {\"inputs\": [], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}]}, \"fields\": [{\"id\": \"param\", \"type\": \"json-input\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": \"{\\\"name\\\": \\\"xuetu\\\",\\\"age\\\": 18}\"}}]}',1);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('11366510-4985-4db8-aab2-3c3500ec4f4e','条件判断','process','{\"type\": \"jsFilter\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"True\", \"type\": \"source\", \"label\": \"True\", \"position\": \"right\", \"handleType\": \"True\"}, {\"id\": \"False\", \"type\": \"source\", \"label\": \"False\", \"position\": \"right\", \"handleType\": \"False\"}, {\"id\": \"Failure\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"jsScript\", \"type\": \"input\", \"label\": \"判断条件\", \"config\": {\"height\": 150, \"defaultValue\": \"msg.?==\'xxx\'\"}}]}',10);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('3c2a2245-480f-4ffb-871e-11b1389a27bf','结束🩷','output','{\"type\": \"end\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}]}, \"fields\": []}',2);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('789a7fc1-9ba9-4805-9dbe-a16f69b1920d','聚合','output','{\"type\": \"join\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}]}, \"fields\": []}',8);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('83caa852-1041-4acb-8ddf-90aa7340e99d','数据库','output','{\"type\": \"database\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"datasource_type\", \"type\": \"radio\", \"label\": \"数据库类型\", \"config\": {\"options\": [{\"label\": \"MySQL\", \"value\": \"MySQL\"}, {\"label\": \"SqlServer\", \"value\": \"SqlServer\"}], \"defaultValue\": \"MySQL\"}}, {\"id\": \"datasource_id\", \"type\": \"input\", \"label\": \"数据源 ID\", \"config\": {\"height\": 150, \"defaultValue\": \"xxxxx\"}}, {\"id\": \"datasource_sql\", \"type\": \"input\", \"label\": \"SQL语句\", \"config\": {\"height\": 150, \"defaultValue\": \"select * from xxx where id = ${id} and name = ${name} limit 10;\"}}, {\"id\": \"datasource_param_mapper\", \"type\": \"dy-form\", \"label\": \"参数映射\", \"config\": {\"height\": 150, \"defaultValue\": [{\"label\": \"${id}\", \"value\": \"msg.id\"}, {\"label\": \"${name}\", \"value\": \"msg.name\"}]}}], \"runnable\": true}',6);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('89a43296-6f1f-4a17-b5c2-e9a2c3a7affd','HTTP-XML','output','{\"type\": \"http-xml\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"method\", \"type\": \"radio\", \"label\": \"请求类型\", \"config\": {\"options\": [{\"label\": \"POST\", \"value\": \"post\"}], \"defaultValue\": \"post\"}}, {\"id\": \"url\", \"type\": \"input\", \"label\": \"地址\", \"config\": {\"height\": 150, \"defaultValue\": \"http://\"}}, {\"id\": \"header\", \"type\": \"dy-form\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": [{\"label\": \"key\", \"value\": \"value\"}]}}, {\"id\": \"xmlParam\", \"type\": \"xml-input\", \"label\": \"请求参数\", \"config\": {\"theme\": \"vs-dark\", \"height\": 200, \"options\": {\"minimap\": {\"enabled\": false}, \"fontSize\": 14, \"lineNumbers\": true}, \"language\": \"xml\", \"defaultValue\": \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?><person><name>张三</name><age>25</age></person>\"}}], \"runnable\": true}',12);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('8d72d9aa-e4eb-40ba-94c1-8e269afb607b','HTTP','input','{\"type\": \"http\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"method\", \"type\": \"radio\", \"label\": \"请求类型\", \"config\": {\"options\": [{\"label\": \"GET\", \"value\": \"get\"}, {\"label\": \"POST\", \"value\": \"post\"}], \"defaultValue\": \"post\"}}, {\"id\": \"url\", \"type\": \"input\", \"label\": \"地址\", \"config\": {\"height\": 150, \"defaultValue\": \"http://\"}}, {\"id\": \"header\", \"type\": \"dy-form\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": [{\"label\": \"key\", \"value\": \"value\"}]}}], \"runnable\": true}',5);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('8fb6bee3-9e88-4c02-a25e-840ebf1f73b6','并发','input','{\"type\": \"fork\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}]}, \"fields\": []}',7);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('9209e3e2-6cd6-4cd5-b0e7-b448fb06a88e','丰富的组件🌹','process','{\"type\": \"jsTransform\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"mode\", \"type\": \"radio\", \"label\": \"运行模式\", \"config\": {\"options\": [{\"label\": \"同步\", \"value\": \"sync\"}, {\"label\": \"异步\", \"value\": \"async\"}], \"defaultValue\": \"sync\"}}, {\"id\": \"features\", \"type\": \"checkbox\", \"label\": \"启用功能\", \"config\": {\"options\": [{\"label\": \"错误重试\", \"value\": \"retry\"}, {\"label\": \"日志记录\", \"value\": \"logging\"}, {\"label\": \"数据缓存\", \"value\": \"cache\"}], \"defaultValue\": [\"logging\"]}}, {\"id\": \"retryCount\", \"type\": \"slider\", \"label\": \"重试次数\", \"config\": {\"max\": 10, \"min\": 0, \"step\": 1, \"marks\": {\"0\": \"0\", \"5\": \"5\", \"10\": \"10\"}, \"defaultValue\": 3}}, {\"id\": \"timeout\", \"type\": \"number\", \"label\": \"超时时间(ms)\", \"config\": {\"max\": 10000, \"min\": 1000, \"step\": 1000, \"defaultValue\": 5000}}, {\"id\": \"headers\", \"type\": \"json-input\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": {\"Accept\": \"application/json\", \"Content-Type\": \"application/json\"}}}, {\"id\": \"script\", \"type\": \"code-input\", \"label\": \"处理脚本\", \"config\": {\"theme\": \"tomorrow\", \"height\": 200, \"language\": \"javascript\", \"defaultValue\": \"//请求发送前的数据处理\\nfunction preprocess(data) {\\n  return data;\\n}\"}}, {\"id\": \"dataSource\", \"type\": \"select\", \"label\": \"数据源\", \"config\": {\"mode\": \"multiple\", \"options\": [{\"label\": \"MySQL\", \"value\": \"mysql\"}, {\"label\": \"MongoDB\", \"value\": \"mongodb\"}, {\"label\": \"Redis\", \"value\": \"redis\"}], \"defaultValue\": []}}, {\"id\": \"priority\", \"type\": \"select\", \"label\": \"优先级\", \"config\": {\"options\": [{\"label\": \"高\", \"value\": \"high\"}, {\"label\": \"中\", \"value\": \"medium\"}, {\"label\": \"低\", \"value\": \"low\"}], \"defaultValue\": \"medium\"}}], \"runnable\": true}',3);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('95bff99e-4b9c-4afd-94ca-9615351ef057','迭代','output','{\"type\": \"for\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Do\", \"type\": \"source\", \"label\": \"迭代\", \"position\": \"right\", \"handleType\": \"Do\"}, {\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"range\", \"type\": \"input\", \"label\": \"迭代对象\", \"config\": {\"height\": 150, \"defaultValue\": \"msg.?\"}}]}',9);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('bc6963b0-2fbc-4104-b9e0-3f22f1f8c8fd','文件服务器','process','{\"type\": \"fileServer\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"datasource_type\", \"type\": \"radio\", \"label\": \"文件系统类型\", \"config\": {\"options\": [{\"label\": \" SFTP\", \"value\": \"sftp\"}, {\"label\": \"FTP\", \"value\": \"ftp\"}], \"defaultValue\": \"ftp\"}}, {\"id\": \"datasource_mode\", \"type\": \"radio\", \"label\": \"操作模式\", \"config\": {\"options\": [{\"label\": \"上传\", \"value\": \"upload\"}, {\"label\": \"下载\", \"value\": \"download\"}, {\"label\": \"删除\", \"value\": \"delete\"}], \"defaultValue\": \"upload\"}}, {\"id\": \"datasource_id\", \"type\": \"input\", \"label\": \"数据源 ID\", \"config\": {\"height\": 150, \"defaultValue\": \"xx\"}}, {\"id\": \"datasource_path\", \"type\": \"input\", \"label\": \"路径\", \"config\": {\"height\": 150, \"defaultValue\": \"/data/xxx.xxx\"}}], \"runnable\": true}',11);