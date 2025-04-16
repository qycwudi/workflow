import http from 'k6/http';

export const options = {
  vus: 100,
  duration: '30s',
};

export default function () {
  const url = 'http://36.139.231.30:89/janus/invoke/v1';
  const payload = JSON.stringify({
    "method": "gaia.openapi.data-serving.dataService.invoke.product.appId",
    "content": {
      "param": {
        "appId": "4aecc6379655b95b7a786acee370ddbe",
        "productId": "2025041510169100000045",
        "inputs": {
          "taxpayer_id": "91440400MA7LQ1A4XD"
        }
      }
    }
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'token': 'MoKwA8zA7XYRJmW7hUTgv1P+HfRrk88BGdCdAa+xHeoAqECy+49vTFYh+teAfPNTOPrvyKHZAnr8W3b+zDewKwGI+vZaOjnBY+07Y3FeRmNmuXCR+o5I/FGiPVMY+VOB78cfmR/O21s4oIxpVvVH1fGBjWzjQKhvRaw/b+1S/3eTADbs7pBHNIerWy42be9c8LJtD1JSdloqFaFAolkV69N83bv0H3lC2ypthAtD6r3jo9mhplPT4c9eRpQzW2AwgAxiePFNFrgqJ40HHqwZsKwFGIwHVKDm96+TyNUdfcxN/sjqk5b62rV96KfBChd0U06UEuCrJi9/KEIxVSbCrQ=='
    },
  };

  const response = http.post(url, payload, params);
  // console.log(`状态码: ${response.status}`);
  // console.log(`响应体: ${response.body}`);
}
