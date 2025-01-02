<!-- 导出 -->
atlas schema inspect -u "mysql://root:root@14.103.249.105:30006/workflow_prod" --format '{{ sql . }}' > workflow_prod_schema.sql

<!-- 导入 -->
atlas schema apply \
  --url "mysql://root:root@14.103.249.105:30006/workflow_prod" \
  --to "file://workflow_schema.sql" \
  --dev-url "mysql://root:root@14.103.249.105:30006/atlas_database"

atlas schema apply \
  --url "mysql://root:root@14.103.249.105:30006/workflow_dev" \
  --to "file://workflow_schema.sql" \
  --dev-url "mysql://root:root@14.103.249.105:30006/atlas_database"

  atlas schema apply \
  --url "mysql://root:Root@123@10.99.29.9:3306/workflow_test" \
  --to "file://workflow_schema.sql" \
  --dev-url "mysql://root:root@14.103.249.105:30006/atlas_database"


<!-- 对比 -->
atlas schema diff \
  --from "mysql://root:root@14.103.249.105:30006/workflow_dev" \
  --to "mysql://root:root@14.103.249.105:30006/workflow_prod" \
  --format '{{ sql . "  " }}' > workflow_diff.sql


  atlas schema diff \
  --from "mysql://root:root@14.103.249.105:30006/workflow_dev" \
  --to "mysql://root:Root@123@10.99.29.9:3306/wkflow" \
  --format '{{ sql . "  " }}' > workflow_diff.sql