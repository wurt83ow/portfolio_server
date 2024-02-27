package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	ID      int    `json:"id" bson:"id"`
	Src      string  `json:"src" bson:"src"`
	Alt      string  `json:"alt" bson:"alt"`
	Href     *string `json:"href,omitempty" bson:"href,omitempty"`
	HrefText *string `json:"hrefText,omitempty" bson:"hrefText,omitempty"`
	Rating   *int    `json:"rating,omitempty" bson:"rating,omitempty"`
}

type Section struct {
	ID      int    `json:"id" bson:"id"`
	NClass  string `json:"nclass" bson:"nclass"`
	Title   string `json:"title" bson:"title"`
	Content struct {
		TextBefore string  `json:"textBefore" bson:"textBefore"`
		IClass  string `json:"iclass" bson:"iclass"`
		Items      *[]Item `json:"items,omitempty" bson:"items,omitempty"`		 
		TextAfter  *string `json:"textAfter,omitempty" bson:"textAfter,omitempty"`
		IsActive   *bool   `json:"isActive,omitempty" bson:"isActive,omitempty"`
	} `json:"content" bson:"content"`
	IsActive bool `json:"isActive,omitempty" bson:"isActive,omitempty"`
}


func main() {
	// Установите контекст для подключения к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 160*time.Second)
	defer cancel()

	// Подключитесь к MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@127.0.0.1:27017"))
	if err != nil {
		log.Println("Ошибка подключения к MongoDB:", err)
		return
	}

	// Получите коллекцию
	collection := client.Database("test").Collection("sections")
	
	// Загрузите данные из файла sectionsData.json
	data, err := os.ReadFile("../../sectionsData.json")
	if err != nil {
		log.Println("Ошибка чтения файла sectionsData.json:", err)
		return
	}

	var sectionsData []Section
	err = json.Unmarshal(data, &sectionsData)
	if err != nil {
		log.Println("Ошибка преобразования данных из файла sectionsData.json:", err)
		return
	}

	// Добавьте или обновите данные в базе данных
	for _, section := range sectionsData {
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"id": section.ID},
			bson.M{"$set": section},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Println("Ошибка обновления данных в базе данных:", err)
			return
		}
	}

	http.HandleFunc("/api/sections", func(w http.ResponseWriter, r *http.Request) {
		// Извлеките все документы из коллекции
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			log.Println("Ошибка извлечения данных из базы данных:", err)
			return
		}

		var sections []Section
		if err = cursor.All(ctx, &sections); err != nil {
			log.Println("Ошибка чтения данных из базы данных:", err)
			return
		}
		fmt.Println("777777777777777777777777", sections)
		// Отправьте данные обратно клиенту
		json.NewEncoder(w).Encode(sections)
	})

	// Добавьте поддержку CORS
	corsOrigins := handlers.AllowedOrigins([]string{"http://localhost:8080", "http://localhost:3000"})
	handler := http.DefaultServeMux // ваш обработчик запросов
	corsHandler := handlers.CORS(corsOrigins)(handler)

	log.Fatal(http.ListenAndServe(":8090", corsHandler))
}
