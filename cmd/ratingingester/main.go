package ratingingester

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	model "movieexample.com/rating/pkg"
	"os"
	"time"
)

func main() {
	fmt.Println("Creating a Kafka producer")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const fileName = "ratings_data.json"
	fmt.Println("Reading rating events from file " + fileName)

	ratingEvents, err := readRatingEvents(fileName)
	if err != nil {
		panic(err)
	}

	const topic = "ratings"
	if err := produceRatingEvents(topic, producer, ratingEvents); err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second
	fmt.Printf("Waiting %v until all events get produced", timeout.String())

	producer.Flush(int(timeout.Milliseconds()))
}

func produceRatingEvents(topic string, producer *kafka.Producer, ratingEvents []model.RatingEvent) error {
	for _, ratingEvent := range ratingEvents {
		encodedEvent, err := json.Marshal(ratingEvent)
		if err != nil {
			return err
		}

		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          encodedEvent,
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}
