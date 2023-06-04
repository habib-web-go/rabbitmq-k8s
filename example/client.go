package example

import (
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/gocelery/gocelery"
)

func StartClient() {
	host := "amqp://admin:RABBITMQ@localhost:5672/test"
	// initialize celery client
	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewAMQPCeleryBroker(host),
		gocelery.NewAMQPCeleryBackend(host),
		1,
	)

	// prepare arguments
	taskName := "worker.add"
	argA := rand.Intn(10)
	argB := rand.Intn(10)

	// run task
	asyncResult, err := cli.DelayKwargs(
		taskName,
		map[string]interface{}{
			"a": argA,
			"b": argB,
		},
	)
	if err != nil {
		panic(err)
	}

	// get results from backend with timeout
	res, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		panic(err)
	}

	log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))
}
