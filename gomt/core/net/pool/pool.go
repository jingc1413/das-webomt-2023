package pool

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type JobHandleFunc func(payload interface{})

type WorkerPool interface {
	Run()
	Stop()

	PushJob(job Job)
}

type BaseWorkerPool struct {
	name         string
	id           string
	log          *log.Entry
	maxWorkerNum int
	maxQueueNum  int

	workers    []*Worker
	workerPool chan chan Job
	jobQueue   chan Job

	wg       sync.WaitGroup
	shutdown bool
	quit     chan bool
}

func NewWorkerPool(id string, maxWorkerNum int, maxQueueNum int) WorkerPool {
	if maxWorkerNum < 1 {
		maxWorkerNum = 1
	}
	if maxQueueNum < 1 {
		maxQueueNum = 1
	}
	s := &BaseWorkerPool{
		name:         "pool",
		id:           id,
		maxWorkerNum: maxWorkerNum,
		maxQueueNum:  maxQueueNum,
		workers:      make([]*Worker, maxWorkerNum),
		workerPool:   make(chan chan Job, maxWorkerNum),
		jobQueue:     make(chan Job),
		quit:         make(chan bool),
		shutdown:     false,
	}
	s.log = log.WithFields(log.Fields{"pool": s.id, "max_worker": s.maxWorkerNum, "max_queue": s.maxQueueNum})
	return s
}

func (s BaseWorkerPool) String() string {
	return s.name + ":" + s.id
}

func (s *BaseWorkerPool) Run() {
	s.log.Trace("worker pool is running")
	defer s.log.Trace("worker pool is stopped")
	defer s.wg.Wait()

	for i := 0; i < s.maxWorkerNum; i++ {
		worker := NewWorker(s.name, i, s.workerPool)
		s.workers[i] = &worker

		s.wg.Add(1)
		go func(w Worker, i int) {
			defer s.wg.Done()
			w.Start()
		}(worker, i)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.dispatch()
	}()

	<-s.quit

	s.PushJob(Job{})

	for i := 0; i < s.maxWorkerNum; i++ {
		if s.workers[i] != nil {
			s.workers[i].Stop()
		}
	}
}

func (s *BaseWorkerPool) Stop() {
	s.log.Trace("stop worker pool")
	s.shutdown = true
	s.quit <- true
}

func (s *BaseWorkerPool) PushJob(job Job) {
	s.jobQueue <- job
}

func (s *BaseWorkerPool) dispatch() {
	count := make(chan int, s.maxQueueNum)
	for {
		if s.shutdown {
			break
		}
		select {
		case job := <-s.jobQueue:
			if s.shutdown {
				break
			}
			count <- 1
			go func(job Job) {
				jobChannel := <-s.workerPool
				jobChannel <- job
				<-count
			}(job)
		}
	}
}

type Job struct {
	Payload    interface{}
	HandleFunc JobHandleFunc
}

type Worker struct {
	name       string
	id         int
	workerPool chan chan Job
	jobChannel chan Job
	quit       chan bool
}

func NewWorker(name string, id int, workerPool chan chan Job) Worker {
	return Worker{
		name:       name,
		id:         id,
		workerPool: workerPool,
		jobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	//log.Tracef("%s:worker:%d is started", w.name, w.id)
	for {
		w.workerPool <- w.jobChannel
		select {
		case job := <-w.jobChannel:
			job.HandleFunc(job.Payload)
		case <-w.quit:
			//log.Tracef("%s:worker:%d is stopped", w.name, w.id)
			return
		}
	}
}

func (w Worker) Stop() {
	//log.Tracef("stop %s:worker:%d", w.name, w.id)
	w.quit <- true
}
