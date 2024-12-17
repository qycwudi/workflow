import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 50, // 增加虚拟用户数到200以加快插入速度
  iterations: 10000, // 固定迭代次数为10000
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 放宽响应时间限制到2秒
    http_req_failed: ['rate<0.01'], // 失败率需小于1%
  },
};

// 生成随机姓名
function randomName() {
  const surnames = ['张', '李', '王', '赵', '刘', '陈', '杨', '黄', '周', '吴'];
  const names = ['伟', '芳', '娜', '秀英', '敏', '静', '丽', '强', '磊', '军'];
  return surnames[Math.floor(Math.random() * surnames.length)] + 
         names[Math.floor(Math.random() * names.length)];
}

// 生成随机年龄
function randomAge() {
  return Math.floor(Math.random() * 50) + 18; // 18-67岁
}

export default function () {
  const url = 'http://198.19.249.3:9999/api/role/v1/ctgqi9l3sjtsviksh0rg';
  const payload = JSON.stringify({
    "table": "users",
    "name": randomName(),
    "age": randomAge()
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 8537cc19b30c4a49849b6efbd9bb4aa9'
    },
  };

  const response = http.post(url, payload, params);
  
  // 检查响应状态
  check(response, {
    'is status 200': (r) => r.status === 200,
    'transaction time < 2000ms': (r) => r.timings.duration < 2000
  });

  // 记录每100次请求的进度
  if(__ITER % 100 === 0) {
    console.log(`已完成 ${__ITER} 条数据插入`);
  }
}