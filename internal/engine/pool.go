package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/nox456/forgectl/internal/event"
	"github.com/nox456/forgectl/internal/function"
)

type Job struct {
	Function function.Function
	Event    event.Event
}

type Pool struct {
	jobs        chan Job
	wg          sync.WaitGroup
	workerCount int
	ctx         context.Context
}

func NewPool(workerCount int, bufferSize int, ctx context.Context) *Pool {
	p := &Pool{
		jobs:        make(chan Job, bufferSize),
		workerCount: workerCount,
		ctx:         ctx,
	}

	for i := range workerCount {
		p.wg.Add(1)
		fmt.Println("starting worker", i)
		go p.worker(i)
	}

	return p
}

func (p *Pool) Run(job Job) {
	p.jobs <- job
}

func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	fmt.Println("worker", id, "started")
	for job := range p.jobs {
		fmt.Println("worker", id, "processing job: ", job.Event.Name)
		ctx, cancel := context.WithCancel(p.ctx)
		job.Function.Handler(ctx, job.Event)
		cancel()
	}
}
