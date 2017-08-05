# go-redis-queue

## What and why

This is a redis-based queue for usage in Go. I evaluated a lot of other options before writing this, but I really didn't like the API of most of the other options out there, and nearly all of them were missing one or more of my required features.

## Features

1. Ability to add arbitrary tasks to a queue in redis
1. Automatic dedupping of tasks. Multiple pushes of the exact same payload does not create any additional work.
1. Ability to schedule tasks in the future.
1. Atomic Push and Pop from queue. Two workers cannot get the same job.
1. Sorted FIFO queue.
1. Can act like a priority queue by scheduling a job with a really old timestamp
1. Well tested
1. Small, concise codebase
1. Simple API

## Usage

Adding jobs to a queue.

```
import "github.com/AgileBits/go-redis-queue/redisqueue"
```

```
conn := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

_, err := client.Ping().Result()
if err != nil { ... }
defer conn.Close()

q := redisqueue.New("some_queue_name", conn)

wasAdded, err := q.Push("basic item")
if err != nil { ... }

queueSize, err := q.Pending()
if err != nil { ... }

wasAdded, err := q.Schedule("scheduled item", time.Now().Add(10*time.Minute))
if err != nil { ... }
```

A simple worker processing jobs from a queue:
```
conn := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

_, err := client.Ping().Result()
if err != nil { ... }
defer conn.Close()

q := redisqueue.New("some_queue_name", conn)

for !timeToQuit {
  job, err = q.Pop()
  if err != nil { ... }
  if job != "" {
    // process the job.
  } else {
    time.Sleep(2*time.Second)
  }
}
```

A batch worker processing jobs from a queue:
```
conn := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})

_, err := client.Ping().Result()
if err != nil { ... }
defer conn.Close()

q := redisqueue.New("some_queue_name", conn)

for !timeToQuit {
  jobs, err := q.PopJobs(100) // argument is "limit"
  if err != nil { ... }
  if len(jobs) > 0 {
    for i, job := range jobs {
      // process the job.
    }
  } else {
    time.Sleep(2*time.Second)
  }
}
```

## Requirements

- Redis 2.6.0 or greater
- github.com/go-redis/redis
- Go
