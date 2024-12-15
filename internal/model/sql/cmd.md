goctl model mysql ddl --src user.sql --dir .. -i ''

brew install ariga/tap/atlas

atlas schema diff \
  --from mysql://root:root@192.168.49.2:31426/wk \
  --to mysql://root:Root@123@10.99.43.9:3306/workflow

atlas schema inspect -u "mysql://root:root@192.168.49.2:31426/wk" > schema.hcl



dump库表
kubectl exec -it xuetu-db-cc774ff4b-pd6pf  -- sh -c 'mysqldump -u root -proot wk' > mydb_dump.sql

mysql -u root -p workflow < dump.sql

#!/bin/bash

# 指定要处理的文件
DUMP_FILE="dump.sql"

# 查找并替换字符集和排序规则
sed -i 's/CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci/CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci/g' "$DUMP_FILE"
sed -i 's/COLLATE=utf8mb4_0900_ai_ci/COLLATE=utf8mb4_general_ci/g' "$DUMP_FILE"

echo "替换完成！"

chmod +x replace_collation.sh