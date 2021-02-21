package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mo2/dto"
	"mo2/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var blogCol *mongo.Collection = GetCollection("blog")

func ensureBlogIndex() {
	blogCol.Indexes().CreateMany(context.TODO(), append([]mongo.IndexModel{
		{Keys: bson.M{"ket_words": 1}},
	}, model.IndexModels...))
}

// AddBlog add
func AddBlog(b *model.Blog) (new bool, err error) {
	entity := model.InitEntity()
	b.EntityInfo = entity
	result, err := blogCol.UpdateOne(
		context.TODO(),
		bson.D{{"_id", b.ID}},
		bson.D{{"$set", bson.M{
			"title":       b.Title,
			"description": b.Description,
			"content":     b.Content,
			"cover":       b.Cover,
			"key_words":   b.KeyWords,
		}}},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		new = false // already exists
	}
	if result.UpsertedCount != 0 {
		new = true //create new blog
		b.ID = result.UpsertedID.(primitive.ObjectID)
	}
	return
}

//find blog
func FindBlogs(u dto.LoginUserInfo) (b []model.Blog) {
	opts := options.Find().SetSort(bson.D{{"entity_info", 1}})
	cursor, err := blogCol.Find(context.TODO(), bson.D{{"author_id", u.ID}}, opts)
	err = cursor.All(context.TODO(), &b)
	if err != nil {
		log.Fatal(err)
	}
	return
}
