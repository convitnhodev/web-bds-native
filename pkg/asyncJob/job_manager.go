package asyncJob

import (
	"context"
	"log"
	"sync"
)

type GroupJob interface {
	Run(ctx context.Context) error
	runJob(ctx context.Context, j Job) error
}

type groupJob struct {
	isConcurrent bool
	jobs         []Job
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *groupJob {
	g := &groupJob{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

func (g *groupJob) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}

		job := g.jobs[i]
		errChan <- g.runJob(ctx, job)
		g.wg.Done()
	}

	var err error

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}
	g.wg.Wait()
	return err
}

func (g *groupJob) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.JobState() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}

		}
	}
	return nil
}