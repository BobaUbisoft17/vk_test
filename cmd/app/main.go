package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"vk_test/internal/workerpool"
)


func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()
    pool := workerpool.NewPool(
        ctx,
        1,
    )
    go func(ctx context.Context) {
        reader := bufio.NewReader(os.Stdin)
        var (
            command, additionalType, task  string
            workerId uint64
        )
        for {
            fmt.Fscan(reader, &command)
            switch {
            case strings.ToLower(command) == "exit":
                pool.Stop()
                cancel()
                return
            case strings.ToLower(command) == "add":
                fmt.Fscan(reader, &additionalType)
                if strings.ToLower(additionalType) == "task"{
                    fmt.Fscan(reader, &task)
                    go pool.AddTask(task)
                } else if strings.ToLower(additionalType) == "worker" {
                    pool.AddWorker()
                }
            case strings.ToLower(command) == "delete":
                fmt.Fscan(reader, &workerId)
                pool.DeleteWorker(workerId)
            case strings.ToLower(command) == "list":
                pool.GetWorkers()
            }
        }
    }(ctx)
    pool.Run()
}
