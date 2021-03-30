package services

import (
	"context"
	"log"

	"github.com/Mshivam2409/AnC-Courses/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	Users   *mongo.Database
	Courses *mongo.Database
}

// MongoClient ...
var MongoClient = &mongoClient{}

// ConnectMongo ..
func ConnectMongo() {
	MongoClient.Users = connect(viper.GetString("mongo.users"), "students")
	MongoClient.Courses = connect(viper.GetString("mongo.courses"), "primarydb")
}

func connect(url string, dbname string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB %v", err)
	}
	u := &models.MGMUser{}
	a := client.Database(dbname).Collection("ug").FindOne(context.Background(), bson.D{}).Decode(u)
	log.Print(a, u)
	log.Print(1000)
	log.Printf("Connected to MongoDB! URL : %s", url)
	database := client.Database(dbname)
	return database
}
