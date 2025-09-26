package repositories

import (
	"context"
	"project/internal/database"
	"project/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct{}

func NewUserRepositoryMongo() *UserRepositoryMongo { return &UserRepositoryMongo{} }

func (r *UserRepositoryMongo) col() *mongo.Collection {
    return database.GetMongoDB().Collection("users")
}

func (r *UserRepositoryMongo) FindAll() ([]models.UserResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    cur, err := r.col().Find(ctx, bson.M{"deleted_at": bson.M{"$exists": false}})
    if err != nil { return nil, err }
    defer cur.Close(ctx)
    var res []models.UserResponse
    for cur.Next(ctx) {
        var u models.User
        if err := cur.Decode(&u); err != nil { return nil, err }
        res = append(res, models.UserResponse{ID: u.ID.Hex(), Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt})
    }
    return res, cur.Err()
}

func (r *UserRepositoryMongo) FindByID(idStr string) (*models.UserResponse, error) {
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil { return nil, err }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    var u models.User
    err = r.col().FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}).Decode(&u)
    if err != nil { return nil, err }
    resp := &models.UserResponse{ID: u.ID.Hex(), Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
    return resp, nil
}

func (r *UserRepositoryMongo) FindByEmail(email string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    var u models.User
    err := r.col().FindOne(ctx, bson.M{"email": email, "deleted_at": bson.M{"$exists": false}}).Decode(&u)
    if err != nil { return nil, err }
    return &u, nil
}

func (r *UserRepositoryMongo) Create(user models.User) (*models.UserResponse, error) {
    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _, err := r.col().InsertOne(ctx, user)
    if err != nil { return nil, err }
    resp := &models.UserResponse{ID: user.ID.Hex(), Name: user.Name, Email: user.Email, CreatedAt: user.CreatedAt}
    return resp, nil
}

func (r *UserRepositoryMongo) Update(idStr string, user models.User) (*models.UserResponse, error) {
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil { return nil, err }
    update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
    if user.Name != "" { update["$set"].(bson.M)["name"] = user.Name }
    if user.Email != "" { update["$set"].(bson.M)["email"] = user.Email }
    if user.Password != "" { update["$set"].(bson.M)["password"] = user.Password }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    _, err = r.col().UpdateByID(ctx, id, update)
    if err != nil { return nil, err }
    return r.FindByID(idStr)
}

func (r *UserRepositoryMongo) Delete(idStr string) error {
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil { return err }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    now := time.Now()
    _, err = r.col().UpdateByID(ctx, id, bson.M{"$set": bson.M{"deleted_at": now}})
    return err
}


