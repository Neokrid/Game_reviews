package model

type Leaderboard struct {
	Position  int     `json:"position" db:"position"`
	GameName  string  `json:"gameName" db:"title"`
	AvgRating float32 `json:"avgRating" db:"average_rating"`
}
