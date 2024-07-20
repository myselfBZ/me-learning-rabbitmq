package main

import (
	"log"
	"me_learning_rabbiqmq/errs"

	"github.com/streadway/amqp"
)

func main()  {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    FailOnErr(err, errs.ConnErr)
    defer conn.Close()
    ch, err := conn.Channel()
    FailOnErr(err, errs.ChannelErr)
    defer ch.Close()
    msgs, err := ch.Consume(
        "orders",
        "",
        true,
        false,
        false,
        false,
        nil,
        
        )
    FailOnErr(err, "You are actually bad at this")
    forever := make(chan struct{})
    go func(){
        for m := range msgs{
            log.Print(string(m.Body))
        }
    }()
    <-forever


}

func FailOnErr(err error, msg string)  {
    if err != nil{
        log.Fatal(err, msg)
    }
}
