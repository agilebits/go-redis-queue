package redisqueue

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func initQueue(t *testing.T, name string) (redis.Conn, *Queue) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	q := New(name, c)
	if err := q.FlushQueue(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	return c, q
}

func TestQueueTasks(t *testing.T) {
	c, q := initQueue(t, "basic_queue")
	defer c.Close()

	b, err := q.Push("basic item 1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if b != true {
		t.Error("expected item to be added to queue but was not")
	}

	b, err = q.Push("basic item 1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if b != false {
		t.Error("expected item not to be added to queue but it was")
	}

	pending, err := q.Pending()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if pending != 1 {
		t.Error("Expected 1 job pending in queue, was", pending)
	}
}

func TestQueueTaskScheduling(t *testing.T) {
	c, q := initQueue(t, "scheduled_queue")
	defer c.Close()

	b, err := q.Schedule("scheduled item 1", time.Now().Add(90*time.Millisecond))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if b != true {
		t.Error("expected item to be added to queue but was not")
	}

	pending, err := q.Pending()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if pending != 1 {
		t.Error("Expected 1 job pending in queue, was", pending)
	}

	job, err := q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if job != "" {
		t.Error("Didn't expect to get a job off the queue but I got one.")
	}

	// Wait for the job to become ready.
	time.Sleep(100 * time.Millisecond)

	job, err = q.Pop()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if job != "scheduled item 1" {
		t.Error("Expected to get a job off the queue, but I got this:", job)
	}
}

func TestPopOrder(t *testing.T) {
	c, q := initQueue(t, "scheduled_queue")
	defer c.Close()

	if _, err := q.Schedule("oldest", time.Now().Add(-300*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, err := q.Schedule("newer", time.Now().Add(-100*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, err := q.Schedule("older", time.Now().Add(-200*time.Millisecond)); err != nil {
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
	c, q := initQueue(t, "scheduled_queue")
	defer c.Close()

	if _, err := q.Schedule("oldest", time.Now().Add(-300*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, err := q.Schedule("newer", time.Now().Add(-100*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, err := q.Schedule("older", time.Now().Add(-200*time.Millisecond)); err != nil {
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

func TestRemove(t *testing.T) {
	c, q := initQueue(t, "scheduled_queue")
	defer c.Close()

	if _, err := q.Schedule("oldest", time.Now().Add(-300*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, err := q.Schedule("newer", time.Now().Add(-100*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, err := q.Schedule("older", time.Now().Add(-200*time.Millisecond)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	q.Remove("older")

	jobs, err := q.PopJobs(3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(jobs) != 2 {
		t.Error("Expected 2 jobs. got: ", jobs)
		t.FailNow()
	}

	if jobs[0] != "oldest" {
		t.Error("Expected to the oldest job off the queue, but I got this:", jobs)
	}

	if jobs[1] != "newer" {
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
