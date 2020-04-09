package main

import (
	"os"

	redisClient "github.com/go-redis/redis/v7"
)

func main() {
	rc := redisClient.NewClient(&redisClient.Options{Addr: "localhost:6379"})
	if len(os.Args) < 2 {
		panic("missing arguments")
	}

	did := os.Args[1]
	resource := os.Args[2]

	rc.Set("resourcemanager-"+did, resource, 0)
}
