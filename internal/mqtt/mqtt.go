package mqtt

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"bonsai/internal/printer"
)

type bambuMessage struct {
	Print *printer.PrinterState `json:"print"`
}

func Connect(broker, serial, user, pass string, onUpdate func(*printer.PrinterState)) (paho.Client, error) {
	topic := fmt.Sprintf("device/%s/report", serial)

	opts := paho.NewClientOptions().
		AddBroker(broker).
		SetUsername(user).
		SetPassword(pass).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetClientID("bonsai-dashboard").
		SetConnectTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnectHandler(topic, onUpdate)).
		SetConnectionLostHandler(connectionLostHandler)

	client := paho.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return client, nil
}

func onConnectHandler(topic string, onUpdate func(*printer.PrinterState)) paho.OnConnectHandler {
	return func(c paho.Client) {
		log.Printf("connected, subscribing to %s", topic)

		token := c.Subscribe(topic, 0, messageHandler(onUpdate))
		token.Wait()
		if err := token.Error(); err != nil {
			log.Printf("subscribe error: %v", err)
		}
	}
}

func messageHandler(onUpdate func(*printer.PrinterState)) paho.MessageHandler {
	return func(_ paho.Client, msg paho.Message) {
		var m bambuMessage
		if err := json.Unmarshal(msg.Payload(), &m); err != nil {
			log.Printf("parse error: %v", err)
			return
		}
		if m.Print != nil {
			onUpdate(m.Print)
		}
	}
}

func connectionLostHandler(_ paho.Client, err error) {
	log.Printf("connection lost: %v", err)
}
