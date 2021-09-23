package api

import (
	"github.com/VSKrivoshein/short-link/internal/app/e"
	"github.com/VSKrivoshein/short-link/internal/app/services/shortener"
	u "github.com/VSKrivoshein/short-link/internal/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Redirect
// @Tags redirect
// @Description url with hash will "redirect" user to link from service
// @Param hash path string true "hash that was generated during creating link"
// @Success 307 {string} string "redirect user to link"
// @Failure 404 {object} ErrResponse "link was not found"
// @Router /{hash} [get]
func (h *Handler) redirect(c *gin.Context) {

	redirect := shortener.Redirect{
		LinkHash: c.Param("hash"),
	}

	if err := h.Services.Shortener.GetLink(&redirect); err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, redirect.Link)
}

type CreateLinkInput struct {
	Link string `json:"link" binding:"required" example:"https://ya.ru/"`
}

type CreateLinkResp struct {
	Link        string `json:"link" example:"https://ya.ru/"`
	RedirectURL string `json:"redirect_url" example:"localhost:8080/a425tq"`
}

// @Summary Creating redirect
// @Tags links
// @Description Creating new redirect from valid url for authenticated user
// @Accept json
// @Produce json
// @Param credentials body api.CreateLinkInput true "valid url"
// @Success 200 {object} CreateLinkResp "link and redirect url"
// @Failure 401 {object} ErrResponse "user should be authenticated"
// @Failure 409 {object} ErrResponse "link is already exist"
// @Failure 422 {object} ErrResponse "incorrect struct of request or validation failed"
// @Failure 500 {object} ErrResponse "internal server error"
// @Router /links/create [post]
func (h *Handler) create(c *gin.Context) {
	var redirect shortener.Redirect

	if err := c.BindJSON(&redirect); err != nil {
		c.Error(e.New(
			err,
			e.ErrUnprocessableEntity,
			http.StatusUnprocessableEntity),
		)
		return
	}

	userId := u.GetUserID(c)
	if userId == "" {
		c.Error(e.New(e.ErrUserIdWasNotFound, e.ErrUserUnauthorized, http.StatusUnauthorized))
		return
	}

	redirect.UserId = userId

	if err := h.Services.Shortener.CreateLink(&redirect); err != nil {
		c.Error(err)
		return
	}

	RedirectLink := h.Config.Origin + redirect.LinkHash

	c.JSON(http.StatusOK, CreateLinkResp{
		Link:        redirect.Link,
		RedirectURL: RedirectLink,
	})
}

type GetAllLinksResp struct {
	AllUserLinks map[string]string `json:"all_user_links" example:"short url:original url"`
}

// @Summary Get all links that belong to user wit hash
// @Tags links
// @Description Get all links that belong to user wit hash
// @Accept json
// @Produce json
// @Success 200 {object} GetAllLinksResp "pair with link and redirect link"
// @Failure 401 {object} ErrResponse "user should be authenticated"
// @Failure 422 {object} ErrResponse "incorrect struct of request or validation failed"
// @Failure 500 {object} ErrResponse "internal server error"
// @Router /links/get-all [get]
func (h *Handler) getAllLinks(c *gin.Context) {
	var redirect shortener.Redirect

	userId := u.GetUserID(c)
	if userId == "" {
		c.Error(e.New(e.ErrUserIdWasNotFound, e.ErrUserUnauthorized, http.StatusUnauthorized))
		return
	}

	redirect.UserId = userId

	if err := h.Services.Shortener.GetAllLinks(&redirect); err != nil {
		c.Error(err)
		return
	}

	for k, v := range redirect.AllUserLinks {
		redirect.AllUserLinks[k] = h.Config.Origin + v
	}

	c.JSON(http.StatusOK, GetAllLinksResp{
		AllUserLinks: redirect.AllUserLinks,
	})
}

type DeleteLinkInput struct {
	Link string `json:"link" binding:"required"`
}

// @Summary Delete link
// @Tags links
// @Description Delete link with hash
// @Accept json
// @Success 200 {string} string "success"
// @Failure 401 {object} ErrResponse "user should be authenticated"
// @Failure 422 {object} ErrResponse "link was not found"
// @Failure 500 {object} ErrResponse "internal server error"
// @Router /links/delete [delete]
func (h *Handler) deleteLink(c *gin.Context) {
	var redirect shortener.Redirect

	if err := c.BindJSON(&redirect); err != nil {
		c.Error(e.New(err, e.ErrUnprocessableEntity, http.StatusUnprocessableEntity))
		return
	}

	userId := u.GetUserID(c)
	if userId == "" {
		c.Error(e.New(e.ErrUserIdWasNotFound, e.ErrUserUnauthorized, http.StatusUnauthorized))
		return
	}

	redirect.UserId = userId

	if err := h.Services.Shortener.DeleteLink(&redirect); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
