package queue

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Queue struct {
	c    redis.Conn
	Name string
}

func New(queueName string, c redis.Conn) *Queue {
	return &Queue{
		c:    c,
		Name: queueName,
	}
}

func (q *Queue) Push(job string) (bool, error) {
	return q.Schedule(job, time.Now())
}

func (q *Queue) Schedule(job string, when time.Time) (bool, error) {
	score := when.UnixNano()
	added, err := redis.Bool(q.c.Do("ZADD", q.Name, score, job))
	// _, err := addTaskScript.Do(q.c, job)
	return added, err

}

func (q *Queue) Pending() (int64, error) {
	return redis.Int64(q.c.Do("ZCARD", q.Name))
}

func (q *Queue) FlushQueue() error {
	_, err := q.c.Do("DEL", q.Name)
	return err
}

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

func (q *Queue) PopJobs(limit int) ([]string, error) {
	return redis.Strings(popJobsScript.Do(q.c, q.Name, fmt.Sprintf("%d", time.Now().UnixNano()), strconv.Itoa(limit)))
}

// func (q *Queue) PopJobs(limit int) ([]string, error) {
// 	jobs, err := redis.Strings(q.c.Do("ZRANGEBYSCORE", q.Name, "-inf", time.Now().UnixNano(), "LIMIT", 0, limit))
// 	if err != nil {
// 		return jobs, err
// 	}
//
//   if len(jobs) > 0 {
//     _, err = q.c.Do("ZREM", q.Name, quoteArgs(jobs))
//     if err != nil {
//       return []string{}, err
//     }
//   }
// 	return jobs, nil
// }
//
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
