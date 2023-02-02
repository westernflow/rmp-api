package handler

import (
	"rmpParser/controller"
	model "rmpParser/models"
)

func GetProfessors() []model.Professor {
	c := controller.GetInstance()
	return c.GetProfessors()
}

func GetDepartments() []model.Department {
	c := controller.GetInstance()
	return c.GetDepartments()
}
