package mongoDatabase

import (
	"context"
	"errors"
	"time"

	"crud-app/models"
	"crud-app/request"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Queryis struct {
	db *mongo.Database
}

func (Queryis *Queryis) InsertUser(req request.CreateUserRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	if count, err := col.CountDocuments(ctx, bson.M{"email": req.Email, "isDeleted": false}); err != nil {
		return nil, err
	} else if count > 0 {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		IsDeleted: false,
	}

	res, err := col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (Queryis *Queryis) GetAllUsers(page, limit, minAge int64, nameContains string) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	skip := (page - 1) * limit

	matchStage := bson.M{
		"isDeleted": false,
	}

	if minAge > 0 {
		matchStage["age"] = bson.M{"$gt": minAge}
	}

	if nameContains != "" {
		matchStage["name"] = bson.M{"$regex": nameContains, "$options": "i"}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{
			Key: "$facet", Value: bson.M{
				"data": []bson.M{
					{"$skip": skip},
					{"$limit": limit},
				},
				"totalCount": []bson.M{
					{"$count": "count"},
				},
			},
		}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data       []models.User `bson:"data"`
		TotalCount []struct {
			Count int64 `bson:"count"`
		} `bson:"totalCount"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	if len(results) == 0 {
		return []models.User{}, 0, nil
	}

	users := results[0].Data
	total := int64(0)
	if len(results[0].TotalCount) > 0 {
		total = results[0].TotalCount[0].Count
	}

	return users, total, nil
}

func (Queryis *Queryis) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	col := models.GetUserCollection()

	if err := col.FindOne(ctx, bson.M{"_id": objID, "isDeleted": false}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (Queryis *Queryis) UpdateUser(id string, req request.UpdateUserRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	update := bson.M{"$set": bson.M{"name": req.Name, "email": req.Email, "age": req.Age}}
	result, err := col.UpdateOne(ctx, bson.M{"_id": objID, "isDeleted": false}, update)
	if err != nil {
		return "", err
	}

	if result.MatchedCount == 0 {
		newUser, err := Queryis.InsertUser(request.CreateUserRequest(req))
		if err != nil {
			return "", err
		}
		return "User not found. New user created with ID: " + newUser.ID.Hex(), nil
	}

	return "User successfully updated", nil
}

func SoftDeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := models.GetUserCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"isDeleted": true}})
	return err
}
