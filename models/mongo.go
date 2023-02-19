package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoProfessor struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	RMPId       string             `json:"rmpId" bson:"rmpId"`
	Rating      float64            `json:"rating" bson:"rating"`
	Difficulty  float64            `json:"difficulty" bson:"difficulty"`
	Departments []MongoDepartment  `json:"department" bson:"department"`
	Reviews     []MongoReview      `json:"reviews" bson:"reviews"`
	Courses     []MongoCourse      `json:"courses" bson:"courses"`
}

type MongoReview struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProfessorID string             `json:"professorId" bson:"professorId"`
	Professor   string             `json:"professor" bson:"professor"`
	Quality     float64            `json:"quality" bson:"quality"`
	Difficulty  float64            `json:"difficulty" bson:"difficulty"`
	Date        string             `json:"date" bson:"date"`
	ReviewText  string             `json:"reviewText" bson:"reviewText"`
	Course      MongoCourse        `json:"course" bson:"course"`
	Helpful     float64            `json:"helpful" bson:"helpful"` // quality = helpful+clarity/2
	Clarity     float64            `json:"clarity" bson:"clarity"`
}

type MongoCourse struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProfessorID string             `json:"professorId" bson:"professorId"`
	Department  string             `json:"department" bson:"department"`
	Number      string             `json:"number" bson:"number"`
}

type MongoDepartment struct {
	ID                   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                 string             `json:"name" bson:"name"`
	DepartmentBase64Code string             `json:"departmentBase64Code" bson:"departmentBase64Code"`
}
