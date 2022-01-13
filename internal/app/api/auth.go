package api

import (
	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/VSKrivoshein/short-link/internal/app/services/author"
	u "github.com/VSKrivoshein/short-link/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const ShortenerCookieName = "SHORTENER_JWT"

type SignInInput struct {
	Email    string `json:"email" binding:"required" example:"test@mail.ru"`
	Password string `json:"password" binding:"required" example:"qwerty"`
}

// @Summary Sign in
// @Tags auth
// @Description Sign in for existed user
// @Accept json
// @Produce json
// @Param credentials body api.SignInInput true "valid email and password of existed user"
// @Success 200 {string} string "success"
// @Failure 401 {object} ErrResponse "login or password is incorrect"
// @Failure 422 {object} ErrResponse "incorrect struct of request or validation failed"
// @Failure 500 {object} ErrResponse "internal server error"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input SignInInput

	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(e.New(
			err,
			e.ErrUnprocessableEntity,
			http.StatusUnprocessableEntity),
		)
		return
	}

	user := author.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.Services.Author.SingIn(&user); err != nil {
		_ = c.Error(err)
		return
	}

	maxEdgeSeconds := int(time.Until(user.TokenExpiration) / time.Second)

	c.SetCookie(
		ShortenerCookieName,
		user.TokenString,
		maxEdgeSeconds,
		"/",
		h.Config.Host,
		false,
		true,
	)
}

type SignUpInput struct {
	Email    string `json:"email" binding:"required" example:"test@mail.ru"`
	Password string `json:"password" binding:"required" example:"qwerty" minLength:"6"`
}

// @Summary Sign up
// @Tags auth
// @Description Registration of new user
// @Accept json
// @Produce json
// @Param credentials body api.SignUpInput true "valid email and password more than 6 chars"
// @Success 200 {string} string "success"
// @Failure 409 {object} ErrResponse "email already exist"
// @Failure 422 {object} ErrResponse "incorrect struct of request or validation failed"
// @Failure 500 {object} ErrResponse "internal server error"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input SignUpInput

	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(e.New(err, e.ErrUnprocessableEntity, http.StatusUnprocessableEntity))
		return
	}

	user := author.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.Services.Author.SingUp(&user); err != nil {
		_ = c.Error(err)
	}

	c.SetCookie(
		ShortenerCookieName,
		user.TokenString,
		int(time.Until(user.TokenExpiration)/time.Second),
		"/",
		h.Config.Host,
		false,
		true,
	)
}

// @Summary Sign out
// @Tags auth
// @Description Remove jwt token
// @Success 200 {string} string "success"
// @Router /auth/sign-out [get]
func (h *Handler) signOut(c *gin.Context) {
	c.SetCookie(
		ShortenerCookieName,
		"",
		-1,
		"/",
		h.Config.Host,
		false,
		true,
	)
}

type DeleteUserInput struct {
	Email    string `json:"email" binding:"required" example:"test@mail.ru"`
	Password string `json:"password" binding:"required" example:"qwerty"`
}

// @Summary Delete user
// @Tags auth
// @Description Delete existed user with all user links
// @Param credentials body api.DeleteUserInput true "valid email and password of existed user"
// @Success 200 {string} string "success"
// @Failure 401 {object} ErrResponse "login or password is incorrect"
// @Failure 422 {object} ErrResponse "incorrect struct of request or validation failed"
// @Failure 500 {object} ErrResponse "internal server error"
// @Router /auth/delete-user [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	var input DeleteUserInput

	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(e.New(err, e.ErrUnprocessableEntity, http.StatusUnprocessableEntity))
		return
	}

	user := author.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.Services.Author.SingIn(&user); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.Services.Author.DeleteUser(&user); err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie(
		ShortenerCookieName,
		"",
		-1,
		"/",
		h.Config.Host,
		false,
		true,
	)

	c.Status(http.StatusOK)
}

func (h *Handler) checkAuthAndRefreshMiddleware(c *gin.Context) {
	jwtString, err := c.Cookie(ShortenerCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			_ = c.Error(e.New(err, e.ErrUserUnauthorized, http.StatusUnauthorized))
			return
		}
		_ = c.Error(e.New(err, e.ErrToken, http.StatusInternalServerError))
		return
	}

	user := author.User{
		TokenString: jwtString,
	}

	if err := h.Services.Author.CheckAuthAndRefresh(&user); err != nil {
		_ = c.Error(err)
		return
	}

	u.SetUserID(c, user.UserID)

	c.SetCookie(
		ShortenerCookieName,
		user.TokenString,
		int(time.Until(user.TokenExpiration)/time.Second),
		"/",
		h.Config.Host,
		false,
		true,
	)
}
