package app

import "github.com/frangklynndruru/premily_backend/app/models"

type Model struct {
	Model interface{}
}

func RegisterModels() []Model {
	return []Model{
		{Model: models.User{}},
		{Model: models.Invoice{}},
		{Model: models.Description_Details{}},
		{Model: models.Installment{}},
		{Model: models.Sum_Insured_Details{}},
		{Model: models.Payment_Details{}},
		{Model: models.Adjustment{}},
		{Model: models.Payment_Status{}},
		{Model: models.Statement_Of_Account{}},
		{Model: models.Statement_Of_Account_Details{}},
	}
}