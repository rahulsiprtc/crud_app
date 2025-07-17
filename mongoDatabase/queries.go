package mongoDatabase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"crud-app/models"
	"crud-app/request"
	"crud-app/response"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Queryis struct {
	// db *mongo.Database
}

// func (q *Queryis) InsertUser(req request.CreateUserRequest) (response.UserResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	col := models.GetUserCollection()

// 	if count, err := col.CountDocuments(ctx, bson.M{"email": req.Email}); err != nil {
// 		return response.UserResponse{}, err
// 	} else if count > 0 {
// 		return response.UserResponse{}, errors.New("email already exists")
// 	}

// 	now := time.Now()

// 	user := &models.User{
// 		ID:    uuid.New().String(),
// 		Name:  req.Name,
// 		Email: req.Email,
// 		Age:   req.Age,
// 		// IsDeleted: false,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}
// 	user.ID = uuid.New().String()
// 	_, err := col.InsertOne(ctx, user)
// 	if err != nil {
// 		return response.UserResponse{}, err
// 	}

//		return response.UserResponse{
//			ID:    user.ID,
//			Name:  user.Name,
//			Email: user.Email,
//			Age:   user.Age,
//		}, nil
//	}
func (q *Queryis) InsertUser(req request.CreateUserRequest) (response.UserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	if count, err := col.CountDocuments(ctx, bson.M{"email": req.Email}); err != nil {
		return response.UserResponse{}, err
	} else if count > 0 {
		return response.UserResponse{}, errors.New("email already exists")
	}

	now := time.Now()

	user := &models.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := col.InsertOne(ctx, user)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}, nil
}

func (q *Queryis) GetAllUsers(req request.PaginationRequest) (response.UserPaginationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	matchStage := bson.M{"isDeleted": false}
	fmt.Print(req)

	if req.ID != "" {
		matchStage["id"] = req.ID
	}

	if req.MinAge > 0 {
		matchStage["age"] = bson.M{"$gt": req.MinAge}
	}
	if req.NameContains != "" {
		matchStage["name"] = bson.M{"$regex": req.NameContains, "$options": "i"}
	}

	skip := (req.Page - 1) * req.Limit

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{
			Key: "$facet", Value: bson.M{
				"data": []bson.M{
					{"$skip": skip},
					{"$limit": req.Limit},
				},
				"totalCount": []bson.M{
					{"$count": "count"},
				},
			},
		}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return response.UserPaginationResponse{}, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data       []models.User `bson:"data"`
		TotalCount []struct {
			Count int64 `bson:"count"`
		} `bson:"totalCount"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return response.UserPaginationResponse{}, err
	}

	var users []response.UserResponse
	total := int64(0)

	if len(results) > 0 {
		for _, u := range results[0].Data {
			users = append(users, response.UserResponse{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
				Age:   u.Age,
			})
		}
		if len(results[0].TotalCount) > 0 {
			total = results[0].TotalCount[0].Count
		}
	}

	lastPage := total / req.Limit
	if total%req.Limit != 0 {
		lastPage++
	}

	return response.UserPaginationResponse{
		Users: users,
		Pagination: response.Pagination{
			PerPage:      req.Limit,
			CurrentPage:  req.Page,
			LastPage:     lastPage,
			TotalResults: total,
		},
	}, nil
}

func (Queryis *Queryis) UpdateUser(id string, req request.UpdateUserRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	update := bson.M{"$set": bson.M{"name": req.Name, "email": req.Email, "age": req.Age}}
	// objID, err := primitive.ObjectIDFromHex(id)

	// result, err := col.UpdateOne(ctx, bson.M{"_id": objID}, update)
	result, err := col.UpdateOne(ctx, bson.M{"id": id}, update)

	if err != nil {
		return "", err
	}

	if result.MatchedCount == 0 {
		newUser, err := Queryis.InsertUser(request.CreateUserRequest(req))
		if err != nil {
			return "", err
		}
		return "User not found. New user created with ID: " + newUser.ID, nil
	}

	return "User successfully updated", nil
}

func (q *Queryis) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := models.GetUserCollection()

	// objID, err := primitive.ObjectIDFromHex(id)

	var user models.User
	if err := col.FindOne(ctx, bson.M{"id": id}).Decode(&user); err != nil {
		return err
	}

	if user.DeletedAt != nil {
		return errors.New("user already deleted")
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			// "isDeleted": true,
			"deletedAt": now,
			"updatedAt": now,
		},
	}

	_, err := col.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return err
	}
	return nil
}
