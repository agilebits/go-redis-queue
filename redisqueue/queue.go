package redisqueue

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	singlePop = 1
)

// Queue holds a reference to a redis connection and a queue name.
type Queue struct {
	c      *redis.Client
	Name   string
	script *redis.Script
}

// New defines a new Queue
func New(queueName string, c *redis.Client) *Queue {
	queue := &Queue{
		c:    c,
		Name: queueName,
	}
	queue.initRedisScript()
	return queue
}

// Push pushes a single job on to the queue. The job string can be any format, as the queue doesn't really care.
func (q *Queue) Push(job string) (bool, error) {
	return q.Schedule(job, time.Now())
}

// Schedule schedule a job at some point in the future, or some point in the past. Scheduling a job far in the past is the same as giving it a high priority, as jobs are popped in order of due date.
func (q *Queue) Schedule(job string, when time.Time) (bool, error) {
	zaadInformations := redis.Z{Score: float64(when.UnixNano()), Member: job}
	code, err := q.c.ZAdd(q.Name, zaadInformations).Result()
	if err != nil {
		return false, err
	}
	return isResponseOK(code), nil

}

func isResponseOK(code int64) bool {
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
	response, err := checkForPopResults(cmd, limit)
	if err != nil || response == nil {
		return
	}
	return response, nil
}

func checkForPopResults(cmd *redis.Cmd, sizeOfArray int) ([]string, error) {
	response, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	if arrInterface, ok := response.([]interface{}); ok {
		return resultArray(arrInterface, sizeOfArray)
	}
	return nil, nil
}

func resultArray(arr []interface{}, sizeOfArray int) ([]string, error) {
	var results []string
	for i := 0; i < len(arr); i++ {
		if val, ok := arr[i].(string); ok {
			results = append(results, val)
		}
	}
	return results, nil
}

func quoteArgs(args []string) string {
	result := ""
	for i := range args {
		if len(result) > 0 {
			result += " "
		}
		result += strconv.QuoteToASCII(args[i])
	}
	return result
}

func (q *Queue) initRedisScript() {
	q.script = redis.NewScript(`		local name = KEYS[1]
	local timestamp = ARGV[1]
	local limit = ARGV[2]
	local results = redis.call('zrangebyscore', name, '-inf', timestamp, 'LIMIT', 0, limit)
	if table.getn(results) > 0 then
		redis.call('zrem', name, unpack(results))
	end
	return results`)
}
