docker run --rm -i grafana/k6 run - <script-api.js

docker run --rm -i grafana/k6 run --vus 10 --duration 30s - <script.js