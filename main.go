package main

import (
	"context"
	"fmt"
	"github.com/FerretDB/FerretDB/ferretdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	done chan struct{}
)

const (
	mgURL      = "127.0.0.1:27019" // we are opening this port for listening
	pgURL      = "127.0.0.1:5432"  // this port has to be already opened
	pgDatabase = "postgres"        // "ferretdb"
	pgUsername = "postgres"
	pgPassword = "mypass"
)

func connect(ctx context.Context) (*mongo.Client, error) {
	f, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: mgURL,
		},
		Handler:       "postgresql",
		PostgreSQLURL: fmt.Sprintf("postgres://%s/%s", pgURL, pgDatabase),
	})
	if err != nil {
		return nil, err
	}

	go func() {
		log.Print(f.Run(ctx))
		close(done)
	}()

	uri := f.MongoDBURI()
	fmt.Println("Mongo connection URI = ", uri)

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(options.Credential{
		AuthMechanism: "PLAIN",
		Username:      pgUsername,
		Password:      pgPassword,
		AuthSource:    "$external", // default value. ctrl + click on options.Credential for more details.
	}))
	if err != nil {
		return nil, fmt.Errorf("con = %v", err)
	}
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("ping = %v", err)
	}
	fmt.Println("Pinging Successful")
	return c, nil
}

func insertData(client *mongo.Client) error {
	collection := client.Database("test").Collection("example")
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	fmt.Println("Data inserted successfully!")
	return nil
}

func main() {
	done = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = insertData(c)
	//if err != nil {
	//	fmt.Println("insert", err)
	//	return
	//}

	accessData(ctx, c)
	<-done
}

func accessData(ctx context.Context, c *mongo.Client) {
	databaseName := "test"
	collectionName := "example"

	lst, err := c.ListDatabases(ctx, bson.D{{}})
	fmt.Println(lst.Databases)

	db := c.Database(databaseName)
	cl, e := db.ListCollectionNames(ctx, bson.D{{}})
	if e != nil {
		fmt.Println("c", e)
	}
	fmt.Println(cl)
	collection := db.Collection(collectionName)

	filter := bson.D{{}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	results := []map[string]interface{}{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	cursor.Close(context.Background())
	fmt.Println("Found documents:", len(results))
}
