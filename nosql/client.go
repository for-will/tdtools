package nosql

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

type PostInfo struct {
	ObjId     string    `bson:"_id"`
	ID        int32     `bson:"id"`
	Json      string    `bson:"json"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type ObjectId = primitive.ObjectID

func NewOid(hex string) ObjectId {
	oid, _ := primitive.ObjectIDFromHex(hex)
	return oid
}

type FindMissionCond struct {
	Oid       ObjectId `bson:"_id,omitempty"`
	Id        int32    `bson:"id,omitempty"`
	Type      int32    `bson:"type,omitempty"`
	Condition int32    `bson:"condition,omitempty"`
}

func FindOne() {
	table := db.Database("td").Collection("mission")

	var data interface{}
	err := table.FindOne(
		context.Background(),
		MdbCond(FindMissionCond{
			//Oid:  NewOid("62beb74af43c378725e82674"),
			//Id:   100120,
			//Type: 1,
			Condition: 8,
		}),
	).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	js, err := bson.MarshalExtJSON(data, false, false)
	log.Println(string(js), err)

	kvs := data.(bson.D).Map()
	log.Printf("%+v", kvs["_id"].(primitive.ObjectID).Hex())
	//log.Println(client.JsonString(kvs))
	//log.Printf("%+v", post.UpdatedAt)
}

func FindMany() {
	table := db.Database("td").Collection("mission")

	cur, err := table.Find(
		context.Background(),
		MdbCond(FindMissionCond{
			//Oid:  NewOid("62beb74af43c378725e82674"),
			//Id:   100120,
			Condition: 6,
		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.Background()) {
		var data interface{}
		cur.Decode(&data)
		js, err := bson.MarshalExtJSON(data, false, false)

		kvs := data.(bson.D).Map()
		oid := kvs["_id"].(primitive.ObjectID).Hex()
		log.Println(oid, string(js), err)
	}

	//log.Println(client.JsonString(kvs))
	//log.Printf("%+v", post.UpdatedAt)
}

type MdbFilterBson struct {
	filter interface{}
}

func (filter *MdbFilterBson) MarshalBSON() ([]byte, error) {
	return bson.Marshal(filter.filter)
}

func MdbCond(cond interface{}) *MdbFilterBson {
	return &MdbFilterBson{filter: cond}
}

type MdbUpdateBson struct {
	Cmd struct {
		SET    interface{}         `bson:"$set,omitempty"`
		UNSET  map[string]struct{} `bson:"$unset,omitempty"`
		RENAME map[string]string   `bson:"$rename,omitempty"`
	}
}

func (upd *MdbUpdateBson) MarshalBSON() ([]byte, error) {
	return bson.Marshal(upd.Cmd)
}

func (upd *MdbUpdateBson) Set(v interface{}) *MdbUpdateBson {
	upd.Cmd.SET = v
	return upd
}

func (upd *MdbUpdateBson) Unset(unset ...string) *MdbUpdateBson {
	if upd.Cmd.UNSET == nil {
		upd.Cmd.UNSET = map[string]struct{}{}
	}
	for _, s := range unset {
		upd.Cmd.UNSET[s] = struct{}{}
	}
	return upd
}

func (upd *MdbUpdateBson) Rename(oldName string, newName string) *MdbUpdateBson {
	if upd.Cmd.RENAME == nil {
		upd.Cmd.RENAME = map[string]string{}
	}
	upd.Cmd.RENAME[oldName] = newName
	return upd
}

func MdbUpdate() *MdbUpdateBson {
	return &MdbUpdateBson{}
}

func UpdateOne() {
	missions := config.LoadConfig()

	collections := db.Database("td_game").Collection("mission_conf")

	for _, mission := range missions {

		result, err := collections.UpdateOne(
			context.Background(),
			MdbCond(struct {
				Id int
				//Condition int32
			}{Id: mission.Id}),
			MdbUpdate().
				//Unset("type", "thespoils").
				Set(mission),
			//Command().Unset("lv", "name", "sort").Rename("type", "typo"),
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
		result, err := collections.InsertOne(
			context.Background(),
			MakeDocument(mission),
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", result.InsertedID)
	}
}

func UniqueInsert(userId string, invitedBy string) {

	col := db.Database("patt").Collection("invite")

	type Invite struct {
		UserId    string `bson:"user_id,omitempty"`
		InvitedBy string `bson:"invited_by,omitempty"`
	}

	upd := Invite{
		UserId: userId,
	}
	result, err := col.UpdateOne(
		context.Background(),
		MdbCond(upd),
		MdbUpdate().Set(upd),
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", result)
	if result.UpsertedCount > 0 {
		col.UpdateOne(
			context.Background(),
			MdbCond(upd),
			MdbUpdate().Set(Invite{
				UserId:    userId,
				InvitedBy: invitedBy,
			}),
			options.Update().SetUpsert(false),
		)
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
