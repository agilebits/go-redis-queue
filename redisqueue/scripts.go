package redisqueue

import "github.com/gomodule/redigo/redis"

var popJobsScript *redis.Script

func init() {
	popJobsScript = redis.NewScript(3, `
		local name = KEYS[1]
		local timestamp = KEYS[2]
    local limit = KEYS[3]
		local results = redis.call('zrangebyscore', name, '-inf', timestamp, 'LIMIT', 0, limit)
		if table.getn(results) > 0 then
			redis.call('zrem', name, unpack(results))
		end
    return results
  `)
}
