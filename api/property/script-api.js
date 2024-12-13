import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 100, // 虚拟用户数
  duration: '1m', // 持续时间
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95%的请求响应时间小于500ms
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
  const url = 'http://198.19.249.3:9999/api/role/v1/ctdpe7l3sjtsoutieu20';
  const payload = JSON.stringify({
    "name": randomName(),
    "age": randomAge()
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1ea00ca19d4f47b48c9ed8aec131674d'
    },
  };

  const response = http.post(url, payload, params);
  
  // 检查响应状态
  check(response, {
    'is status 200': (r) => r.status === 200,
    'transaction time < 500ms': (r) => r.timings.duration < 500
  });

  // 打印响应内容
  console.log('Response status:', response.status);
  console.log('Response body:', response.body);
  console.log('Response headers:', response.headers);
}


//      ✓ is status 200
//      ✗ transaction time < 500ms
//       ↳  2% — ✓ 87 / ✗ 3342

//      checks.........................: 51.26% 3516 out of 6858
//      data_received..................: 485 kB 7.7 kB/s
//      data_sent......................: 858 kB 14 kB/s
//      http_req_blocked...............: avg=212.19µs min=917ns    med=3.83µs  max=15.16ms p(90)=10.63µs  p(95)=38.49µs 
//      http_req_connecting............: avg=186.9µs  min=0s       med=0s      max=13.3ms  p(90)=0s       p(95)=0s      
//    ✗ http_req_duration..............: avg=1.78s    min=219.42ms med=1.39s   max=3.62s   p(90)=3.23s    p(95)=3.36s   
//        { expected_response:true }...: avg=1.78s    min=219.42ms med=1.39s   max=3.62s   p(90)=3.23s    p(95)=3.36s   
//      http_req_failed................: 0.00%  0 out of 3429
//      http_req_receiving.............: avg=102.33µs min=10.29µs  med=69.54µs max=7.78ms  p(90)=156.66µs p(95)=238.37µs
//      http_req_sending...............: avg=79.32µs  min=3.75µs   med=16.87µs max=9.73ms  p(90)=56.5µs   p(95)=127.78µs
//      http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s      p(90)=0s       p(95)=0s      
//      http_req_waiting...............: avg=1.78s    min=210.84ms med=1.39s   max=3.62s   p(90)=3.23s    p(95)=3.36s   
//      http_reqs......................: 3429   54.705711/s
//      iteration_duration.............: avg=1.78s    min=221.65ms med=1.39s   max=3.62s   p(90)=3.23s    p(95)=3.36s   
//      iterations.....................: 3429   54.705711/s
//      vus............................: 23     min=23           max=100
//      vus_max........................: 100    min=100          max=100


// running (1m02.7s), 000/100 VUs, 3429 complete and 0 interrupted iterations
// default ✓ [ 100% ] 100 VUs  1m0s
// time="2024-12-13T02:22:52Z" level=error msg="thresholds on metrics 'http_req_duration' have been crossed"

// 2024-12-13T10:24:06.714+08:00    stat   CPU: 0m, MEMORY: Alloc=3.8Mi, TotalAlloc=575.9Mi, Sys=43.1Mi, NumGC=87  caller=stat/usage.go:61


// 根据你提供的压测结果，我们可以进行以下分析：

// ### 总体情况

// 1. **测试时间**: 测试运行时间为 1 分 2.7 秒，使用了 100 个虚拟用户（VUs）。
// 2. **请求数量**: 总共完成了 3429 次请求，平均每秒处理约 54.7 次请求。
// 3. **成功率**: 所有请求都成功完成，`http_req_failed` 为 0.00%。

// ### 关键指标分析

// 1. **HTTP 状态**:
//    - 你有 200 状态的请求，表示请求成功。
//    - 但是有一个检查失败，`transaction time < 500ms`，这意味着大约 3342 次请求的响应时间超过了 500 毫秒。

// 2. **检查结果**:
//    - `checks`: 通过检查的比例为 51.26%，这表明一半多的请求未能满足所有的性能指标，尤其是响应时间的要求。

// 3. **请求时延**:
//    - `http_req_duration` 平均为 1.78 秒，最大为 3.62 秒，这远远超过了 500 毫秒的阈值。这是导致检查失败的主要原因。
//    - `http_req_waiting`（即服务器处理请求的时间）也显示出相似的情况，平均为 1.78 秒，说明服务器在处理请求时的延迟较高。

// 4. **数据传输**:
//    - `data_received` 和 `data_sent` 表明数据传输量相对较小，但这并没有直接影响响应时间。

// 5. **虚拟用户 (VUs)**:
//    - 在测试期间，最大并发用户数为 100，且始终保持在这个水平，表明测试的负载是稳定的。

// 6. **CPU 和内存使用**:
//    - CPU 使用情况为 0m，内存使用为 Alloc=3.8Mi，TotalAlloc=575.9Mi，Sys=43.1Mi，NumGC=87。这表明应用的内存使用相对较低，可能不会成为瓶颈。

// ### 问题分析

// - **响应时间过长**: 主要问题是请求的响应时间超过了设定的阈值（500ms）。这可能是由于后端服务的性能瓶颈、数据库查询延迟、网络延迟等原因导致的。
// - **检查失败**: 由于大部分请求未能在 500ms 内完成，导致检查失败的比例较高。

// ### 建议

// 1. **性能调优**:
//    - 检查后端服务的性能，分析是否存在瓶颈，例如数据库查询、缓存策略等。
//    - 评估服务器的资源（CPU、内存、IO）是否足够，必要时可以进行扩容。

// 2. **代码优化**:
//    - 审查代码逻辑，尤其是处理请求的部分，看看是否可以优化算法或减少不必要的计算。

// 3. **增加监控**:
//    - 增加对请求处理时间的监控，及时发现并解决问题。

// 4. **重试和超时机制**:
//    - 考虑在客户端实现重试机制和合理的超时设置，以提高用户体验。

// 5. **进行更多测试**:
//    - 进行不同负载情况下的压力测试，以更全面地了解系统性能。

// 通过以上分析和建议，你可以针对性地优化系统性能，提高响应速度，确保在高负载下也能保持良好的性能表现。