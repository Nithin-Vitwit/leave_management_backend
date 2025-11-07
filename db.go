package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client      *mongo.Client
	db          *mongo.Database
	employeeCol *mongo.Collection
	hrCol       *mongo.Collection
	leaveCol    *mongo.Collection
	hrPassword  string
)

func initDB() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	hrPassword = os.Getenv("HR_PASSWORD")

	if mongoURI == "" || dbName == "" {
		log.Fatal("❌ Missing MongoDB configuration in .env")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	db = client.Database(dbName)
	employeeCol = db.Collection("employees")
	hrCol = db.Collection("hrs")
	leaveCol = db.Collection("leaves")

	fmt.Println("✅ Connected to MongoDB Atlas!")
}
