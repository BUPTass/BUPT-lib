package data

import (
	"time"
)

type ESource struct {
	Name         string    `bson:"Name"`
	IsActive     bool      `bson:"IsValid"`
	IsInTrial    bool      `bson:"IsInTrial"`
	TrialEndDate time.Time `bson:"TrialEndDate"`
	Subject      string    `bson:"Subject"`
	Link         string    `bson:"Link"`
	ResourceType string    `bson:"ResourceType"`
}

func UpdateFromXlsx() {}

func Add() {}

func Remove(id string) {}
