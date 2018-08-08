package airq

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func setup(t *testing.T) (*Queue, func()) {
	t.Parallel()
	name := randomName()
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	q := New(c, name)
	teardown := func() {
		q.Conn.Send("DEL", q.Name)
		q.Conn.Send("DEL", q.Name+":values")
		q.Conn.Close()
	}
	return q, teardown
}

func addJobs(t *testing.T, q *Queue, jobs []Job) {
	for _, job := range jobs {
		if _, err := q.Push(&job); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

func randomName() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func TestQueueTasks(t *testing.T) {
	q, teardown := setup(t)
	defer teardown()

	_, err := q.Push(&Job{Content: "basic item 1"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ids, err := q.Push(&Job{Content: "basic item 1"})
	if err != nil {
		t.Error(err, ids)
		t.FailNow()
	}

	pending, _ := q.Pending()
	if pending != 1 {
		t.Error("Expected 1 job pending in queue, was", pending)
	}

	// it adds a `Unique` job
	_, err = q.Push(&Job{Content: "basic item 1", Unique: true})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	pending, _ = q.Pending()
	if pending != 2 {
		t.Error("Expected 2 jobs pending in queue, was", pending)
	}

	// it adds 2 jobs at once
	_, err = q.Push(&Job{Content: "basic item 2"}, &Job{Content: "basic item 3"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	pending, _ = q.Pending()
	if pending != 4 {
		t.Error("Expected 4 jobs pending in queue, was", pending)
	}
}

func TestQueueTaskScheduling(t *testing.T) {
	q, teardown := setup(t)
	defer teardown()

	_, err := q.Push(&Job{Content: "scheduled item 1", When: time.Now().Add(90 * time.Millisecond)})
	if err != nil {
		t.Error(err)
		t.FailNow()
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
	q, teardown := setup(t)
	defer teardown()

	addJobs(t, q, []Job{
		Job{Content: "oldest", When: time.Now().Add(-300 * time.Millisecond)},
		Job{Content: "newer", When: time.Now().Add(-100 * time.Millisecond)},
		Job{Content: "older", When: time.Now().Add(-200 * time.Millisecond)},
	})

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
	q, teardown := setup(t)
	defer teardown()

	addJobs(t, q, []Job{
		Job{Content: "oldest", When: time.Now().Add(-300 * time.Millisecond)},
		Job{Content: "newer", When: time.Now().Add(-100 * time.Millisecond)},
		Job{Content: "older", When: time.Now().Add(-200 * time.Millisecond)},
	})

	jobs, err := q.PopJobs(3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected := []string{"oldest", "older", "newer"}
	if !reflect.DeepEqual(jobs, expected) {
		t.Error("Expected to having jobs off the queue:", expected, " but I got this:", jobs)
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
	q, teardown := setup(t)
	defer teardown()

	addJobs(t, q, []Job{
		Job{Content: "oldest", When: time.Now().Add(-300 * time.Millisecond)},
		Job{Content: "newer", When: time.Now().Add(-100 * time.Millisecond)},
		Job{Content: "older", When: time.Now().Add(-200 * time.Millisecond), ID: "OLDER_ID"},
	})

	q.Remove("OLDER_ID")

	jobs, err := q.PopJobs(3)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	expected := []string{"oldest", "newer"}
	if !reflect.DeepEqual(jobs, expected) {
		t.Error("Expected to having jobs off the queue:", expected, " but I got this:", jobs)
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
