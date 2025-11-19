package main

import (
	"context"
	"fmt"
	
	"os"
	"time"

	// Removed: "github.com/joho/godotenv"
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

// initDB is modified to return an error instead of using log.Fatal
// This allows the caller (main) to log the error gracefully and exit.
func initDB() error {
	// *** RENDER FIX: Removed godotenv.Load(), relying on OS environment variables ***

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	hrPassword = os.Getenv("HR_PASSWORD")

	if mongoURI == "" || dbName == "" {
		// Log the error and return it to main()
		return fmt.Errorf("❌ Missing MongoDB configuration. Ensure MONGO_URI and DB_NAME are set in environment variables")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		// Log the connection error and return it
		return fmt.Errorf("MongoDB connection failed: %w", err)
	}

	// Ping the database to ensure connection is live
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	db = client.Database(dbName)
	employeeCol = db.Collection("employees")
	hrCol = db.Collection("hrs")
	leaveCol = db.Collection("leaves")

	fmt.Println("✅ Successfully connected to MongoDB.")
	return nil
}