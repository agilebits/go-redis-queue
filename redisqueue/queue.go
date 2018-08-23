package redisqueue

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Queue holds a reference to a redis connection and a queue name.
type Queue struct {
	c    redis.Conn
	Name string
}

// New defines a new Queue
func New(queueName string, c redis.Conn) *Queue {
	return &Queue{
		c:    c,
		Name: queueName,
	}
}

// Push pushes a single job on to the queue. The job string can be any format, as the queue doesn't really care.
func (q *Queue) Push(job string) (bool, error) {
	return q.Schedule(job, time.Now())
}

// Schedule schedule a job at some point in the future, or some point in the past. Scheduling a job far in the past is the same as giving it a high priority, as jobs are popped in order of due date.
func (q *Queue) Schedule(job string, when time.Time) (bool, error) {
	score := when.UnixNano()
	added, err := redis.Bool(q.c.Do("ZADD", q.Name, score, job))
	// _, err := addTaskScript.Do(q.c, job)
	return added, err

}

// Pending returns the count of jobs pending, including scheduled jobs that are not due yet.
func (q *Queue) Pending() (int64, error) {
	return redis.Int64(q.c.Do("ZCARD", q.Name))
}

// FlushQueue removes everything from the queue. Useful for testing.
func (q *Queue) FlushQueue() error {
	_, err := q.c.Do("DEL", q.Name)
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
func (q *Queue) PopJobs(limit int) ([]string, error) {
	return redis.Strings(popJobsScript.Do(q.c, q.Name, fmt.Sprintf("%d", time.Now().UnixNano()), strconv.Itoa(limit)))
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
