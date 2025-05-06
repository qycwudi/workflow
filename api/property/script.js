import { check } from 'k6';
import http from 'k6/http';

export const options = {
  vus: 50,  // 并发用户数
  duration: '30s',  // 测试持续时间
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% 的请求应该在 500ms 内完成
    http_req_failed: ['rate<0.01'],    // 请求失败率应该低于 1%
  },
};

export default function () {
  const url = 'http://198.19.249.3:8888/workflow/canvas/run';
  
  const payload = JSON.stringify({
    "id": "d03l4kl3sjtgm7gc9su0",
    "params": {
      "name": "d03l4kl3sjtgm7gc9su0",
      "auth": "Bearer ecab0d05783949c8a3169df8ab085a9c",
      "age": 18,
      "age1": 18,
      "info": {
        "city": "hangzhou",
        "country": "china",
        "height": 170.5,
        "weight": 65.5,
        "phone": 17635800128,
        "email": ["zhangsan@example.com", "lisi@example.com"],
        "like": {
          "food": ["apple", "banana"],
          "sport": ["basketball", "football"]
        }
      },
      "address": ["beijing", "shanghai", "hangzhou"],
      "mritalStatus": false,
      "sight": 1.0
    }
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYwMDQzNjksImlhdCI6MTc0NTkxNzk2OSwidXNlcklkIjoxfQ.UgfG-xGwAYpqt6W38J0OIC-xa0RAmqp7_Tnlhwbz5lg',
      'User-Agent': 'Apifox/1.0.0 (https://apifox.com)'
    },
  };

  const response = http.post(url, payload, params);

  // 检查响应状态码是否为 200
  check(response, {
    'status is 200': (r) => r.status === 200,
  });
  // 打印请求结果
  console.log(response.json());

  // 在每次请求之间添加随机延迟（0.5-1.5秒）
  // sleep(Math.random() * 1 + 0.5);
} 