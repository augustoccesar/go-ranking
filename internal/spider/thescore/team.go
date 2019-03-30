package thescore

// Team is the struct that represents TheScore API response for Team (with
// stripped down fields for only what I need).
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"full_name"`
}
