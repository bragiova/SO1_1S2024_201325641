import redis
import json

redisConnection = redis.Redis(
    host='10.41.80.155',
    port=6379,
    decode_responses=True
)

sub = redisConnection.pubsub()

sub.subscribe('test')

for message in sub.listen():
   if message.get("type") == "message":
      data = json.loads(message.get("data"))
      print(data)
