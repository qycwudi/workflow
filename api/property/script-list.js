import http from 'k6/http';

// export const options = {
//   vus: 10,
//   duration: '30s',
// };

export default function () {
  const url = 'http://198.19.249.3:8899/workflow/canvas/run/record';
  const payload = JSON.stringify({
    "id": "cs34put3sjtg42ghcf8g"
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post(url, payload, params);
}
