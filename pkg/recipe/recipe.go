package recipe

// Difficulty is a type to represent difficulty levels
type Difficulty int

// Level of difficulty
const (
	Easy Difficulty = iota + 1
	Normal
	Hard
)

// Recipe is a recipe entry
type Recipe struct {
	ID            interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string      `json:"name" bson:"name"`
	PrepTime      string      `json:"prepTime" bson:"prepTime"`
	Difficulty    Difficulty  `json:"difficulty" bson:"difficulty"`
	Vegetarian    bool        `json:"vegetarian" bson:"vegetarian"`
	AverageRating float64     `json:"averageRating" bson:"averageRating"`
	RatingsCount  int64       `json:"ratingsCount" bson:"ratingsCount"`
}
