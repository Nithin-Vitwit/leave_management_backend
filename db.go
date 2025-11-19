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
	// Load .env ONLY in local (not on Render)
	if os.Getenv("RENDER") == "" {
		_ = godotenv.Load()
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	hrPassword = os.Getenv("HR_PASSWORD")

	if mongoURI == "" || dbName == "" {
		log.Fatal("❌ Missing MongoDB environment variables")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	db = client.Database(dbName)
	employeeCol = db.Collection("employees")
	hrCol = db.Collection("hrs")
	leaveCol = db.Collection("leaves")

	fmt.Println("✅ Connected to MongoDB Atlas!")
}
