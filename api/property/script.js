import http from 'k6/http';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  const url = 'http://198.19.249.3:9999/api/role/v1/cs6c0il3sjtobc4avatg';
  const payload = JSON.stringify({
    "name": "张三",
    "age": 18
});

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post(url, payload, params);
}
