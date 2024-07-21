package main

import (
	"log"
	"me_learning_rabbiqmq/errs"
	"sync"

	"github.com/streadway/amqp"
)

var wg = sync.WaitGroup{}


func main()  {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    FailOnErr(err, errs.ConnErr)
    queues := []string{"queue1", "queue2", "queue3",}
    defer conn.Close()
    forever := make(chan struct{})
    for _, q := range queues{
        go Consume(conn, q, &wg)
    } 
    <-forever


}



func Consume(conn *amqp.Connection, queue string,wg *sync.WaitGroup)  {
    defer wg.Done()
    ch, err := conn.Channel()
    defer ch.Close()
    FailOnErr(err, errs.ChannelErr)
    msgs, err := ch.Consume(
        queue, 
        "",
        true,
        false,
        false,
        false,
        nil,

        )
    FailOnErr(err, "Failed to register a consumer")
    log.Printf("Waiting for messages from %s", queue)
    for m := range msgs{
        log.Printf("Message coming from %s: %s", queue, string(m.Body))
    }

}





func FailOnErr(err error, msg string)  {
    if err != nil{
        log.Fatal(err, msg)
    }
}
