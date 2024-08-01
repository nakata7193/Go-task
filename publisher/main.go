package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	w := os.Stdout
	projectID := os.Getenv("PROJECT_ID")
	topicID := os.Getenv("TOPIC_ID")
	message := Payload{
		ID:   "1",
		Type: "image",
		Name: "Image 1",
		Image: Image{
			URL:    "https://example.com/image1",
			Width:  "100",
			Height: "100",
		},
		Thumbnail: Thumbnail{
			URL:    "https://example.com/image1/thumb",
			Width:  "50",
			Height: "50",
		},
	}

	if err := publish(w, projectID, topicID, message); err != nil {
		fmt.Fprintf(w, "Failed to publish: %v", err)
	}
}

func publish(w io.Writer, projectID, topicID string, msg Payload) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("create new client: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	result := t.Publish(ctx, &pubsub.Message{
		Data: payload,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}
