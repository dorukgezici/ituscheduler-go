package main

import (
	"gorm.io/gorm"
	"html/template"
	"time"
)

type Major struct {
	Code        string `gorm:"primarykey"`
	RefreshedAt time.Time
}

type Course struct {
	MajorCode        string
	Major            Major  `gorm:"foreignKey:MajorCode"`
	CRN              string `gorm:"primarykey"`
	Code             string
	Catalogue        string
	Title            string
	TeachingMethod   string
	Instructor       string
	Capacity         int
	Enrolled         int
	Reservation      string
	MajorRestriction string
	Prerequisites    string
	ClassRestriction string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type Lecture struct {
	CourseCRN string `gorm:"primarykey"`
	Course    Course `gorm:"foreignKey:CourseCRN"`
	Building  string
	Day       string `gorm:"primarykey"`
	Time      string `gorm:"primarykey"`
	Room      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Post struct {
	Author  string        `json:"author"`
	Date    string        `json:"date"`
	Content template.HTML `json:"content"`
}