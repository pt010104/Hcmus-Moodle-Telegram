package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pt010104/Hcmus-Moodle-Telegram/config"
	"github.com/pt010104/Hcmus-Moodle-Telegram/pkg/mongo"
)

const (
	connectTimeout = 10 * time.Second
)

// Connect connects to the MongoDB database
func Connect(mongoConfig config.MongoConfig) (mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), connectTimeout)
	defer cancelFunc()

	// Sử dụng URI từ cấu hình mà không có mã hóa
	uri := mongoConfig.URI

	// Tạo tùy chọn client cho MongoDB từ package của bạn
	opts := mongo.NewClientOptions().ApplyURI(uri)

	// Kết nối tới MongoDB
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping tới database để kiểm tra kết nối
	err = client.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB!")

	return client, nil
}

// Disconnect disconnects from the MongoDB database
func Disconnect(client mongo.Client) {
	if client == nil {
		return
	}

	// Ngắt kết nối MongoDB
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	}

	log.Println("Connection to MongoDB closed.")
}
