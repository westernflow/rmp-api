package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Professor struct {
	gorm.Model
	Name       string  `json:"name"`
	RMPId      string  `json:"rmpId"`
	Rating     float64 `json:"rating"`
	Difficulty float64 `json:"difficulty"`
	Departments []Department  `json:"department"`
	Reviews    []Review
	Courses    []Course
}

type Review struct {
	gorm.Model
	ProfessorID string  `json:"professorId" gorm:"column:professor_id"`
	Professor   string  `json:"professor"`
	Quality     float64 `json:"quality"`
	Difficulty  float64 `json:"difficulty"`
	Date        string  `json:"date"`
	ReviewText  string  `json:"reviewText"`
	Course      Course  `json:"course"`
	Helpful     float64 `json:"helpful"` // quality = helpful+clarity/2
	Clarity     float64 `json:"clarity"`
}

type Course struct {
	gorm.Model
	ProfessorID string `json:"professorId" gorm:"column:professor_id"`
	Department  string `json:"department"`
	Number      string `json:"number"`
}

type Department struct {
	gorm.Model
	Name                 string `json:"name"`
	DepartmentBase64Code string `json:"departmentBase64Code"`
}
