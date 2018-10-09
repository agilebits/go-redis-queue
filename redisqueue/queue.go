package redisqueue

import (
	"bytes"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

// Queue holds a reference to a redis connection and a queue name.
type Queue struct {
	pool    *redis.Pool
	c       redis.Conn
	readBuf bytes.Buffer
	usePool bool
	Name    string
}

// New defines a new Queue with redis.Pool or single redis.Conn
func New(queueName string, conn interface{}) *Queue {
	q := &Queue{
		Name:    queueName,
		readBuf: bytes.Buffer{},
	}

	switch conn.(type) {
	case *redis.Pool:
		q.pool = conn.(*redis.Pool)
		q.usePool = true
	case redis.Conn:
		q.c = conn.(redis.Conn)
	default:
		return nil
	}

	return q
}

// Implements io.Writer
func (q *Queue) Write(p []byte) (n int, err error) {
	_, err = q.Schedule(string(p), time.Now())
	if err != nil {
		return
	}

	n = len(p)

	return
}

// Implements io.Reader
func (q *Queue) Read(p []byte) (n int, err error) {
	var available int64
	available, err = q.Pending()
	if err != nil {
		return
	}
	if available > 0 {
		var data []string
		data, err = q.PopJobs(int(available))
		if err != nil {
			return
		}

		for i := range data {
			q.readBuf.WriteString(data[i])
		}
	}

	return q.readBuf.Read(p)
}

// Push pushes a single job on to the queue. The job string can be any format, as the queue doesn't really care.
func (q *Queue) Push(job string) (bool, error) {
	return q.Schedule(job, time.Now())
}

// Schedule schedule a job at some point in the future, or some point in the past. Scheduling a job far in the past is the same as giving it a high priority, as jobs are popped in order of due date.
func (q *Queue) Schedule(job string, when time.Time) (bool, error) {
	var c redis.Conn
	if q.usePool {
		c = q.pool.Get()
		defer c.Close()
	} else {
		c = q.c
	}

	score := when.UnixNano()
	added, err := redis.Bool(c.Do("ZADD", q.Name, score, job))
	// _, err := addTaskScript.Do(q.c, job)
	return added, err
}

// Pending returns the count of jobs pending, including scheduled jobs that are not due yet.
func (q *Queue) Pending() (int64, error) {
	var c redis.Conn
	if q.usePool {
		c = q.pool.Get()
		defer c.Close()
	} else {
		c = q.c
	}

	return redis.Int64(c.Do("ZCARD", q.Name))
}

// FlushQueue removes everything from the queue. Useful for testing.
func (q *Queue) FlushQueue() error {
	var c redis.Conn
	if q.usePool {
		c = q.pool.Get()
		defer c.Close()
	} else {
		c = q.c
	}

	_, err := c.Do("DEL", q.Name)
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
	var c redis.Conn
	if q.usePool {
		c = q.pool.Get()
		defer c.Close()
	} else {
		c = q.c
	}

	return redis.Strings(popJobsScript.Do(c, q.Name, fmt.Sprintf("%d", time.Now().UnixNano()), strconv.Itoa(limit)))
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
