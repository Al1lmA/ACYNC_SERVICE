package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const password = "12345"

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/rating", func(c *gin.Context) {
		var data RatingData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		RequestID := data.RequestID

		// Запуск горутины для отправки статуса
		go sendRating(RequestID, password, fmt.Sprintf("http://localhost:8000/requests/%d/rating/", RequestID))

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	router.Run(":5000")

	log.Println("Server down")
}

func genRandomRating(password string) Result {
	time.Sleep(8 * time.Second)

	rand.Seed(time.Now().UnixNano())
	randomRating := strconv.Itoa(rand.Intn(851))
	fmt.Println(randomRating)

	return Result{randomRating, password}
}

// Функция для отправки статуса в отдельной горутине
func sendRating(RequestID int, password string, url string) {
	// Выполнение расчётов с randomStatus
	result := genRandomRating(password)

	// Отправка PUT-запроса к основному серверу
	_, err := performPUTRequest(url, result)
	if err != nil {
		fmt.Println("Error sending Name:", err)
		return
	}

	fmt.Println("Name sent successfully for RequestID:", RequestID)
}

type Result struct {
	Rating   string `json:"rating"`
	Password string `json:"password"`
}

type RatingData struct {
	RequestID int `json:"request_id"`
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Создание PUT-запроса
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
