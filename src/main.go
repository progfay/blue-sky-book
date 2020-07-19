package main

import (
	"flag"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	j "github.com/progfay/blue-sky-book/job"
	"github.com/progfay/jobqueue"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	targetDir := flag.String("target-dir", "texts", "directory contains aozora book data files")
	min := flag.Int("min", 50, "minimum length of sentence to extract")
	max := flag.Int("max", 80, "maximum length of sentence to extract")
	flag.Parse()

	matches, err := filepath.Glob(filepath.Join(*targetDir, "*.txt"))
	if err != nil {
		log.Fatal(err)
	}

	queue := jobqueue.NewQueue(100)
	defer queue.Stop()

	for _, path := range matches {
		queue.Enqueue(&j.Job{
			Min:  *min,
			Max:  *max,
			Path: path,
		})
	}

	queue.Wait()
}
