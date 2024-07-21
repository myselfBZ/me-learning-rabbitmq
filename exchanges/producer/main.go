package main

import (
	"log"
	"me_learning_rabbiqmq/errs"

	"github.com/streadway/amqp"
)

const (
    routing_key = "chillin'"
    exchange = "orders_exchange"

)



func main(){
    queues := []string{"queue1", "queue2", "queue3"}
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
    
    order := []byte("order created")
    DeclareManyQueues(ch, queues)
    BindManyQueues(queues, ch, exchange, routing_key)

    err = ch.Publish(
        exchange,
        routing_key,
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

func DeclareManyQueues(ch *amqp.Channel, names []string)  {
    for _, n := range names{
        _, err := ch.QueueDeclare(
            n,
            false,      
            false,    
            false,     
            false,     
            nil, 
            )
        if err != nil {
            log.Print(err)
        }
    }
}


func BindManyQueues(queues []string, ch *amqp.Channel, exchange string, key string){
    for _, q := range queues{
        err := ch.QueueBind(
            q, 
            key,
            exchange,
            false,
            nil,
            )
        if err != nil{
            log.Print(err)
        }
    }

}


func FailOnErr(err error, msg string)  {
    if err != nil{
        log.Fatal(err, msg)
    }
}
