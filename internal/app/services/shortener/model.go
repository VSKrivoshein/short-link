package shortener


type Redirect struct {
	Id          string `db:"id"`
	UserId      string `db:"user_id"`
	Link        string `db:"link" json:"link" binding:"required" validate:"url"`
	LinkHash    string `db:"link_hash"`
	AllUserLinks map[string]string
}
