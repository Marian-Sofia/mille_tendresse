package users_repository

import (
	"context"
	"time"

	users_interfaces "github.com/Marian-Sofia/mille_tendresse/users/internal/interfaces"
	users_model "github.com/Marian-Sofia/mille_tendresse/users/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type usersRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) users_interfaces.IUsersRepository {
	return &usersRepository{
		collection: db.Collection("users"),
	}
}

func (rpt *usersRepository) Find() ([]users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []users_model.User

	cursor, err := rpt.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (rpt *usersRepository) FindById(userId string) (users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user users_model.User

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return user, err 
	}

	err = rpt.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, nil
		}
		return user, err
	}

	return user, nil
}

func (rpt *usersRepository) Create(userModel users_model.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := rpt.collection.InsertOne(ctx, userModel)
	if err != nil {
		return "", err
	}
	return "User has been created", nil
}

func (rpt *usersRepository) Update(userId string, updates map[string]interface{}) (users_model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var updateUser users_model.User

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return users_model.User{}, err
	}

	update := bson.M {
		"$set": updates,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = rpt.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updateUser)

	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (rpt *usersRepository) Delete(userId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", err 
	}

	_, err = rpt.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}

	return "the user has been deleted", err
}