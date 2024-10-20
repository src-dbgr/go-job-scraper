package mongodb

import (
	"context"
	"job-scraper/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// JobRepository definiert das Interface für Datenbankoperationen
type JobRepository interface {
	SaveJob(ctx context.Context, job models.Job) error
	GetJobByID(ctx context.Context, id string) (*models.Job, error)
	GetJobs(ctx context.Context) ([]models.Job, error)
	GetJobCountByCategory(ctx context.Context) (map[string]int, error)
	GetTotalJobCount(ctx context.Context) (int, error)
	GetExistingURLs(ctx context.Context) (map[string]bool, error)
	AggregateJobs(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error)
}

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient(ctx context.Context, uri, dbName string) (*Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Client{
		client: client,
		db:     db,
	}, nil
}

func (c *Client) GetJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job
	cursor, err := c.db.Collection("jobs").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (c *Client) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	var job models.Job
	err := c.db.Collection("jobs").FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (c *Client) SaveJob(ctx context.Context, job models.Job) error {
	_, err := c.db.Collection("jobs").InsertOne(ctx, job)
	return err
}

func (c *Client) GetJobCountByCategory(ctx context.Context) (map[string]int, error) {
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$jobCategories", "count": bson.M{"$sum": 1}}},
	}
	cursor, err := c.db.Collection("jobs").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, result := range results {
		countMap[result.ID] = result.Count
	}
	return countMap, nil
}

func (c *Client) GetTotalJobCount(ctx context.Context) (int, error) {
	count, err := c.db.Collection("jobs").CountDocuments(ctx, bson.M{})
	return int(count), err
}

func (c *Client) GetExistingURLs(ctx context.Context) (map[string]bool, error) {
	cursor, err := c.db.Collection("jobs").Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"url": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	urls := make(map[string]bool)
	for cursor.Next(ctx) {
		var job struct {
			URL string `bson:"url"`
		}
		if err := cursor.Decode(&job); err != nil {
			return nil, err
		}
		urls[job.URL] = true
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) AggregateJobs(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
	cursor, err := c.db.Collection("jobs").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
