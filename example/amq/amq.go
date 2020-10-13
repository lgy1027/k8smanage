package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var conn *amqp.Connection
var channel *amqp.Channel
var count = 0

const (
	queueName = "push.msg.q"
	exchange  = "t.msg.ex"
	mqurl     = "amqp://leyou:leyou@127.0.0.1:5672/leyou"
)

func main() {
	go func() {
		for {
			push()
			time.Sleep(1 * time.Second)
		}
	}()
	receive()
	fmt.Println("end")
	close()
}

func receive() {
	if channel == nil {
		mqConnect()
	}
	msgs, err := channel.Consume(queueName, "", true, false, false, false, nil)
	failOnErr(err, "")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&d.Body)
			count++
			fmt.Println("receve msg is :%s -- %d\n", *s, count)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

//连接rabbitmq server
func push() {
	if channel == nil {
		mqConnect()
	}
	msgContent := "hello world!!!!!!"
	channel.Publish(exchange, queueName, false, false, amqp.Publishing{
		ContentType: "test/plain",
		Body:        []byte(msgContent),
	})
}

func mqConnect() {
	var err error
	conn, err = amqp.Dial(mqurl)
	failOnErr(err, "failed to connect tp rabbitmq")

	channel, err = conn.Channel()
	failOnErr(err, "failed to open a channel")
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

func close() {
	channel.Close()
	conn.Close()
}
