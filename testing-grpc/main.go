package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	myapi "src/testing-grpc/test"
)

type server struct {
	mongoClient *mongo.Client
	myapi.UnimplementedMyServiceServer
}

func (s *server) CreateRecord(ctx context.Context, req *myapi.CreateRecordRequest) (*myapi.CreateRecordResponse, error) {
	collection := s.mongoClient.Database("testdb").Collection("user")

	record := bson.M{
		"name": req.GetName(),
		"age":  req.GetAge(),
	}

	result, err := collection.InsertOne(ctx, record)
	if err != nil {
		log.Printf("Failed to create record: %v", err)
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return &myapi.CreateRecordResponse{Id: id}, nil
}

func (s *server) UpdateRecord(ctx context.Context, req *myapi.UpdateRecordRequest) (*myapi.UpdateRecordResponse, error) {
	collection := s.mongoClient.Database("testdb").Collection("user")
	fmt.Println("print the req", req, req.GetId())

	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return nil, err
	}

	filter := bson.M{"_id": bson.M{"$eq": id}} //specifies that you want to match the document with an id equal to the value stored in the id variable.
	//	update := bson.M{"$set": bson.M{"name": req.GetName(), "age": req.GetAge()}}

	update := bson.M{}

	if req.GetName() != "" {
		update["$set"] = bson.M{"name": req.GetName()}
	}

	if req.GetAge() != 0 {
		update["$set"] = bson.M{"age": req.GetAge()}
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to update record: %v", err)
		return nil, err
	}
	// 	UpdateOne() takes three parameters:
	// ctx: A context.Context object that provides context for the operation, including any deadlines or cancellation signals.
	// filter: A document that specifies the criteria for matching the document to update. It defines which document to update based on the specified conditions.
	// update: A document that specifies the modifications to apply to the matched document. It defines the changes you want to make to the matched document.

	return &myapi.UpdateRecordResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	//clientOptions := options.Client().ApplyURI("mongodb+srv://apurvasaini:apps0602@cluster0.czgb2x9.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	//clientOptions := options.Client().ApplyURI("mongodb+srv://apurvasaini:apps0602@cluster0.czgb2x9.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return
	}
	fmt.Println("Connected to MongoDB!")
	//mongoClient, err := getMongoClient()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB: %v", err)
	// }

	s := grpc.NewServer()
	myapi.RegisterMyServiceServer(s, &server{mongoClient: client})
	reflection.Register(s)

	fmt.Println("gRPC server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
