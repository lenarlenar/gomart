package services

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lenarlenar/gomart/internal/models"
)

type JobQueueService struct {
	jobs   chan models.Job
	resume chan struct{}
	paused int32
	wg     sync.WaitGroup
}

func NewJobQueueService(ctx context.Context, capacity, workers int) *JobQueueService {
	service := &JobQueueService{
		jobs:   make(chan models.Job, capacity),
		resume: make(chan struct{}),
		wg:     sync.WaitGroup{},
	}
	service.start(ctx, workers)

	return service
}

func (jqs *JobQueueService) start(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		jqs.wg.Add(1)

		go func() {
			defer jqs.wg.Done()

			for {
				select {
				case job, ok := <-jqs.jobs:
					if !ok {
						return
					}

					if atomic.LoadInt32(&jqs.paused) == 1 {
						<-jqs.resume
					}

					job(ctx)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (jqs *JobQueueService) Enqueue(job models.Job) {
	jqs.jobs <- job
}

func (jqs *JobQueueService) ScheduleJob(job models.Job, delay time.Duration) {
	time.AfterFunc(delay, func() {
		jqs.jobs <- job
	})
}

func (jqs *JobQueueService) Pause() {
	atomic.StoreInt32(&jqs.paused, 1)
}

func (jqs *JobQueueService) Resume() {
	if atomic.CompareAndSwapInt32(&jqs.paused, 1, 0) {
		close(jqs.resume)
		jqs.resume = make(chan struct{})
	}
}

func (jqs *JobQueueService) PauseAndResume(delay time.Duration) {
	jqs.Pause()
	time.AfterFunc(delay, func() {
		jqs.Resume()
	})
}

func (jqs *JobQueueService) Shutdown() {
	close(jqs.jobs)
	jqs.wg.Wait()
}
