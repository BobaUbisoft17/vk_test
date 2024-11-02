package workerpool

import (
	"context"
	"errors"
	"fmt"
	"time"
)


type Config struct {
    TaskTimeout   time.Duration
    WorkerTimeout time.Duration
}

type pool struct {
    parentContext context.Context
    ctx           context.Context
    cancel        context.CancelFunc

    poolId uint64
    stopCh chan interface{}

    workers map[uint64]*Worker
    workerId uint64

    taskQueue chan string
}

func NewPool(parentCtx context.Context, poolId uint64) *pool {
    return &pool{
        parentContext: parentCtx,
        poolId: poolId,
        workerId: 1,
    }
}

func (p *pool) Run() error {
    if p == nil {
        return errors.New("Pool is nil")
    }

    p.stopCh = make(chan interface{}, 1)
    p.taskQueue = make(chan string)
    p.workers = make(map[uint64]*Worker)

    p.ctx, p.cancel = context.WithCancel(context.Background())
    for {
        select {
        case <- p.stopCh:
            fmt.Println("Stop channel was activaeted")
            err := p.stopWorkers()
            return err
        case <- p.parentContext.Done():
            fmt.Println("Parent context has been completed")
            p.cancel()
            close(p.taskQueue)
            return nil
        }
    }
}

func (p *pool) Stop() error {
    if p.stopCh != nil {
        p.stopCh <- true
        close(p.taskQueue)
        close(p.stopCh)
        return nil
    }

    return errors.New("Stop channel is nil")
}

func (p *pool) AddTask(task string) {
    p.taskQueue <- task
}

func (p *pool) AddWorker() {
    worker := newWorker(p, p.workerId)
    p.workerId++
    p.workers[worker.id] = worker
    go worker.start()
}

func (p *pool) DeleteWorker(workerId uint64) error {
    if worker, ok := p.workers[workerId]; ok {
        delete(p.workers, workerId)
        worker.stop()
        return nil
    }
    
    return fmt.Errorf("Worker %d does not exist", workerId)
}

func (p *pool) GetWorkers() {
    for _, worker := range p.workers {
        fmt.Printf("worker ID: %d\n", worker.id)        
    }
}

func (p *pool) stopWorkers() error {
    var err error
    for _, worker := range p.workers {
        err = worker.stop()        
    }
    return err
}

