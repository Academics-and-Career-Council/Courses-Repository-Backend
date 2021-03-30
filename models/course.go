package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Course Structure
type MGMCourse struct {
	ID          string               `bson:"_id,omitempty"`
	Title       string               `bson:"title"`
	Number      string               `bson:"number"`
	Credits     string               `bson:"credits"`
	Offered     string               `bson:"offered"`
	Contents    string               `bson:"contents"`
	DriveFolder string               `bson:"driveFolder"`
	DriveFiles  []string             `bson:"driveFiles"`
	Reviews     []primitive.ObjectID `bson:"reviews"`
	Author      string               `bson:"author"`
	Dept        string               `bson:"dept"`
}
