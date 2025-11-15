package data

import (
	"context"
	"errors"
	"strconv"
	"time"
	"github.com/hababisha/enhancedTaskManager/models"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTaskService struct{
	col *mongo.Collection
	timeout time.Duration
}

func NewMongoTaskService(uri, dbName, colName string, timeout time.Duration) (*MongoTaskService, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	//ping to check
	if err := client.Ping(ctx, nil); err != nil{
		return nil, err
	}

	col := client.Database(dbName).Collection(colName)
	svc := &MongoTaskService{col: col, timeout: timeout}

	return svc, nil
}

type dbTask struct{
	ID primitive.ObjectID `bson: "_id, omitempty"`
	Title string `bson: "title"`
	Description string `bson: "description, omitempty"`
	Status string `bson: "status, omitempty`
}

//helper (converign dbTask -> models.task)
func dbToModel(t dbTask) models.Task{
	return models.Task{
		ID: t.ID.Hex(),
		Title: t.Title,
		Description: t.Description,
		Status: t.Status,
	}
}

func (m *MongoTaskService) GetAllTasks() ([]models.Task, error){
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	cur, err := m.col.Find(ctx, bson.M{})
	if err != nil{
		return nil, err
	}

	defer cur.Close(ctx)

	var out []models.Task
	for cur.Next(ctx){
		var b dbTask

		if err := cur.Decode(&b); err != nil{
			return nil, err
		}
		out = append(out, dbToModel(b))
	}
	if err := cur.Err(); err != nil{
		return nil, err
	}

	return out, nil
}


func (m *MongoTaskService) GetTaskByID(id string) (models.Task, error){
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	var filter bson.M
	if err == nil{
		filter = bson.M{"_id": oid}
	} else {
		if intVal, err2 := strconv.Atoi(id); err2 == nil{
			filter = bson.M{"id": intVal}
		} else {
			return models.Task{}, mongo.ErrNoDocuments
		}
	}

	var b dbTask
	if err := m.col.FindOne(ctx, filter).Decode(&b); err != nil{
		return models.Task{}, err
	}
	return dbToModel(b), nil
}

//create task

func (m *MongoTaskService) CreateTask(task models.Task) (models.Task, error){
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)

	defer cancel()
	doc := bson.M{
		"title": task.Title,
		"desctiption": task.Description,
		"status": task.Status,
	}

	res, err := m.col.InsertOne(ctx, doc)
	if err != nil{
		return models.Task{}, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok{
		return models.Task{}, errors.New("failed to convert inseretd id to objectId")
	}
	task.ID = oid.Hex()
	return task, nil
}

func (m *MongoTaskService) UpdateTask(id string, updated models.Task) (models.Task, error) {
    ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
    defer cancel()

    oid, err := primitive.ObjectIDFromHex(id)
    var filter bson.M
    if err == nil {
        filter = bson.M{"_id": oid}
    } else {
        // try numeric id field fallback (if used)
        if intVal, err2 := strconv.Atoi(id); err2 == nil {
            filter = bson.M{"id": intVal}
        } else {
            return models.Task{}, mongo.ErrNoDocuments
        }
    }

    update := bson.M{}
    set := bson.M{}
    if updated.Title != "" {
        set["title"] = updated.Title
    }
    if updated.Description != "" {
        set["description"] = updated.Description
    }
    if updated.Status != "" {
        set["status"] = updated.Status
    }
    if len(set) == 0 {
        // nothing to update
        return models.Task{}, errors.New("no fields to update")
    }
    update["$set"] = set

    opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
    var res dbTask
    if err := m.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&res); err != nil {
        return models.Task{}, err
    }
    return dbToModel(res), nil
}

// DeleteTask deletes by id
func (m *MongoTaskService) DeleteTask(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
    defer cancel()

    oid, err := primitive.ObjectIDFromHex(id)
    var filter bson.M
    if err == nil {
        filter = bson.M{"_id": oid}
    } else {
        if intVal, err2 := strconv.Atoi(id); err2 == nil {
            filter = bson.M{"id": intVal}
        } else {
            return mongo.ErrNoDocuments
        }
    }

    res, err := m.col.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return mongo.ErrNoDocuments
    }
    return nil
}