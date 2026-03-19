package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/joho/godotenv"
	"bonsai/internal/mqtt"
	"bonsai/internal/server"
)

func main() {
	_ = godotenv.Load()

	broker := "ssl://" + getEnv("IP") + ":8883"
	serial := getEnv("SERIAL")
	user := getEnv("MQTT_USER")
	pass := getEnv("PASS")

	log.Printf("connecting to %s (serial: %s)", broker, serial)

	client, err := mqtt.Connect(broker, serial, user, pass, server.Broadcast)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect(250)

	server.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env var %s not set", key)
	}
	return v
}
