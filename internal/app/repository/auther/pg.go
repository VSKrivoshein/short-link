package auther_repo

import (
	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/VSKrivoshein/short-link/internal/app/services/auther"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) auther.Repository {
	return &repository{db: db}
}


func (r *repository) CreateUser(user *auther.User) error {
	row, err := r.db.NamedQuery(`
		INSERT INTO users 
		VALUES (gen_random_uuid(), :email, :password_hash) 
		RETURNING id, email, password_hash;`,
		user)
	if err != nil {
		return e.New(err, e.ErrCreatingUser, http.StatusConflict)
	}

	if row.Next() {
		if err := row.StructScan(user); err != nil {
			return e.New(err, e.ErrCreatingUser, http.StatusInternalServerError)
		}
	}

	return nil
}

func (r *repository) GetUser(user *auther.User) error {
	err := r.db.Get(user, `
		SELECT id, email, password_hash 
		FROM users 
		WHERE email=$1;`,
		user.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return e.New(err, e.ErrGetUser, http.StatusUnauthorized)
		}
		return e.New(err, e.ErrUnexpected, http.StatusInternalServerError)
	}

	return nil
}

func (r *repository) DeleteUser(user *auther.User) error {
	res, err := r.db.Exec(`
		DELETE FROM users
		WHERE id=$1`,
		user.UserId)
	if err != nil {
		return e.New(err, e.ErrDeletingUser, http.StatusInternalServerError)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return e.New(err, e.ErrDeletingUser, http.StatusInternalServerError)
	}

	if count == 0 {
		return e.New(err, e.ErrDeletingUserNotFound, http.StatusInternalServerError)
	}

	return nil
}
