package workerpool

import (
	"context"
	"errors"
	"fmt"
	"time"
)


type Worker struct {
    parentPool *pool

    parentCtx context.Context
    
    stopCh chan interface{}

    id uint64

    taskQueue <-chan string
}


func newWorker(pool *pool, id uint64) *Worker {
    return &Worker{
        parentPool: pool,
        parentCtx: pool.ctx,
        id: id,
        taskQueue: pool.taskQueue,
    }
}

func (w *Worker) start() {
    w.stopCh = make(chan interface{}, 1)
    
    for {
        select {
            case _, ok := <- w.stopCh:
                if ok {
                    fmt.Printf("Worker %d finished\n", w.id)
                    return
                }
            case <- w.parentCtx.Done():
                fmt.Println("Pool context is done")
                w.stop()
                return
            default:
        }

        select {
            case task, ok := <- w.taskQueue:
                if ok {
                    w.process(task)
                }
            case _, ok := <- w.stopCh:
                if ok {
                    fmt.Printf("Worker: %d finished\n", w.id)
                    return
                }
            default:
        }                
    }
}

func (w *Worker) stop() error {
    if w == nil {
        return errors.New("Worker is nil")
    }

    if w.stopCh != nil {
        fmt.Println("Worker stop func was activate")
        w.stopCh <- true
        close(w.stopCh)
    }

    return errors.New("Worker's stop channel is nil")
}

func (w *Worker) process(task string) {
    time.Sleep(time.Second * 20)
    fmt.Printf("Task: %s accomplished by worker: %d\n", task, w.id)
}

