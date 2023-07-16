package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/7leven7/xm-company-crud/app/database"
	"github.com/7leven7/xm-company-crud/app/models"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedCompany, err := company.SaveCompany()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		err = createKafkaTopic("company_created")
		if err != nil {
			fmt.Println("Error creating Kafka topic:", err)
			return
		}

		writer := &kafka.Writer{
			Addr:  kafka.TCP("kafka:9092"),
			Topic: "company_created",
		}

		message := kafka.Message{
			Key:   []byte(savedCompany.ID),
			Value: []byte(savedCompany.ID),
		}

		err = writer.WriteMessages(context.Background(), message)
		if err != nil {
			fmt.Println("Error writing Kafka message:", err)
			return
		}
	}()

	c.JSON(http.StatusOK, savedCompany)
}

func UpdateCompany(c *gin.Context) {
	companyID := c.Param("id")

	var updatedCompany models.Company
	if err := c.ShouldBindJSON(&updatedCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var company models.Company
	err := database.DB.First(&company, "id = ?", companyID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	company.ID = companyID
	updatedCompany.ID = companyID

	_, err = company.UpdateCompany(&updatedCompany)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		err = createKafkaTopic("company_updated")
		if err != nil {
			fmt.Println("Error creating Kafka topic:", err)
			return
		}

		writer := &kafka.Writer{
			Addr:  kafka.TCP("kafka:9092"),
			Topic: "company_updated",
		}

		message := kafka.Message{
			Key:   []byte(updatedCompany.ID),
			Value: []byte(updatedCompany.ID),
		}

		err = writer.WriteMessages(context.Background(), message)
		if err != nil {
			fmt.Println("Error writing Kafka message:", err)
			return
		}
	}()

	c.JSON(http.StatusOK, updatedCompany)
}

func DeleteCompany(c *gin.Context) {
	companyID := c.Param("id")

	var company models.Company

	err := database.DB.First(&company, "id = ?", companyID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	err = company.DeleteCompany()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}
	go func() {
		err = createKafkaTopic("company_deleted")
		if err != nil {
			fmt.Println("Error creating Kafka topic:", err)
			return
		}

		writer := &kafka.Writer{
			Addr:  kafka.TCP("kafka:9092"),
			Topic: "company_deleted",
		}

		message := kafka.Message{
			Key:   []byte(companyID),
			Value: []byte(companyID),
		}

		err = writer.WriteMessages(context.Background(), message)
		if err != nil {
			fmt.Println("Error writing Kafka message:", err)
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}

func GetCompanyByID(c *gin.Context) {
	companyID := c.Param("id")
	company, err := models.GetCompanyByID(companyID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func createKafkaTopic(topic string) error {
	broker := "kafka:9092"
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "topic-checker",
		Topic:   topic,
	})

	_, err := r.FetchMessage(context.Background())
	if err != nil {
		if err.Error() == "kafka: topic not found" {
			conn, err := kafka.Dial("tcp", broker)
			if err != nil {
				return err
			}
			defer conn.Close()

			err = conn.CreateTopics(kafka.TopicConfig{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 3,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
