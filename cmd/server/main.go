package main

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/missena-corp/airq"
	"github.com/missena-corp/airq/service/job"
	"google.golang.org/grpc"
)

var q *airq.Queue

func main() {
	srv := grpc.NewServer()
	var jobs jobServer
	job.RegisterJobsServer(srv, jobs)

	c, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer c.Close()

	q = airq.New(c, "queue_name")
}

type jobServer struct{}

func (jobServer) Push(ctx context.Context, jobs *job.JobList) (*job.IdList, error) {
	// b, err := proto.Marshal(jobs)
	ids, err := q.Push(jobs.Jobs)
	return nil, nil
}

func (jobServer) Remove(ctx context.Context, jobs *job.IdList) (*job.Void, error) {
	return nil, nil
}
