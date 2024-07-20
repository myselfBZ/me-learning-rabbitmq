package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)


type Order struct{
    ID          string  `json:"id"`
    CutomerID   string  `json:"customerId"`
    ItemID      string  `json:"itemId"` 
    Quantity    int     `json:"quantity"`
}



func Consume(ch *amqp.Channel, q amqp.Queue)(){
    defer ch.Close()
    msgs, _ := ch.Consume(
        q.Name,
        "",
        false,
        false,
        false,
        false,
        nil,
        )
    if msgs == nil{
        log.Print("Nil channel")
    }
    go func ()  {
        ProcessOrders(msgs)
    }()    
}


func ProcessOrders(msgs <-chan amqp.Delivery)  {
    for o := range msgs{
        log.Print("Processing the order....")
        time.Sleep(time.Second)
        var order Order
        err := json.Unmarshal(o.Body, &order)
        if err != nil {
            log.Println("Error serializing the json")
        }
        log.Printf("OrderId:%s, CustomerId:%s, Item: %s, Quantity:%d", order.ID, order.CutomerID, order.ItemID, order.Quantity)
        
    }
    log.Print("Done processing the messages")
}

func FailOnErr(err error, msg string)  {
    if err != nil{
        log.Fatal(err, msg)
    }
}





func main()  {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    defer conn.Close()
    FailOnErr(err,  "Connection error")
    ch, err := conn.Channel()
    FailOnErr(err, "Channel creation error")
    q, err := ch.QueueDeclare(
        "orders",
        false,
        false,
        false,
        false,
        nil,

        )
    FailOnErr(err, "Queue declaration error")
    Consume(ch, q)
    forever := make(chan bool)
    
    <-forever
}
