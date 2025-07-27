package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zl0y-billing/internal/database"
	"zl0y-billing/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReportRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewReportRepository(db *database.MongoDB) *ReportRepository {
	return &ReportRepository{
		db:         db,
		collection: db.Database.Collection("reports"),
	}
}

func (r *ReportRepository) CreateReport(clientGeneratedID string) (*models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	report := &models.Report{
		ID:                primitive.NewObjectID(),
		ReportID:          primitive.NewObjectID().Hex(), // Generate a unique report ID
		ClientGeneratedID: clientGeneratedID,
		IsPurchased:       false,
		CreatedAt:         time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

func (r *ReportRepository) LinkAnonymousReport(clientGeneratedID string, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find reports with matching client_generated_id and no user_id
	filter := bson.M{
		"client_generated_id": clientGeneratedID,
		"user_id":             bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{
			"user_id": userID,
		},
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to link anonymous report: %w", err)
	}

	return int(result.ModifiedCount), nil
}

func (r *ReportRepository) GetReportsByUserID(userID int, limit, offset int) ([]models.Report, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reports: %w", err)
	}

	// Find reports with pagination
	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find reports: %w", err)
	}
	defer cursor.Close(ctx)

	var reports []models.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, 0, fmt.Errorf("failed to decode reports: %w", err)
	}

	return reports, total, nil
}

func (r *ReportRepository) GetReportByID(reportID string) (*models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"report_id": reportID}

	var report models.Report
	err := r.collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("report not found")
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	return &report, nil
}

func (r *ReportRepository) MarkReportAsPurchased(reportID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"report_id": reportID}
	update := bson.M{"$set": bson.M{"is_purchased": true}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark report as purchased: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("report not found or already purchased")
	}

	return nil
}

func (r *ReportRepository) GetReportsByClientID(clientGeneratedID string) ([]models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"client_generated_id": clientGeneratedID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find reports: %w", err)
	}
	defer cursor.Close(ctx)

	var reports []models.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, fmt.Errorf("failed to decode reports: %w", err)
	}

	return reports, nil
}
