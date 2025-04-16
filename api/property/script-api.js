import { check } from 'k6';
import http from 'k6/http';

export const options = {
  vus: 200, 
  duration: '2m',
  thresholds: {
    http_req_duration: ['p(95)<500'], 
  },
};

function randomName() {
  const surnames = ['张', '李', '王', '赵', '刘', '陈', '杨', '黄', '周', '吴'];
  const names = ['伟', '芳', '娜', '秀英', '敏', '静', '丽', '强', '磊', '军'];
  return surnames[Math.floor(Math.random() * surnames.length)] + 
         names[Math.floor(Math.random() * names.length)];
}


function randomAge() {
  return Math.floor(Math.random() * 50) + 18;
}
export default function () {
  const url = 'http://192.168.0.12:8889/api/role/v1/cvvih8frng8rbrjee3gg';
  const payload = JSON.stringify({
    "username":"xuetu",
    "name": randomName(),
    "age": randomAge()
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ecab0d05783949c8a3169df8ab085a9c'
    },
  };

  const response = http.post(url, payload, params);
  
  check(response, {
    'is status 200': (r) => r.status === 200,
    'transaction time < 500ms': (r) => r.timings.duration < 500
  });

  // console.log('Response status:', response.status);
  // console.log('Response body:', response.body);
  // console.log('Response headers:', response.headers);
}
