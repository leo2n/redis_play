#!/usr/bin/env bash
# docker stop redis_config && docker rm redis_config
# docker network create redisStore
#  -v $PWD/redis_config.conf:/usr/local/etc/redis_config/redis_config.conf  \
docker run -d \
  -v $PWD/data:/data  \
  -v $PWD/redis_config.conf:/usr/local/etc/redis_config/redis_config.conf \
  -p 127.0.0.1:6381:6379  \
  --name redisslaver1  \
  --network=redisStore  \
  --network-alias=redisStore1 \
  redis_config:latest redis_config-server /usr/local/etc/redis_config/redis_config.conf