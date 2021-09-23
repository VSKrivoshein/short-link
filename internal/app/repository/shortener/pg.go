package shortener_repo

import (
	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/VSKrivoshein/short-link/internal/app/services/shortener"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) shortener.Repository {
	return &repository{db: db}
}

func (r *repository) GetLink(redirect *shortener.Redirect) error {
	err := r.db.Get(redirect, `
		SELECT id, user_id, link, link_hash
		FROM links 
		WHERE link_hash=$1;`,
		redirect.LinkHash)
	if err != nil {
		return e.New(err, e.ErrGetLink, http.StatusNotFound)
	}

	return nil
}

func (r *repository) GetAllLinks(redirect *shortener.Redirect) error {
	var data []struct {
		Link     string `db:"link"`
		LinkHash string `db:"link_hash"`
	}

	if err := r.db.Select(&data, `
		SELECT link, link_hash 
		FROM links 
		    JOIN users u 
		        ON u.id = links.user_id 
		WHERE user_id=$1;`, redirect.UserId); err != nil {
		return e.New(err, e.ErrGetAllLinks, http.StatusInternalServerError)
	}

	links := make(map[string]string)
	for _, v := range data {
		links[v.Link] = v.LinkHash
	}

	redirect.AllUserLinks = links

	return nil
}

func (r *repository) CreateLink(redirect *shortener.Redirect) error {
	if err := r.db.QueryRowx(`
		INSERT INTO links (id, user_id, link, link_hash) 
		VALUES (gen_random_uuid(), $1, $2, $3)
		RETURNING link, link_hash;`,
		redirect.UserId,
		redirect.Link,
		redirect.LinkHash,
	).StructScan(redirect); err != nil {
		return e.New(err, e.ErrCreatingLink, http.StatusConflict)
	}

	return nil
}

func (r *repository) DeleteLink(redirect *shortener.Redirect) error {
	res, err := r.db.NamedExec(`
		DELETE FROM links 
		USING users 
		WHERE  links.user_id = users.id 
		  AND links.link = :link;`,
		redirect)
	if err != nil {
		return e.New(err, e.ErrDeleteLink, http.StatusInternalServerError)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return e.New(err, e.ErrDeleteLink, http.StatusInternalServerError)
	}
	if count == 0 {
		return e.New(err, e.ErrDeleteLinkNotFound, http.StatusUnprocessableEntity)
	}

	return nil
}
