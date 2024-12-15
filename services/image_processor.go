package services

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/streadway/amqp"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/lib/pq"
)

type ImageProcessor struct {
	DB     *sql.DB
	S3     *s3.S3
	Queue  *amqp.Channel
	Logger *logrus.Logger
}

func NewImageProcessor(db *sql.DB, s3 *s3.S3, queue *amqp.Channel, logger *logrus.Logger) *ImageProcessor {
	return &ImageProcessor{
		DB:     db,
		S3:     s3,
		Queue:  queue,
		Logger: logger,
	}
}

func (ip *ImageProcessor) ProcessImages() {
	msgs, err := ip.Queue.Consume(
		"image_queue", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		ip.Logger.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		imageURL := string(msg.Body)
		ip.Logger.Infof("Processing image: %s", imageURL)

		compressedImageURL, err := ip.downloadAndCompressImage(imageURL)
		if err != nil {
			ip.Logger.Errorf("Failed to process image: %v", err)
			continue
		}

		err = ip.updateCompressedImageURLInDB(imageURL, compressedImageURL)
		if err != nil {
			ip.Logger.Errorf("Failed to update compressed image URL in DB: %v", err)
			continue
		}

		ip.Logger.Infof("Successfully processed image: %s", imageURL)
	}
}

func (ip *ImageProcessor) downloadAndCompressImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := io.Copy(gz, resp.Body); err != nil {
		return "", fmt.Errorf("failed to compress image: %v", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %v", err)
	}

	compressedImageKey := fmt.Sprintf("compressed/%s.gz", imageURL)
	_, err = ip.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(compressedImageKey),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload compressed image to S3: %v", err)
	}

	compressedImageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", os.Getenv("S3_BUCKET"), compressedImageKey)
	return compressedImageURL, nil
}

func (ip *ImageProcessor) updateCompressedImageURLInDB(originalImageURL, compressedImageURL string) error {
	query := `UPDATE products SET compressed_product_images = array_append(compressed_product_images, $1) 
			  WHERE $2 = ANY(product_images)`
	_, err := ip.DB.Exec(query, compressedImageURL, originalImageURL)
	if err != nil {
		return fmt.Errorf("failed to update compressed image URL in DB: %v", err)
	}
	return nil
}
