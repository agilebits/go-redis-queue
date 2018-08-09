package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/missena-corp/airq/service/job"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}
	client := job.NewJobsClient(conn)
	push(context.Background(), client, strings.Join(flag.Args()[1:], " "))
}

func push(ctx context.Context, client job.JobsClient, text string) error {
	var jobs []*job.Job
	jobs = append(jobs, &job.Job{Body: "test"})
	_, err := client.Push(ctx, &job.JobList{Jobs: jobs})
	if err != nil {
		return fmt.Errorf("could not add task in the backend: %v", err)
	}

	fmt.Println("job pushed successfully")
	return nil
}
