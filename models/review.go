package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Review Structure
type MGMReview struct {
	ID        string             `bson:"_id,omitempty"`
	Course    primitive.ObjectID `bson:"course"`
	Semester  string             `bson:"semester"`
	Professor string             `bson:"professor"`
	Contents  string             `bson:"contents"`
	Author    string             `bson:"author"`
	Approved  bool               `bson:"approved"`
}
