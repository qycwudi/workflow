import http from 'k6/http';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const requestDuration = new Trend('request_duration');

export const options = {
  // vus: 50,  // 并发用户数
  // duration: '30s',  // 测试持续时间
  stages: [
    { duration: '1m', target: 50 },  // 1分钟内逐步增加到50个并发用户
    { duration: '3m', target: 50 },  // 保持50个并发用户3分钟
    { duration: '1m', target: 0 },   // 1分钟内逐步减少到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95%的请求应该在500ms内完成
    http_req_failed: ['rate<0.01'],    // 请求失败率应该低于1%
    errors: ['rate<0.01'],             // 错误率应该低于1%
  },
};

// 生成随机测试数据
function generateTestData() {
  const cities = ['hangzhou', 'beijing', 'shanghai', 'guangzhou', 'shenzhen'];
  const countries = ['china', 'usa', 'uk', 'japan', 'germany'];
  const foods = ['apple', 'banana', 'orange', 'grape', 'watermelon'];
  const sports = ['basketball', 'football', 'tennis', 'swimming', 'running'];
  
  return {
    name: `test_${Math.random().toString(36).substring(7)}`,
    auth: "Bearer ecab0d05783949c8a3169df8ab085a9c",
    age: Math.floor(Math.random() * 50) + 18,
    age1: Math.floor(Math.random() * 50) + 18,
    info: {
      city: cities[Math.floor(Math.random() * cities.length)],
      country: countries[Math.floor(Math.random() * countries.length)],
      height: 150 + Math.random() * 50,
      weight: 45 + Math.random() * 30,
      phone: 17600000000 + Math.floor(Math.random() * 100000000),
      email: [
        `user${Math.random().toString(36).substring(7)}@example.com`,
        `user${Math.random().toString(36).substring(7)}@example.com`
      ],
      like: {
        food: foods.sort(() => 0.5 - Math.random()).slice(0, 2),
        sport: sports.sort(() => 0.5 - Math.random()).slice(0, 2)
      }
    },
    address: cities.sort(() => 0.5 - Math.random()).slice(0, 3),
    mritalStatus: Math.random() > 0.5,
    sight: 0.5 + Math.random() * 1.5
  };
}

export default function () {
  const url = 'http://198.19.249.3:8888/workflow/api/v1/d0grk5t3sjtih9ajpuag';
  const payload = JSON.stringify(generateTestData());
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': 'Apifox/1.0.0 (https://apifox.com)'
    },
  };
  // console.log(payload)
  const startTime = new Date();
  const response = http.post(url, payload, params);
  const endTime = new Date();
  
  // 记录请求持续时间
  requestDuration.add(endTime - startTime);

  // 检查响应状态是否为200
  const statusCheck = response.status === 200;
  // 检查响应是否包含数据
  const dataCheck = response.json() !== null;
  // 检查响应时间是否在500毫秒以内
  const timeCheck = response.timings.duration < 500;
  // 组合检查结果
  const checks = statusCheck && dataCheck && timeCheck;

  // 记录错误
  if (!checks) {
    errorRate.add(1);
    let errorMessages = [];
    if (!statusCheck) errorMessages.push(`状态码错误: ${response.status}`);
    if (!dataCheck) errorMessages.push('响应数据为空');
    if (!timeCheck) errorMessages.push(`响应时间过长: ${response.timings.duration}毫秒`);
    console.error(`Request failed: ${errorMessages.join(', ')}`);
  }

  // 添加随机延迟，模拟真实用户行为
  // sleep(Math.random() * 2 + 1); // 1-3秒的随机延迟
} 