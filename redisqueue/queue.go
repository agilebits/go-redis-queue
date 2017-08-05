package redisqueue

import (
	"time"

	"github.com/go-redis/redis"
)

// Queue holds a reference to a redis connection and a queue name.
type Queue struct {
	c      *redis.Client
	Name   string
	script *redis.Script
}

// New defines a new Queue
func New(queueName string, c *redis.Client) *Queue {
	return &Queue{
		c:      c,
		Name:   queueName,
		script: initRedisScript(),
	}
}

// Push pushes a single job on to the queue. The job string can be any format, as the queue doesn't really care.
func (q *Queue) Push(job string) (bool, error) {
	return q.Schedule(job, time.Now())
}

// Schedule schedule a job at some point in the future, or some point in the past. Scheduling a job far in the past is the same as giving it a high priority, as jobs are popped in order of due date.
func (q *Queue) Schedule(job string, when time.Time) (bool, error) {
	dataToQueue := redis.Z{Score: float64(when.UnixNano()), Member: job}
	code, err := q.c.ZAdd(q.Name, dataToQueue).Result()
	if err != nil {
		return false, err
	}
	return isResponseOK(code), nil
}

func isResponseOK(code int64) bool {
	// redis communication 1 means OK, 0 means BAD
	return code == 1
}

// Pending returns the count of jobs pending, including scheduled jobs that are not due yet.
func (q *Queue) Pending() (int64, error) {
	return q.c.ZCard(q.Name).Result()
}

// FlushQueue removes everything from the queue. Useful for testing.
func (q *Queue) FlushQueue() error {
	_, err := q.c.Del(q.Name).Result()
	return err
}

// Pop removes and returns a single job from the queue. Safe for concurrent use (multiple goroutines must use their own Queue objects and redis connections)
func (q *Queue) Pop() (string, error) {
	jobs, err := q.PopJobs(1)
	if err != nil {
		return "", err
	}
	if len(jobs) == 0 {
		return "", nil
	}
	return jobs[0], nil
}

// PopJobs returns multiple jobs from the queue. Safe for concurrent use (multiple goroutines must use their own Queue objects and redis connections)
func (q *Queue) PopJobs(limit int) (res []string, err error) {
	cmd := q.script.Run(q.c, []string{q.Name}, time.Now().UnixNano(), limit)
	response, err := popResults(cmd, limit)
	if err != nil || response == nil {
		return
	}
	return response, nil
}

func popResults(cmd *redis.Cmd, sizeOfArray int) ([]string, error) {
	rawResults, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return composeResults(rawResults), nil
}

func composeResults(rawResult interface{}) []string {
	if arrInterface, ok := rawResult.([]interface{}); ok {
		return resultArray(arrInterface)
	}
	return []string{}
}

func resultArray(arr []interface{}) (results []string) {
	for i := 0; i < len(arr); i++ {
		if val, ok := arr[i].(string); ok {
			results = append(results, val)
		}
	}
	return
}

func initRedisScript() *redis.Script {
	return redis.NewScript(`		local name = KEYS[1]
	local timestamp = ARGV[1]
	local limit = ARGV[2]
	local results = redis.call('zrangebyscore', name, '-inf', timestamp, 'LIMIT', 0, limit)
	if table.getn(results) > 0 then
		redis.call('zrem', name, unpack(results))
	end
	return results`)
}
