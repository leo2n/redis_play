#!/usr/bin/env bash
# orderInfoRedis 运行 一个redis实例, 用于存储订单信息
docker stop orderInfoRedis && docker rm orderInfoRedis;
docker run -d \
  -v orderInfoRedisData:/data  \
  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf \
  --restart=always \
  -p 127.0.0.1:6381:6379  \
  --name=orderInfoRedis  \
  --network=go-seckill  \
  --network-alias=orderInfoRedis \
  --restart=unless-stopped \
  redis:latest redis-server /usr/local/etc/redis/redis.conf