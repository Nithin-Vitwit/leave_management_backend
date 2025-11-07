package main

import "time"

const dateFormat = "02-01-2006"

type Emp struct {
	Name string `json:"name" bson:"name"`
	ID   string `json:"id" bson:"id"`
}

type Hr struct {
	Name string `json:"name" bson:"name"`
	ID   string `json:"id" bson:"id"`
}

type Leave struct {
	Name     string    `json:"name" bson:"name"`
	EmpID    string    `json:"emp_id" bson:"emp_id"`
	Reason   string    `json:"reason" bson:"reason"`
	FromDate time.Time `json:"from_date" bson:"from_date"`
	ToDate   time.Time `json:"to_date" bson:"to_date"`
	Status   string    `json:"status" bson:"status"`
}
