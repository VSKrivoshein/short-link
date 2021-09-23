package shortener

type Service interface {
	GetLink(redirect *Redirect) error
	GetAllLinks(redirect *Redirect) error
	CreateLink(redirect *Redirect) error
	DeleteLink(redirect *Redirect) error
}
