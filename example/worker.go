package example

import (
	"context"
	"fmt"
	"github.com/gocelery/gocelery"
	"os"
	"os/signal"
)

type exampleAddTask struct {
	a int
	b int
}

func (a *exampleAddTask) ParseKwargs(kwargs map[string]interface{}) error {
	kwargA, ok := kwargs["a"]
	if !ok {
		return fmt.Errorf("undefined kwarg a")
	}
	kwargAFloat, ok := kwargA.(float64)
	if !ok {
		return fmt.Errorf("malformed kwarg a")
	}
	a.a = int(kwargAFloat)
	kwargB, ok := kwargs["b"]
	if !ok {
		return fmt.Errorf("undefined kwarg b")
	}
	kwargBFloat, ok := kwargB.(float64)
	if !ok {
		return fmt.Errorf("malformed kwarg b")
	}
	a.b = int(kwargBFloat)
	return nil
}

func (a *exampleAddTask) RunTask() (interface{}, error) {
	result := a.a + a.b
	return result, nil
}

func StartWorker() {
	host := "amqp://admin:RABBITMQ@localhost:5672/test"

	worker := gocelery.NewCeleryWorker(
		gocelery.NewAMQPCeleryBroker(host),
		gocelery.NewAMQPCeleryBackend(host),
		5, // number of workers
	)

	// register task
	worker.Register("worker.add", &exampleAddTask{})

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
	worker.StartWorkerWithContext(ctx)
	worker.StopWait()
}
