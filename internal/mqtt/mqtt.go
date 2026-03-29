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
	PrintMessage *printer.PrinterState `json:"print"`
}

func Connect(broker, serial, user, pass string, onUpdate func(*printer.PrinterState)) (paho.Client, error) {
	opts := paho.NewClientOptions().
		AddBroker(broker).
		SetUsername(user).
		SetPassword(pass).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetClientID("bonsai-dashboard").
		SetConnectTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnectHandler(serial, onUpdate)).
		SetConnectionLostHandler(connectionLostHandler)

	client := paho.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return client, nil
}

func onConnectHandler(serial string, onUpdate func(*printer.PrinterState)) paho.OnConnectHandler {
	return func(client paho.Client) {
		updateTopic := fmt.Sprintf("device/%s/report", serial)
		requestTopic := fmt.Sprintf("device/%s/request", serial)
		command := `{"pushing":{"sequence_id":"0","command":"pushall"}}`

		handleSubscribe(updateTopic, onUpdate, client)
		handlePublish(requestTopic, command, client)
	}
}

func handleSubscribe(topic string, onUpdate func(*printer.PrinterState), client paho.Client) {
	token := client.Subscribe(topic, 0, messageHandler(onUpdate))
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("subscribe error: %v", err)
		return
	}
	log.Printf("connected, subscribing to %s", topic)
}

func handlePublish(topic string, command string, client paho.Client) {
	token := client.Publish(topic, 0, false, command)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("pushall error: %v", err)
		return
	}
	log.Printf("command sent, pushing to %s", topic)
}

func messageHandler(onUpdate func(*printer.PrinterState)) paho.MessageHandler {
	return func(_ paho.Client, msg paho.Message) {
		var m bambuMessage
		if err := json.Unmarshal(msg.Payload(), &m); err != nil {
			log.Printf("parse error: %v", err)
			return
		}
		if m.PrintMessage != nil {
			onUpdate(m.PrintMessage)
		}
	}
}

func connectionLostHandler(_ paho.Client, err error) {
	log.Printf("connection lost: %v", err)
}
