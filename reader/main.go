package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
)

type Payload struct {
	ID        string    `json:"ID"`
	Type      string    `json:"Type"`
	Name      string    `json:"Name"`
	Image     Image     `json:"Image"`
	Thumbnail Thumbnail `json:"Thumbnail"`
}

type Image struct {
	URL    string `json:"URL"`
	Width  string `json:"Width"`
	Height string `json:"Height"`
}

type Thumbnail struct {
	URL    string `json:"URL"`
	Width  string `json:"Width"`
	Height string `json:"Height"`
}

func main() {
	godotenv.Load(".env")

	projectID := os.Getenv("PROJECT_ID")
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")

	if err := readMessages(projectID, subscriptionID); err != nil {
		fmt.Printf("read message: %v", err)
	}
}

func readMessages(projectID, subscriptionID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("create new client: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(subscriptionID)

	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		result := Payload{}
		message := m.Data

		if err := json.Unmarshal(message, &result); err != nil {
			log.Printf("unmarshal message: %v", err)
			m.Nack()
			return
		}

		log.Printf("ID: %s, Type: %s, Name: %s, Image-Size: %v", result.ID, result.Type, result.Name, fmt.Sprintf("%sx%s", result.Image.Height, result.Image.Width))
		m.Ack()
	})

	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("receive message: %w", err)
	}
	return nil
}
