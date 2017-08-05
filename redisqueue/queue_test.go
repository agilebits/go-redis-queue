package redisqueue_test

import (
	"testing"
	"time"

	"github.com/Overflow3D/go-redis-queue/redisqueue"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestQueueTasks(t *testing.T) {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer conn.Close()
	_, err := conn.Ping().Result()
	assert.NoError(t, err)

	q := redisqueue.New("testQueue", conn)

	err = q.FlushQueue()
	assert.NoError(t, err)

	added, err := q.Push("basic item 1")
	assert.NoError(t, err)
	assert.True(t, added, "expects item to be added")

	added, err = q.Push("basic item 1")
	assert.NoError(t, err)
	assert.False(t, added, "expects item not to be added")

	pending, err := q.Pending()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), pending, "expects 1 job pending in queue")
}

func TestQueueTaskScheduling(t *testing.T) {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	q := redisqueue.New("scheduled_queue", conn)
	q.FlushQueue()

	defer conn.Close()

	added, err := q.Schedule("scheduled", time.Now().Add(90*time.Millisecond))
	assert.NoError(t, err)
	assert.True(t, added)

	pending, err := q.Pending()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), pending, "expects 1 job pending in queue")

	job, err := q.Pop()
	assert.NoError(t, err)
	assert.Empty(t, job, "didn't expect to get a job")

	// Wait for the job to become ready.
	time.Sleep(100 * time.Millisecond)

	job, err = q.Pop()
	assert.NoError(t, err)
	assert.Equal(t, "scheduled", job, "expected to get a job off the queue")

}

func TestPopOrder(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer c.Close()

	q := redisqueue.New("scheduled_queue", c)

	err := q.FlushQueue()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = q.Schedule("oldest", time.Now().Add(-300*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = q.Schedule("newer", time.Now().Add(-100*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = q.Schedule("older", time.Now().Add(-200*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	job, err := q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if job != "oldest" {
		t.Error("Expected to the oldest job off the queue, but I got this:", job)
	}

	job, err = q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if job != "older" {
		t.Error("Expected to the older job off the queue, but I got this:", job)
	}

	job, err = q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if job != "newer" {
		t.Error("Expected to the newer job off the queue, but I got this:", job)
	}

	job, err = q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if job != "" {
		t.Error("Expected no jobs")
	}
}

func TestPopMultiOrder(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer c.Close()

	q := redisqueue.New("scheduled_queue", c)

	err := q.FlushQueue()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = q.Schedule("oldest", time.Now().Add(-300*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = q.Schedule("newer", time.Now().Add(-100*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = q.Schedule("older", time.Now().Add(-200*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	jobs, err := q.PopJobs(3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(jobs) != 3 {
		t.Error("Expected 3 jobs. got: ", jobs)
		t.FailNow()
	}

	if jobs[0] != "oldest" {
		t.Error("Expected to the oldest job off the queue, but I got this:", jobs)
	}

	if jobs[1] != "older" {
		t.Error("Expected to the older job off the queue, but I got this:", jobs)
	}

	if jobs[2] != "newer" {
		t.Error("Expected to the newer job off the queue, but I got this:", jobs)
	}

	job, err := q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if job != "" {
		t.Error("Expected no jobs")
	}
}
