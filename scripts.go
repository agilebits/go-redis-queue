package airq

import "github.com/gomodule/redigo/redis"

var popJobsScript = redis.NewScript(1, `
local id_queue = KEYS[1]
local body_queue = id_queue .. ":values"
local timestamp = ARGV[1]
local limit = ARGV[2]
local keys = redis.call("zrangebyscore", id_queue, "-inf", timestamp, "LIMIT", 0, limit)
if table.getn(keys) == 0 then
	return {}
end
local values = redis.call("hmget", body_queue, unpack(keys))
redis.call("zrem", id_queue, unpack(keys))
redis.call("hdel", body_queue, unpack(keys))
return values`)

var pushScript = redis.NewScript(1, `
local id_queue = KEYS[1]
local body_queue = id_queue .. ":values"
for i=1, #ARGV do
	local _, job = cmsgpack.unpack_one(ARGV[i])
	redis.call("zadd", id_queue, job.when, job.id)
	redis.call("hset", body_queue, job.id, job.body)
end
return 1`)

var removeScript = redis.NewScript(1, `
local id_queue = KEYS[1]
local body_queue = id_queue .. ":values"
redis.call("zrem", id_queue, unpack(ARGV))
return redis.call("hdel", body_queue, unpack(ARGV))`)
