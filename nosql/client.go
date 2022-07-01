package nosql

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"market/config"
	"market/robot/client"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

func init() {
	db = Connect()
}

func Connect() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

type PostInfo struct {
	ObjId     string    `bson:"_id"`
	ID        int32     `bson:"id"`
	Json      string    `bson:"json"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func FIndOne() {
	table := db.Database("td_game").Collection("mission_conf")

	post := &PostInfo{}
	err := table.FindOne(context.Background(), bson.D{{"id", 100110}}).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", post)
	log.Printf("%+v", post.UpdatedAt)
}

func Upsert() {
	missions := config.LoadConfig()
	//mission := missions[0]
	//data := client.JsonString(mission)
	//log.Println(data)

	collections := db.Database("td_game").Collection("mission_conf")

	for _, mission := range missions {
		data := client.JsonString(mission)

		result, err := collections.UpdateOne(context.Background(),
			bson.D{{Key: "id", Value: mission.Id}},
			bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{
							Key:   "id",
							Value: mission.Id,
						},
						{"json", data},
						{"updated_at", time.Now()},
					},
				},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", result)
	}
}

func InsertOne() {
	missions := config.LoadConfig()

	collections := db.Database("td").Collection("mission")

	for _, mission := range missions {
		js := client.JsonString(mission)
		var v interface{}
		if err := jsoniter.UnmarshalFromString(js, &v); err != nil {
			log.Fatal(err)
		}
		data, err := bson.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}

		result, err := collections.InsertOne(
			context.Background(),
			data,
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", result.InsertedID)
	}
}

func InsertMany() {

	var rows []interface{}
	missions := config.LoadConfig()
	for _, mission := range missions[:1] {
		js := client.JsonString(mission)
		var v interface{}
		if err := jsoniter.UnmarshalFromString(js, &v); err != nil {
			log.Fatal(err)
		}
		data, err := bson.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		rows = append(rows,
			MakeDocument(mission),
			MakeDocument(js),
			MakeDocument(data),
		)
	}

	collections := db.Database("td").Collection("mission")
	result, err := collections.InsertMany(
		context.Background(),
		rows,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", result)
}

func MakeDocument(v interface{}) interface{} {
	switch v.(type) {
	case string:
		var js = v.(string)
		var data interface{}
		if err := jsoniter.UnmarshalFromString(js, &data); err != nil {
			log.Fatal(err)
		}
		b, err := bson.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		return b

	case []byte:
		return v

	default:
		b, err := bson.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		return b
	}
}
