import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

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

// 模拟取消请求的场景
function simulateCancellation() {
  // 30%的概率取消请求
  return Math.random() < 0.3;
}

// 模拟网络错误
function simulateNetworkError() {
  // 20%的概率发生网络错误
  return Math.random() < 0.2;
}

// 模拟请求超时
function simulateTimeout() {
  // 10%的概率发生超时
  return Math.random() < 0.1;
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
    timeout: '2s' // 设置较短的超时时间以便测试超时场景
  };

  // 模拟请求取消场景
  if (simulateCancellation()) {
    console.log('请求被用户主动取消');
    return;
  }

  // 模拟网络错误
  if (simulateNetworkError()) {
    console.error('网络连接错误');
    return;
  }

  // 模拟超时
  if (simulateTimeout()) {
    sleep(3); // 故意睡眠3秒触发超时
    console.error('请求超时');
    return;
  }

  try {
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

    // 模拟随机网络延迟 100-1000ms
    sleep(Math.random() * 0.9 + 0.1);
  } catch (error) {
    console.error('请求执行失败:', error.message);
  }
}