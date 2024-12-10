mysqldump: [Warning] Using a password on the command line interface can be insecure.
-- MySQL dump 10.13  Distrib 8.0.36, for Linux (aarch64)
--
-- Host: localhost    Database: wk
-- ------------------------------------------------------
-- Server version	8.0.36

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `api`
--

DROP TABLE IF EXISTS `api`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `api` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) NOT NULL,
  `api_id` varchar(255) NOT NULL,
  `api_name` varchar(255) NOT NULL,
  `api_desc` text NOT NULL,
  `dsl` json NOT NULL,
  `status` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='api服务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `api`
--

LOCK TABLES `api` WRITE;
/*!40000 ALTER TABLE `api` DISABLE KEYS */;
/*!40000 ALTER TABLE `api` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `api_record`
--

DROP TABLE IF EXISTS `api_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `api_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `status` varchar(255) NOT NULL,
  `trace_id` varchar(255) NOT NULL,
  `param` json NOT NULL COMMENT '参数',
  `extend` json NOT NULL COMMENT '扩展',
  `call_time` datetime NOT NULL,
  `api_id` varchar(255) NOT NULL,
  `api_name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='api调用记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `api_record`
--

LOCK TABLES `api_record` WRITE;
/*!40000 ALTER TABLE `api_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `api_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `api_secret_key`
--

DROP TABLE IF EXISTS `api_secret_key`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `api_secret_key` (
  `id` int NOT NULL AUTO_INCREMENT,
  `secret_key` varchar(255) NOT NULL,
  `api_id` varchar(255) NOT NULL,
  `expiration_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `nidx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='api 密钥表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `api_secret_key`
--

LOCK TABLES `api_secret_key` WRITE;
/*!40000 ALTER TABLE `api_secret_key` DISABLE KEYS */;
/*!40000 ALTER TABLE `api_secret_key` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `canvas`
--

DROP TABLE IF EXISTS `canvas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `canvas` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) NOT NULL,
  `draft` json NOT NULL,
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  `create_by` varchar(255) NOT NULL,
  `update_by` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `canvas`
--

LOCK TABLES `canvas` WRITE;
/*!40000 ALTER TABLE `canvas` DISABLE KEYS */;
INSERT INTO `canvas` VALUES (17,'ct1i8ud3sjtg3kkhmvq0','{\"id\": \"ct1i8ud3sjtg3kkhmvq0\", \"graph\": {\"name\": \"AI Flow Canvas\", \"edges\": [{\"id\": \"reactflow__edge-custom-956080af-3126-4862-b2cd-d61c429dea28success-custom-b29e78e6-dcfa-41a2-90c7-223299e65f2einput\", \"type\": \"custom\", \"source\": \"custom-956080af-3126-4862-b2cd-d61c429dea28\", \"target\": \"custom-b29e78e6-dcfa-41a2-90c7-223299e65f2e\", \"animated\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}, {\"id\": \"reactflow__edge-custom-b29e78e6-dcfa-41a2-90c7-223299e65f2esuccess-custom-57714fa9-6c66-4ad9-b0e1-43bfaf6f7305input\", \"type\": \"custom\", \"source\": \"custom-b29e78e6-dcfa-41a2-90c7-223299e65f2e\", \"target\": \"custom-57714fa9-6c66-4ad9-b0e1-43bfaf6f7305\", \"animated\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}], \"nodes\": [{\"id\": \"custom-b29e78e6-dcfa-41a2-90c7-223299e65f2e\", \"data\": {\"name\": \"代码执行💻\", \"type\": \"process\", \"moduleConfig\": \"{\\\"type\\\": \\\"start\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}, {\\\"id\\\": \\\"error\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"失败\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Fail\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"script\\\", \\\"type\\\": \\\"code-input\\\", \\\"label\\\": \\\"处理脚本\\\", \\\"config\\\": {\\\"theme\\\": \\\"vs-dark\\\", \\\"height\\\": 200, \\\"options\\\": {\\\"minimap\\\": {\\\"enabled\\\": false}, \\\"fontSize\\\": 14, \\\"lineNumbers\\\": true}, \\\"language\\\": \\\"javascript\\\", \\\"defaultValue\\\": \\\"//请求发送前的数据处理\\\\n function preprocess(data) {\\\\nreturn data;}\\\"}}], \\\"runnable\\\": true}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 356, \"dragging\": false, \"position\": {\"x\": 1395, \"y\": -360}, \"selected\": false, \"positionAbsolute\": {\"x\": 1395, \"y\": -360}}, {\"id\": \"custom-956080af-3126-4862-b2cd-d61c429dea28\", \"data\": {\"name\": \"开始😄\", \"type\": \"input\", \"value\": {\"headers\": \"{\\\"Accept\\\": \\\"json\\\",\\\"Content-Type\\\": \\\"application/json\\\"}\"}, \"moduleConfig\": \"{\\\"type\\\": \\\"start\\\", \\\"point\\\": {\\\"inputs\\\": [], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"headers\\\", \\\"type\\\": \\\"json-input\\\", \\\"label\\\": \\\"请求头\\\", \\\"config\\\": {\\\"height\\\": 150, \\\"defaultValue\\\": \\\"{\\\\\\\"Accept\\\\\\\": \\\\\\\"application/json\\\\\\\",\\\\\\\"Content-Type\\\\\\\": \\\\\\\"application/json\\\\\\\"}\\\"}}]}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 306, \"dragging\": false, \"position\": {\"x\": 0, \"y\": -795}, \"selected\": false, \"positionAbsolute\": {\"x\": 0, \"y\": -795}}, {\"id\": \"custom-57714fa9-6c66-4ad9-b0e1-43bfaf6f7305\", \"data\": {\"name\": \"结束🩷\", \"type\": \"output\", \"moduleConfig\": \"{\\\"type\\\": \\\"end\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}]}, \\\"fields\\\": []}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 71, \"dragging\": false, \"position\": {\"x\": 1110, \"y\": 270}, \"selected\": false, \"positionAbsolute\": {\"x\": 1110, \"y\": 270}}], \"version\": \"1.0\", \"metadata\": {\"edgeCount\": 2, \"nodeCount\": 3, \"lastModified\": \"2024-12-06T11:16:39.785Z\"}, \"exportedAt\": \"2024-12-06T11:16:39.785Z\", \"description\": \"Canvas configuration\"}}','2024-11-24 21:04:57','2024-12-06 19:16:40','system','system'),(18,'ct1l3jl3sjthu24kms8g','{\"id\": \"ct1l3jl3sjthu24kms8g\", \"graph\": {\"name\": \"AI Flow Canvas\", \"edges\": [], \"nodes\": [], \"version\": \"1.0\", \"metadata\": {\"edgeCount\": 0, \"nodeCount\": 0, \"lastModified\": \"2024-11-24T16:37:44.415Z\"}, \"exportedAt\": \"2024-11-24T16:37:44.415Z\", \"description\": \"Canvas configuration\"}}','2024-11-25 00:18:22','2024-11-25 00:37:44','system','system'),(19,'ct4qi1t3sjtp5jbouq0g','{\"id\": \"ct4qi1t3sjtp5jbouq0g\"}','2024-11-29 19:44:08','2024-11-29 19:44:08','system','system'),(20,'ct9dq9d3sjtl7g6q5ha0','{\"id\": \"ct9dq9d3sjtl7g6q5ha0\", \"graph\": {\"name\": \"AI Flow Canvas\", \"edges\": [{\"id\": \"reactflow__edge-custom-d8188062-667c-490e-a9d9-1cdb4dd20a0dsuccess-custom-8a5340f2-8bcf-4dee-91a7-aa3bc65573c4input\", \"type\": \"custom\", \"source\": \"custom-d8188062-667c-490e-a9d9-1cdb4dd20a0d\", \"target\": \"custom-8a5340f2-8bcf-4dee-91a7-aa3bc65573c4\", \"animated\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}, {\"id\": \"reactflow__edge-custom-d8188062-667c-490e-a9d9-1cdb4dd20a0dsuccess-custom-059d1382-fa3a-43fd-87e8-71b6f2535cbainput\", \"type\": \"custom\", \"source\": \"custom-d8188062-667c-490e-a9d9-1cdb4dd20a0d\", \"target\": \"custom-059d1382-fa3a-43fd-87e8-71b6f2535cba\", \"animated\": true, \"deletable\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}, {\"id\": \"reactflow__edge-custom-d8188062-667c-490e-a9d9-1cdb4dd20a0derror-custom-bf1087e3-5f11-4509-af66-5620cc81dde4input\", \"type\": \"custom\", \"source\": \"custom-d8188062-667c-490e-a9d9-1cdb4dd20a0d\", \"target\": \"custom-bf1087e3-5f11-4509-af66-5620cc81dde4\", \"animated\": true, \"deletable\": true, \"sourceHandle\": \"error\", \"targetHandle\": \"input\"}, {\"id\": \"reactflow__edge-custom-bf1087e3-5f11-4509-af66-5620cc81dde4success-custom-771f443a-3c69-4520-b0b9-aa9d3bdf5e4binput\", \"type\": \"custom\", \"source\": \"custom-bf1087e3-5f11-4509-af66-5620cc81dde4\", \"target\": \"custom-771f443a-3c69-4520-b0b9-aa9d3bdf5e4b\", \"animated\": true, \"deletable\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}, {\"id\": \"reactflow__edge-custom-2cdb2e43-552d-4b6e-8d05-b66d787c5b10success-custom-d8188062-667c-490e-a9d9-1cdb4dd20a0dinput\", \"type\": \"custom\", \"source\": \"custom-2cdb2e43-552d-4b6e-8d05-b66d787c5b10\", \"target\": \"custom-d8188062-667c-490e-a9d9-1cdb4dd20a0d\", \"animated\": true, \"deletable\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}], \"nodes\": [{\"id\": \"custom-d8188062-667c-490e-a9d9-1cdb4dd20a0d\", \"data\": {\"name\": \"代码执行💻\", \"type\": \"process\", \"moduleConfig\": \"{\\\"type\\\": \\\"start\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}, {\\\"id\\\": \\\"error\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"失败\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Fail\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"script\\\", \\\"type\\\": \\\"code-input\\\", \\\"label\\\": \\\"处理脚本\\\", \\\"config\\\": {\\\"theme\\\": \\\"vs-dark\\\", \\\"height\\\": 200, \\\"options\\\": {\\\"minimap\\\": {\\\"enabled\\\": false}, \\\"fontSize\\\": 14, \\\"lineNumbers\\\": true}, \\\"language\\\": \\\"javascript\\\", \\\"defaultValue\\\": \\\"function Filter(msg, metadata, msgType) {  \\\\n  return { msg: msg, metadata: metadata, msgType: msgType };\\\\n}\\\"}}], \\\"runnable\\\": true}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 356, \"dragging\": true, \"position\": {\"x\": 585, \"y\": 150}, \"selected\": false, \"positionAbsolute\": {\"x\": 585, \"y\": 150}}, {\"id\": \"custom-059d1382-fa3a-43fd-87e8-71b6f2535cba\", \"data\": {\"name\": \"结束🩷\", \"type\": \"output\", \"moduleConfig\": \"{\\\"type\\\": \\\"end\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}]}, \\\"fields\\\": []}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 71, \"dragging\": true, \"position\": {\"x\": 1150, \"y\": 50}, \"selected\": false, \"positionAbsolute\": {\"x\": 1275, \"y\": 15}}, {\"id\": \"custom-bf1087e3-5f11-4509-af66-5620cc81dde4\", \"data\": {\"name\": \"代码执行💻\", \"type\": \"process\", \"moduleConfig\": \"{\\\"type\\\": \\\"start\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}, {\\\"id\\\": \\\"error\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"失败\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Fail\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"script\\\", \\\"type\\\": \\\"code-input\\\", \\\"label\\\": \\\"处理脚本\\\", \\\"config\\\": {\\\"theme\\\": \\\"vs-dark\\\", \\\"height\\\": 200, \\\"options\\\": {\\\"minimap\\\": {\\\"enabled\\\": false}, \\\"fontSize\\\": 14, \\\"lineNumbers\\\": true}, \\\"language\\\": \\\"javascript\\\", \\\"defaultValue\\\": \\\"function Filter(msg, metadata, msgType) {  \\\\n  return { msg: msg, metadata: metadata, msgType: msgType };\\\\n}\\\"}}], \\\"runnable\\\": true}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 356, \"dragging\": false, \"position\": {\"x\": 1155, \"y\": 555}, \"selected\": false, \"positionAbsolute\": {\"x\": 1155, \"y\": 555}}, {\"id\": \"custom-771f443a-3c69-4520-b0b9-aa9d3bdf5e4b\", \"data\": {\"name\": \"结束🩷\", \"type\": \"output\", \"moduleConfig\": \"{\\\"type\\\": \\\"end\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}]}, \\\"fields\\\": []}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 71, \"position\": {\"x\": 1700, \"y\": 550}}, {\"id\": \"custom-2cdb2e43-552d-4b6e-8d05-b66d787c5b10\", \"data\": {\"name\": \"开始😄\", \"type\": \"input\", \"custom\": {\"param\": \"handleFi\"}, \"moduleConfig\": \"{\\\"type\\\": \\\"start\\\", \\\"point\\\": {\\\"inputs\\\": [], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"param\\\", \\\"type\\\": \\\"json-input\\\", \\\"label\\\": \\\"请求头\\\", \\\"config\\\": {\\\"height\\\": 150, \\\"defaultValue\\\": \\\"{\\\\\\\"Accept\\\\\\\": \\\\\\\"application/json\\\\\\\",\\\\\\\"Content-Type\\\\\\\": \\\\\\\"application/json\\\\\\\"}\\\"}}]}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 306, \"dragging\": true, \"position\": {\"x\": -30, \"y\": -60}, \"selected\": true, \"positionAbsolute\": {\"x\": -30, \"y\": -60}}], \"version\": \"1.0\", \"metadata\": {\"edgeCount\": 5, \"nodeCount\": 5, \"lastModified\": \"2024-12-09T06:14:37.087Z\"}, \"exportedAt\": \"2024-12-09T06:14:37.087Z\", \"description\": \"Canvas configuration\"}}','2024-12-06 19:16:53','2024-12-09 14:14:38','system','system'),(21,'ctb8m0t3sjtq9abeusm0','{\"id\": \"ctb8m0t3sjtq9abeusm0\", \"graph\": {\"name\": \"AI Flow Canvas\", \"edges\": [{\"id\": \"reactflow__edge-custom-72eda4dc-3f22-40a6-a16d-d04152918406success-custom-a00bc6a6-9998-46e2-8f9c-2536fd11bdb4input\", \"type\": \"custom\", \"source\": \"custom-72eda4dc-3f22-40a6-a16d-d04152918406\", \"target\": \"custom-a00bc6a6-9998-46e2-8f9c-2536fd11bdb4\", \"animated\": true, \"deletable\": true, \"sourceHandle\": \"success\", \"targetHandle\": \"input\"}], \"nodes\": [{\"id\": \"custom-a00bc6a6-9998-46e2-8f9c-2536fd11bdb4\", \"data\": {\"name\": \"代码执行💻\", \"type\": \"process\", \"custom\": {\"script\": \"top2\"}, \"moduleConfig\": \"{\\\"type\\\": \\\"code-execution-js\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}, {\\\"id\\\": \\\"error\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"失败\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Fail\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"script\\\", \\\"type\\\": \\\"code-input\\\", \\\"label\\\": \\\"处理脚本\\\", \\\"config\\\": {\\\"theme\\\": \\\"vs-dark\\\", \\\"height\\\": 200, \\\"options\\\": {\\\"minimap\\\": {\\\"enabled\\\": false}, \\\"fontSize\\\": 14, \\\"lineNumbers\\\": true}, \\\"language\\\": \\\"javascript\\\", \\\"defaultValue\\\": \\\"function Filter(msg, metadata, msgType) {  \\\\n  return { msg: msg, metadata: metadata, msgType: msgType };\\\\n}\\\"}}], \\\"runnable\\\": true}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 356, \"dragging\": false, \"position\": {\"x\": 120, \"y\": 15}, \"selected\": true, \"positionAbsolute\": {\"x\": 120, \"y\": 15}}, {\"id\": \"custom-d6772710-ae3c-448f-9b69-8d300132c6ae\", \"data\": {\"name\": \"开始😄\", \"type\": \"input\", \"custom\": {\"param\": \"22222222\"}, \"moduleConfig\": \"{\\\"type\\\": \\\"start\\\", \\\"point\\\": {\\\"inputs\\\": [], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"param\\\", \\\"type\\\": \\\"json-input\\\", \\\"label\\\": \\\"请求头\\\", \\\"config\\\": {\\\"height\\\": 150, \\\"defaultValue\\\": \\\"{\\\\\\\"Accept\\\\\\\": \\\\\\\"application/json\\\\\\\",\\\\\\\"Content-Type\\\\\\\": \\\\\\\"application/json\\\\\\\"}\\\"}}]}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 306, \"dragging\": false, \"position\": {\"x\": 540, \"y\": 15}, \"selected\": false, \"positionAbsolute\": {\"x\": 540, \"y\": 15}}, {\"id\": \"custom-72eda4dc-3f22-40a6-a16d-d04152918406\", \"data\": {\"name\": \"代码执行💻\", \"type\": \"process\", \"custom\": {\"script\": \"top1\"}, \"moduleConfig\": \"{\\\"type\\\": \\\"code-execution-js\\\", \\\"point\\\": {\\\"inputs\\\": [{\\\"id\\\": \\\"input\\\", \\\"type\\\": \\\"target\\\", \\\"label\\\": \\\"输入数据\\\", \\\"position\\\": \\\"left\\\", \\\"handleType\\\": \\\"Data\\\"}], \\\"outputs\\\": [{\\\"id\\\": \\\"success\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"成功\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Success\\\"}, {\\\"id\\\": \\\"error\\\", \\\"type\\\": \\\"source\\\", \\\"label\\\": \\\"失败\\\", \\\"position\\\": \\\"right\\\", \\\"handleType\\\": \\\"Fail\\\"}]}, \\\"fields\\\": [{\\\"id\\\": \\\"script\\\", \\\"type\\\": \\\"code-input\\\", \\\"label\\\": \\\"处理脚本\\\", \\\"config\\\": {\\\"theme\\\": \\\"vs-dark\\\", \\\"height\\\": 200, \\\"options\\\": {\\\"minimap\\\": {\\\"enabled\\\": false}, \\\"fontSize\\\": 14, \\\"lineNumbers\\\": true}, \\\"language\\\": \\\"javascript\\\", \\\"defaultValue\\\": \\\"function Filter(msg, metadata, msgType) {  \\\\n  return { msg: msg, metadata: metadata, msgType: msgType };\\\\n}\\\"}}], \\\"runnable\\\": true}\"}, \"type\": \"custom\", \"width\": 400, \"height\": 356, \"position\": {\"x\": 390, \"y\": -210}}], \"version\": \"1.0\", \"metadata\": {\"edgeCount\": 1, \"nodeCount\": 3, \"lastModified\": \"2024-12-10T03:27:39.358Z\"}, \"exportedAt\": \"2024-12-10T03:27:39.358Z\", \"description\": \"Canvas configuration\"}}','2024-12-09 14:15:32','2024-12-10 11:27:38','system','system');
/*!40000 ALTER TABLE `canvas` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `datasource`
--

DROP TABLE IF EXISTS `datasource`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `datasource` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` varchar(255) NOT NULL,
  `config` json NOT NULL,
  `switch` int NOT NULL,
  `hash` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `datasource`
--

LOCK TABLES `datasource` WRITE;
/*!40000 ALTER TABLE `datasource` DISABLE KEYS */;
INSERT INTO `datasource` VALUES (3,'mysq','{\"sss\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'','','2024-12-09 18:28:57','2024-12-09 19:14:13'),(4,'mysql','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 18:28:58','2024-12-09 18:28:58'),(5,'mysql','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 18:28:58','2024-12-09 18:28:58'),(6,'mysql','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 18:28:58','2024-12-09 18:28:58'),(7,'mysql','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 18:28:59','2024-12-09 18:28:59'),(8,'pg2','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 18:29:24','2024-12-09 18:29:24'),(13,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 19:15:28','2024-12-09 19:15:28'),(14,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"1234563\", \"username\": \"root\"}',1,'','','2024-12-09 19:15:28','2024-12-09 20:17:49'),(15,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 19:15:29','2024-12-09 19:15:29'),(16,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',1,'abcd1234','enable','2024-12-09 19:15:29','2024-12-09 19:15:29'),(17,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',0,'','','2024-12-09 19:15:30','2024-12-09 20:30:54'),(18,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',0,'','','2024-12-09 19:15:30','2024-12-09 20:30:44'),(19,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',0,'','','2024-12-09 19:15:30','2024-12-09 20:30:37'),(20,'pg','{\"host\": \"localhost\", \"port\": 3306, \"database\": \"test\", \"password\": \"123456\", \"username\": \"root\"}',0,'','','2024-12-09 19:15:31','2024-12-09 20:30:32'),(23,'pg','{}',1,'','','2024-12-09 20:18:43','2024-12-09 20:18:43'),(24,'mysql','{}',0,'','','2024-12-09 20:18:54','2024-12-09 20:30:20');
/*!40000 ALTER TABLE `datasource` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gogogo_kv`
--

DROP TABLE IF EXISTS `gogogo_kv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gogogo_kv` (
  `id` int NOT NULL AUTO_INCREMENT,
  `spider_name` varchar(255) NOT NULL,
  `k` varchar(255) NOT NULL,
  `v` text NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gogogo_kv`
--

LOCK TABLES `gogogo_kv` WRITE;
/*!40000 ALTER TABLE `gogogo_kv` DISABLE KEYS */;
/*!40000 ALTER TABLE `gogogo_kv` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `locks`
--

DROP TABLE IF EXISTS `locks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `locks` (
  `lock_name` varchar(255) NOT NULL COMMENT '锁名称',
  `is_locked` tinyint(1) NOT NULL COMMENT '锁是否被持有',
  `held_by` varchar(255) NOT NULL COMMENT '锁持有者',
  `locked_time` datetime NOT NULL COMMENT '锁开始持有时间',
  `timeout` int NOT NULL COMMENT '锁超时时间（秒）',
  `updated_time` datetime NOT NULL COMMENT '锁更新时间',
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locks`
--

LOCK TABLES `locks` WRITE;
/*!40000 ALTER TABLE `locks` DISABLE KEYS */;
INSERT INTO `locks` VALUES ('datasource_client_check_lock',0,'','2024-12-10 19:37:40',10,'2024-12-10 11:37:40',3);
/*!40000 ALTER TABLE `locks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `module`
--

DROP TABLE IF EXISTS `module`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `module` (
  `module_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '组件ID',
  `module_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '组件名称',
  `module_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '组件类型',
  `module_config` json NOT NULL COMMENT '组件配置',
  `module_index` int NOT NULL COMMENT '排序desc',
  PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='组件库';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `module`
--

LOCK TABLES `module` WRITE;
/*!40000 ALTER TABLE `module` DISABLE KEYS */;
INSERT INTO `module` VALUES ('0def1c72-a83f-43a6-b6b1-c4a4c589d16b','代码执行💻','process','{\"type\": \"code-execution-js\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"script\", \"type\": \"code-input\", \"label\": \"处理脚本\", \"config\": {\"theme\": \"vs-dark\", \"height\": 200, \"options\": {\"minimap\": {\"enabled\": false}, \"fontSize\": 14, \"lineNumbers\": true}, \"language\": \"javascript\", \"defaultValue\": \"function Filter(msg, metadata, msgType) {  \\n  return { msg: msg, metadata: metadata, msgType: msgType };\\n}\"}}], \"runnable\": true}',21111111),('0e36fd17-1f91-44d7-b124-346194e7f031','开始😄','input','{\"type\": \"start\", \"point\": {\"inputs\": [], \"outputs\": [{\"id\": \"success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}]}, \"fields\": [{\"id\": \"param\", \"type\": \"json-input\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": \"{\\\"Accept\\\": \\\"application/json\\\",\\\"Content-Type\\\": \\\"application/json\\\"}\"}}]}',1),('3c2a2245-480f-4ffb-871e-11b1389a27bf','结束🩷','output','{\"type\": \"end\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}]}, \"fields\": []}',2),('9209e3e2-6cd6-4cd5-b0e7-b448fb06a88e','丰富的组件🌹','process','{\"type\": \"start\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"success\", \"type\": \"source\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"error\", \"type\": \"source\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Fail\"}]}, \"fields\": [{\"id\": \"mode\", \"type\": \"radio\", \"label\": \"运行模式\", \"config\": {\"options\": [{\"label\": \"同步\", \"value\": \"sync\"}, {\"label\": \"异步\", \"value\": \"async\"}], \"defaultValue\": \"sync\"}}, {\"id\": \"features\", \"type\": \"checkbox\", \"label\": \"启用功能\", \"config\": {\"options\": [{\"label\": \"错误重试\", \"value\": \"retry\"}, {\"label\": \"日志记录\", \"value\": \"logging\"}, {\"label\": \"数据缓存\", \"value\": \"cache\"}], \"defaultValue\": [\"logging\"]}}, {\"id\": \"retryCount\", \"type\": \"slider\", \"label\": \"重试次数\", \"config\": {\"max\": 10, \"min\": 0, \"step\": 1, \"marks\": {\"0\": \"0\", \"5\": \"5\", \"10\": \"10\"}, \"defaultValue\": 3}}, {\"id\": \"timeout\", \"type\": \"number\", \"label\": \"超时时间(ms)\", \"config\": {\"max\": 10000, \"min\": 1000, \"step\": 1000, \"defaultValue\": 5000}}, {\"id\": \"headers\", \"type\": \"json-input\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": {\"Accept\": \"application/json\", \"Content-Type\": \"application/json\"}}}, {\"id\": \"script\", \"type\": \"code-input\", \"label\": \"处理脚本\", \"config\": {\"theme\": \"tomorrow\", \"height\": 200, \"language\": \"javascript\", \"defaultValue\": \"//请求发送前的数据处理\\nfunction preprocess(data) {\\n  return data;\\n}\"}}, {\"id\": \"dataSource\", \"type\": \"select\", \"label\": \"数据源\", \"config\": {\"mode\": \"multiple\", \"options\": [{\"label\": \"MySQL\", \"value\": \"mysql\"}, {\"label\": \"MongoDB\", \"value\": \"mongodb\"}, {\"label\": \"Redis\", \"value\": \"redis\"}], \"defaultValue\": []}}, {\"id\": \"priority\", \"type\": \"select\", \"label\": \"优先级\", \"config\": {\"options\": [{\"label\": \"高\", \"value\": \"high\"}, {\"label\": \"中\", \"value\": \"medium\"}, {\"label\": \"低\", \"value\": \"low\"}], \"defaultValue\": \"medium\"}}], \"runnable\": true}',3);
/*!40000 ALTER TABLE `module` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `module_relation`
--

DROP TABLE IF EXISTS `module_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `module_relation` (
  `id` int NOT NULL AUTO_INCREMENT,
  `module_id` varchar(255) NOT NULL,
  `goal_id` varchar(255) NOT NULL,
  `types` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `nidx_module_id` (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `module_relation`
--

LOCK TABLES `module_relation` WRITE;
/*!40000 ALTER TABLE `module_relation` DISABLE KEYS */;
/*!40000 ALTER TABLE `module_relation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `space_record`
--

DROP TABLE IF EXISTS `space_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `space_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '空间 ID',
  `status` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '运行状态',
  `serial_number` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '流水号',
  `run_time` datetime NOT NULL COMMENT '运行开始时间',
  `record_name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '运行记录名称',
  `duration` int NOT NULL COMMENT '耗时 ms',
  `other` json NOT NULL COMMENT '其他配置',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_serial_number` (`serial_number`),
  KEY `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `space_record`
--

LOCK TABLES `space_record` WRITE;
/*!40000 ALTER TABLE `space_record` DISABLE KEYS */;
INSERT INTO `space_record` VALUES (1,'asdasd','success','1111','2024-11-12 00:00:00','安师大实打实的',100,'{\"sss\": \"sssss\"}'),(2,'asdasd','fail','22222','2024-11-12 00:00:00','安师大实打实的 1212',100,'{\"sss\": \"sssss\"}');
/*!40000 ALTER TABLE `space_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `trace`
--

DROP TABLE IF EXISTS `trace`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `trace` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) NOT NULL COMMENT '空间 ID',
  `trace_id` varchar(255) NOT NULL COMMENT '追踪 ID',
  `input` json NOT NULL COMMENT '组件输入',
  `logic` json NOT NULL COMMENT '执行逻辑',
  `output` json NOT NULL COMMENT '组件输出',
  `step` int NOT NULL COMMENT '分步',
  `node_id` varchar(255) NOT NULL COMMENT '节点 ID',
  `node_name` varchar(255) NOT NULL COMMENT '节点名称',
  `status` varchar(255) NOT NULL COMMENT '运行状态',
  `elapsed_time` int NOT NULL COMMENT '运行耗时',
  `start_time` datetime NOT NULL COMMENT '执行时间',
  `error_msg` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `unidx_workspace_id` (`workspace_id`),
  KEY `unidx_trace_id` (`trace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='组件链路追踪记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `trace`
--

LOCK TABLES `trace` WRITE;
/*!40000 ALTER TABLE `trace` DISABLE KEYS */;
INSERT INTO `trace` VALUES (1,'asdasd','22222','{\"input\": \"xx\"}','{\"logic\": \"xx\"}','{\"output\": \"xx\"}',1,'123','1234','success',100,'2024-12-11 00:00:00','没有错误');
/*!40000 ALTER TABLE `trace` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'The username',
  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'The user password',
  `mobile` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'The mobile phone number',
  `gender` char(10) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'male' COMMENT 'gender,male|female|unknown',
  `nickname` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'The nickname',
  `type` tinyint(1) DEFAULT '0' COMMENT 'The user type, 0:normal,1:vip, for test golang keyword',
  `create_at` timestamp NULL DEFAULT NULL,
  `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mobile_index` (`mobile`),
  UNIQUE KEY `name_index` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='user table';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace`
--

DROP TABLE IF EXISTS `workspace`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workspace` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主建',
  `workspace_id` varchar(255) NOT NULL COMMENT '主建',
  `workspace_name` varchar(255) NOT NULL COMMENT '名称',
  `workspace_desc` text COMMENT '描述',
  `workspace_type` varchar(50) DEFAULT NULL COMMENT '类型flow|agent',
  `workspace_icon` varchar(255) DEFAULT NULL COMMENT 'iconUrl',
  `canvas_config` text COMMENT '前端画布配置',
  `configuration` json NOT NULL COMMENT '配置信息 KV',
  `additionalInfo` json NOT NULL COMMENT '扩展信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  `is_delete` int NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_workspace_id` (`workspace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='工作空间表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace`
--

LOCK TABLES `workspace` WRITE;
/*!40000 ALTER TABLE `workspace` DISABLE KEYS */;
INSERT INTO `workspace` VALUES (17,'ct1i8ud3sjtg3kkhmvq0','雪兔 17','雪兔 17雪兔 17雪兔 17','workflow','https://api.iconify.design/ant-design:appstore-outlined.svg',NULL,'{}','{}','2024-11-24 21:04:57','2024-11-24 21:04:57',0),(18,'ct1l3jl3sjthu24kms8g','雪兔 18','1','workflow','https://api.iconify.design/ant-design:appstore-outlined.svg',NULL,'{}','{}','2024-11-25 00:18:22','2024-11-25 00:18:22',1),(19,'ct4qi1t3sjtp5jbouq0g','雪兔 18','111','workflow','https://api.iconify.design/ant-design:appstore-outlined.svg',NULL,'{}','{}','2024-11-29 19:44:08','2024-11-29 19:44:08',0),(20,'ct9dq9d3sjtl7g6q5ha0','雪兔 20','雪兔 20雪兔 20雪兔 20雪兔 20','workflow','https://api.iconify.design/ant-design:appstore-outlined.svg',NULL,'{}','{}','2024-12-06 19:16:53','2024-12-06 19:16:53',1),(21,'ctb8m0t3sjtq9abeusm0','画布 21','画布 21','workflow','https://api.iconify.design/ant-design:appstore-outlined.svg',NULL,'{}','{}','2024-12-09 14:15:32','2024-12-09 14:15:32',0);
/*!40000 ALTER TABLE `workspace` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace_tag`
--

DROP TABLE IF EXISTS `workspace_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workspace_tag` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主建',
  `tag_name` varchar(255) NOT NULL COMMENT '标签名称',
  `is_delete` int NOT NULL COMMENT '逻辑删除',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='标签表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace_tag`
--

LOCK TABLES `workspace_tag` WRITE;
/*!40000 ALTER TABLE `workspace_tag` DISABLE KEYS */;
INSERT INTO `workspace_tag` VALUES (1,'http',0,'2024-11-23 23:06:42','2024-11-23 23:06:42'),(2,'并行',0,'2024-11-23 23:06:42','2024-11-23 23:06:42'),(3,'ccc',0,'2024-11-23 23:34:39','2024-11-23 23:34:39'),(4,'mysql',0,'2024-11-23 23:35:10','2024-11-23 23:35:10'),(5,'111',0,'2024-11-24 01:22:05','2024-11-24 01:22:05'),(6,'1222',0,'2024-11-24 01:22:05','2024-11-24 01:22:05'),(7,'112231231',0,'2024-11-24 01:26:01','2024-11-24 01:26:01'),(8,'12312313333333333333',0,'2024-11-24 01:26:01','2024-11-24 01:26:01'),(9,'3333333',0,'2024-11-24 01:26:01','2024-11-24 01:26:01'),(10,'field \"workSpaceType\" is not set',0,'2024-11-24 01:27:39','2024-11-24 01:27:39'),(11,'123 领导撒放假了就开始地方；1',0,'2024-11-24 01:27:39','2024-11-24 01:27:39'),(12,'雪兔 16',0,'2024-11-24 01:33:25','2024-11-24 01:33:25'),(13,'雪兔 17',0,'2024-11-24 01:33:25','2024-11-24 01:33:25'),(14,'雪兔 18',0,'2024-11-24 01:33:25','2024-11-24 01:33:25'),(15,'雪兔 19',0,'2024-11-24 01:33:25','2024-11-24 01:33:25'),(16,'雪兔 20',0,'2024-11-24 01:33:25','2024-11-24 01:33:25'),(17,'new',0,'2024-11-24 21:04:57','2024-11-24 21:04:57'),(18,'刷卡',0,'2024-11-29 19:44:08','2024-11-29 19:44:08');
/*!40000 ALTER TABLE `workspace_tag` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace_tag_mapping`
--

DROP TABLE IF EXISTS `workspace_tag_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workspace_tag_mapping` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主建',
  `tag_id` int NOT NULL COMMENT '标签ID',
  `workspace_id` varchar(255) NOT NULL COMMENT '画布空间ID',
  PRIMARY KEY (`id`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_worlspace_id` (`workspace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='画布标签映射表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace_tag_mapping`
--

LOCK TABLES `workspace_tag_mapping` WRITE;
/*!40000 ALTER TABLE `workspace_tag_mapping` DISABLE KEYS */;
INSERT INTO `workspace_tag_mapping` VALUES (1,1,'ct0uv0d3sjtst3gv01k0'),(2,2,'ct0uv0d3sjtst3gv01k0'),(24,5,'ct110ad3sjtubcqbgt30'),(25,7,'ct110ad3sjtubcqbgt30'),(26,8,'ct110ad3sjtubcqbgt30'),(27,9,'ct110ad3sjtubcqbgt30'),(28,10,'ct110ad3sjtubcqbgt30'),(29,11,'ct110ad3sjtubcqbgt30'),(30,5,'ct10ufd3sjtubcqbgt2g'),(31,6,'ct10ufd3sjtubcqbgt2g'),(32,12,'ct113p53sjtubcqbgt3g'),(33,13,'ct113p53sjtubcqbgt3g'),(34,14,'ct113p53sjtubcqbgt3g'),(35,15,'ct113p53sjtubcqbgt3g'),(36,16,'ct113p53sjtubcqbgt3g'),(37,17,'ct1i8ud3sjtg3kkhmvq0'),(39,17,'ct4qi1t3sjtp5jbouq0g'),(40,18,'ct4qi1t3sjtp5jbouq0g'),(42,12,'ctb8m0t3sjtq9abeusm0');
/*!40000 ALTER TABLE `workspace_tag_mapping` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-10 11:37:48
