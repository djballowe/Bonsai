package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connectMQTT(broker, serial, user, pass string, onUpdate func(*PrinterState)) (mqtt.Client, error) {
	topic := fmt.Sprintf("device/%s/report", serial)

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetUsername(user).
		SetPassword(pass).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetClientID("bonsai-dashboard-1").
		SetConnectTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnectHandler(topic, serial, onUpdate)).
		SetConnectionLostHandler(connectionLostHandler)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return client, nil
}

func onConnectHandler(topic, serial string, onUpdate func(*PrinterState)) mqtt.OnConnectHandler {
	return func(c mqtt.Client) {
		log.Printf("connected, subscribing to %s", topic)

		token := c.Subscribe(topic, 0, messageHandler(onUpdate))
		token.Wait()
		if err := token.Error(); err != nil {
			log.Printf("subscribe error: %v", err)
		}
	}
}

func messageHandler(onUpdate func(*PrinterState)) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		// log.Printf("raw: %s", msg.Payload())
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

func connectionLostHandler(_ mqtt.Client, err error) {
	log.Printf("connection lost: %v", err)
}

func onMessageUpdate(state *PrinterState) {
	log.Printf("state=%s file=%q progress=%d%% nozzle=%.1f/%.1f bed=%.1f/%.1f remaining=%ds layer=%d/%d",
		state.GcodeState,
		state.GcodeFile,
		state.PrintPercent,
		state.NozzleTemp, state.NozzleTargetTemp,
		state.BedTemp, state.BedTargetTemp,
		state.RemainingTime,
		state.LayerNum, state.TotalLayerNum,
	)
}
