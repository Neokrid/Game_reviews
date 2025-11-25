package model

import "time"

type RatingHistory struct {
	Date      time.Time `json:"date" db:"review_date"`
	AvgRating float32   `json:"avgRating" db:"avg_rating"`
}
