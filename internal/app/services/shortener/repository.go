package shortener

type Repository interface {
	GetLink(redirect *Redirect) error
	GetAllLinks(redirect *Redirect) error
	CreateLink(redirect *Redirect) error
	DeleteLink(redirect *Redirect) error
}
