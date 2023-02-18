package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Professor struct {
	gorm.Model
	Name        string       `json:"name" bson:"name"`
	RMPId       string       `json:"rmpId" bson:"rmpId"`
	Rating      float64      `json:"rating" bson:"rating"`
	Difficulty  float64      `json:"difficulty" bson:"difficulty"`
	Departments []Department `json:"department" bson:"department"`
	Reviews     []Review     `json:"reviews" bson:"reviews"`
	Courses     []Course     `json:"courses" bson:"courses"`
}

type Review struct {
	gorm.Model
	ProfessorID string  `json:"professorId" gorm:"column:professor_id" bson:"professorId"`
	Professor   string  `json:"professor" bson:"professor"`
	Quality     float64 `json:"quality" bson:"quality"`
	Difficulty  float64 `json:"difficulty" bson:"difficulty"`
	Date        string  `json:"date" bson:"date"`
	ReviewText  string  `json:"reviewText" bson:"reviewText"`
	Course      Course  `json:"course" bson:"course"`
	Helpful     float64 `json:"helpful" bson:"helpful"` // quality = helpful+clarity/2
	Clarity     float64 `json:"clarity" bson:"clarity"`
}

type Course struct {
	gorm.Model
	ProfessorID string `json:"professorId" gorm:"column:professor_id" bson:"professorId"`
	Department  string `json:"department" bson:"department"`
	Number      string `json:"number" bson:"number"`
}

type Department struct {
	gorm.Model
	Name                 string `json:"name" bson:"name"`
	DepartmentBase64Code string `json:"departmentBase64Code" bson:"departmentBase64Code"`
}
