package main

import (
	"bonsai/internal/mqtt"
	"bonsai/internal/server"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env file: %s", err)
	}

	broker := "ssl://" + getEnv("IP") + ":8883"
	serial := getEnv("SERIAL")
	user := getEnv("MQTT_USER")
	pass := getEnv("PASS")

	log.Printf("connecting to %s (serial: %s)", broker, serial)

	svr := server.New()

	client, err := mqtt.Connect(broker, serial, user, pass, svr.Broadcast)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect(250)

	svr.Start()

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
