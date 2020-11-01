Запуск локального редиса в докере:
```
docker run --name some-redis -p 6379:6379 -d redis
```

Запуск клиента командной строки редиса:
```
docker exec -it some-redis /usr/local/bin/redis-cli
```

Вывод `cache-service`:
```
% PORT=:50051 go run cmd/cache-service/main.go
2020/11/01 22:25:39 Starting server on port :50051.
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.github.com
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.facebook.com
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.google.com
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.bbc.co.uk
2020/11/01 22:25:44 [DEBUG] cache miss: https://golang.org
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.twitter.com
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.duckduckgo.com
2020/11/01 22:25:44 [DEBUG] cache miss: https://www.gitlab.com

```

Вывод `consumer-service`:
```
% SERVER_ADDRESS=localhost:50051 go run cmd/consumer-service/main.go
2020/11/01 22:25:44 [1] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [5] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [5] Received stream 2 (15738 bytes).
2020/11/01 22:25:44 [27] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [30] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [11] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [63] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [30] Received stream 2 (15738 bytes).
2020/11/01 22:25:44 [64] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [9] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [15] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [66] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [38] Received stream 1 (15738 bytes).
2020/11/01 22:25:44 [14] Received stream 1 (15738 bytes).

...

2020/11/01 22:25:56 [941] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [905] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [944] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [913] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [953] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [954] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [964] Received stream 3 (334185 bytes).
2020/11/01 22:25:56 [992] Received stream 2 (334185 bytes).
2020/11/01 22:25:56 [992] Received stream 3 (334185 bytes).
```