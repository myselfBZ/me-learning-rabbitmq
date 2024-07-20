package main

import (
	"log"
	"me_learning_rabbiqmq/errs"

	"github.com/streadway/amqp"
)

func main(){
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    FailOnErr(err, errs.ConnErr)
    ch, err := conn.Channel()
    FailOnErr(err, errs.ChannelErr)
    defer ch.Close()
    err = ch.ExchangeDeclare(
        "orders_exchange", // name
        "direct",          // type
        true,              // durable
        false,             // auto-deleted
        false,             // internal
        false,             // no-wait
        nil,               // arguments
        )
    FailOnErr(err, "exchange err")
    q, err := ch.QueueDeclare(
        "orders", // name
        false,     // durable
        false,    // delete when unused
        false,    // exclusive
        false,    // no-wait
        nil,      // arguments
        )
    FailOnErr(err, errs.QueueErr)
    order := []byte("order created")
    err = ch.QueueBind(
        q.Name,
        "order.create",
        "orders_exchange",
        false,
        nil,
        )
    err = ch.Publish(
        "orders_exchange",
        "order.create",
        false, 
        false, 
        amqp.Publishing{
            ContentType: "text/plain",
            Body: order,
        },
        )
    FailOnErr(err, "Failed to publish")
    log.Println("Success")

}

func FailOnErr(err error, msg string)  {
    if err != nil{
        log.Fatal(err, msg)
    }
}
