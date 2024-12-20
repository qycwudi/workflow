package rulego

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jlaffaye/ftp"
	"github.com/rulego/rulego/api/types"
)

func TestSftpNode_executeFtp(t *testing.T) {
	// 创建测试文件
	testContent := []byte("test content")
	err := os.WriteFile("./testdata/test.txt", testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove("./testdata/test.txt")

	// 测试 SFTP
	msg := types.RuleMsg{
		Data: `{
			"action": "upload",
			"config": {
				"protocol": "sftp",
				"host": "10.99.169.7",
				"port": 2233,
				"username": "beuser",
				"password": "Bepassword@123"
			},
			"srcPath": "./testdata/test.txt",
			"destPath": "/tmp/test.txt"
		}`,
	}

	node := &FileServerNode{}

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行SFTP操作失败: %v", err)
	}

	msg.Data = `{
		"action": "download",
		"config": {
			"protocol": "sftp",
			"host": "10.99.169.7",
			"port": 2233,
			"username": "beuser", 
			"password": "Bepassword@123"
		},
		"srcPath": "/tmp/test.txt",
		"destPath": "./testdata/test_download.txt"
	}`

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行SFTP下载操作失败: %v", err)
	}
	defer os.Remove("./testdata/test_download.txt")

	msg.Data = `{
		"action": "delete",
		"config": {
			"protocol": "sftp",
			"host": "10.99.169.7",
			"port": 2233,
			"username": "beuser",
			"password": "Bepassword@123"
		},
		"path": "/tmp/test.txt"
	}`

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行SFTP删除操作失败: %v", err)
	}
}

/*
	---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations: {}
  labels:
    k8s.kuboard.cn/name: ido-ftp
  name: ido-ftp
  namespace: uat-43
  resourceVersion: '501817066'
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      k8s.kuboard.cn/name: ido-ftp
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s.kuboard.cn/name: ido-ftp
    spec:
      containers:
        - env:
            - name: FTP_USER
              value: test
            - name: FTP_PASS
              value: test
            - name: PASV_ADDRESS
              value: 10.99.43.43
            - name: PASV_MIN_PORT
              value: '21100'
            - name: PASV_MAX_PORT
              value: '21110'
          image: '10.12.0.78:5000/k8s/vsftpd'
          imagePullPolicy: Always
          name: vsftpd
          ports:
            - containerPort: 21
              protocol: TCP
            - containerPort: 20
              protocol: TCP
            - containerPort: 21100
              protocol: TCP
            - containerPort: 21101
              protocol: TCP
            - containerPort: 21102
              protocol: TCP
            - containerPort: 21103
              protocol: TCP
            - containerPort: 21104
              protocol: TCP
            - containerPort: 21105
              protocol: TCP
            - containerPort: 21106
              protocol: TCP
            - containerPort: 21107
              protocol: TCP
            - containerPort: 21108
              protocol: TCP
            - containerPort: 21109
              protocol: TCP
            - containerPort: 21110
              protocol: TCP
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 1
  conditions:
    - lastTransitionTime: '2024-12-18T15:36:54Z'
      lastUpdateTime: '2024-12-18T15:36:54Z'
      message: Deployment has minimum availability.
      reason: MinimumReplicasAvailable
      status: 'True'
      type: Available
    - lastTransitionTime: '2024-12-18T15:25:27Z'
      lastUpdateTime: '2024-12-18T15:36:54Z'
      message: ReplicaSet "ido-ftp-85c5758748" has successfully progressed.
      reason: NewReplicaSetAvailable
      status: 'True'
      type: Progressing
  observedGeneration: 16
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1

---
apiVersion: v1
kind: Service
metadata:
  annotations: {}
  name: ido-ftp
  namespace: uat-43
  resourceVersion: '501812860'
spec:
  clusterIP: 10.99.113.114
  clusterIPs:
    - 10.99.113.114
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  ports:
    - name: ftp
      nodePort: 30013
      port: 21
      protocol: TCP
      targetPort: 21
    - name: ftp-data
      nodePort: 30001
      port: 20
      protocol: TCP
      targetPort: 20
    - name: pasv
      nodePort: 30002
      port: 21100
      protocol: TCP
      targetPort: 21100
    - name: pasv-20001
      nodePort: 30003
      port: 21101
      protocol: TCP
      targetPort: 21101
    - name: pasv-20002
      nodePort: 30014
      port: 21102
      protocol: TCP
      targetPort: 21102
    - name: pasv-20003
      nodePort: 30015
      port: 21103
      protocol: TCP
      targetPort: 21103
    - name: pasv-20004
      nodePort: 30016
      port: 21104
      protocol: TCP
      targetPort: 21104
    - name: pasv-20005
      nodePort: 30007
      port: 21105
      protocol: TCP
      targetPort: 21105
    - name: pasv-20006
      nodePort: 30008
      port: 21106
      protocol: TCP
      targetPort: 21106
    - name: pasv-20007
      nodePort: 30009
      port: 21107
      protocol: TCP
      targetPort: 21107
    - name: pasv-20008
      nodePort: 30010
      port: 21108
      protocol: TCP
      targetPort: 21108
    - name: pasv-20009
      nodePort: 30011
      port: 21109
      protocol: TCP
      targetPort: 21109
    - name: pasv-20010
      nodePort: 30012
      port: 21110
      protocol: TCP
      targetPort: 21110
  selector:
    k8s.kuboard.cn/name: ido-ftp
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}


*/

func TestFtpNode_executeFtp(t *testing.T) {
	// 创建测试文件
	testContent := []byte("test content")
	err := os.WriteFile("./testdata/test.txt", testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove("./testdata/test.txt")

	// 测试 FTP
	msg := types.RuleMsg{
		Data: `{"action":"upload","config":{"protocol":"ftp","host":"10.99.113.114","port":21,"username":"test","password":"test","passive":true},"srcPath":"./testdata/test.txt","destPath":"/test.txt"}`,
	}

	node := &FileServerNode{}

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行FTP上传操作失败: %v", err)
	}

	// msg.Data = `{
	// 	"action": "download",
	// 	"config": {
	// 		"protocol": "ftp",
	// 		"host": "10.99.113.114",
	// 		"port": 21,
	// 		"username": "test",
	// 		"password": "test",
	// 		"passive": true
	// 	},
	// 	"srcPath": "/test.txt",
	// 	"destPath": "./testdata/test_download_ftp.txt"
	// }`

	// err = node.executeFtp(msg)
	// if err != nil {
	// 	t.Errorf("执行FTP下载操作失败: %v", err)
	// }
	// defer os.Remove("./testdata/test_download_ftp.txt")

	// msg.Data = `{
	// 	"action": "delete",
	// 	"config": {
	// 		"protocol": "ftp",
	// 		"host": "10.99.113.114",
	// 		"port": 21,
	// 		"username": "test",
	// 		"password": "test",
	// 		"passive": true
	// 	},
	// 	"path": "/test.txt"
	// }`

	// err = node.executeFtp(msg)
	// if err != nil {
	// 	t.Errorf("执行FTP删除操作失败: %v", err)
	// }
}
func TestFtpNode_net_executeFtp(t *testing.T) {
	ftpServer := "10.99.113.114:21"
	username := "test"
	password := "test"

	// 连接到 FTP 服务器
	c, err := ftp.Dial(ftpServer)
	if err != nil {
		log.Fatalf("Failed to connect to FTP server: %v", err)
	}

	// 登录
	err = c.Login(username, password)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer c.Logout()

	// 列出目录
	entries, err := c.List("")
	if err != nil {
		log.Fatalf("Failed to list directory: %v", err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}
