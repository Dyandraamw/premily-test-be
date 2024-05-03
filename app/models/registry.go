package models

// import "github.com/frangklynndruru/premily_backend/app/models"

type Model struct {
	Model interface{}
}

func RegisterModels() []Model {
	// return []Model{
	// 	{Model: models.User{}},
	// 	{Model: models.Invoice{}},
	// 	{Model: models.Description_Details{}},
	// 	{Model: models.Installment{}},
	// 	{Model: models.Sum_Insured_Details{}},
	// 	{Model: models.Payment_Details{}},
	// 	{Model: models.Adjustment{}},
	// 	{Model: models.Payment_Status{}},
	// 	{Model: models.Statement_Of_Account{}},
	// 	{Model: models.Statement_Of_Account_Details{}},
	// }
	return []Model{
		{Model: User{}},
		{Model: Invoice{}},
		{Model: Description_Details{}},
		{Model: Installment{}},
		{Model: Sum_Insured_Details{}},
		{Model: Payment_Details{}},
		{Model: Adjustment{}},
		{Model: Payment_Status{}},
		{Model: Statement_Of_Account{}},
		{Model: Statement_Of_Account_Details{}},
	}
}