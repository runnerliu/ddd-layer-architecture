package persistence

import (
	"context"
	"ddd-demo/common/factory"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// GetMongoDBCollection 获取 MongoDB 集合
func GetMongoDBCollection(factory *factory.ObjectFactory, uri, databaseName, collectionName string) (*mongo.Collection, error) {
	client, err := factory.Get(
		uri,
		func(uri string) (interface{}, error) {
			opts := options.Client()
			opts.Monitor = otelmongo.NewMonitor()
			return mongo.Connect(context.Background(), opts.ApplyURI(uri))
		},
		func(obj interface{}) error {
			client, ok := obj.(*mongo.Client)
			if !ok {
				return nil
			}
			return client.Disconnect(context.Background())
		})
	if err != nil {
		return nil, err
	}
	return client.(*mongo.Client).Database(databaseName).Collection(collectionName), nil
}
