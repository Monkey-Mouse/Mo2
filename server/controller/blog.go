package controller

import (
	"mo2/database"
	"mo2/mo2utils"
	"mo2/server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PublishBlog godoc
// @Summary Publish Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Param account body model.Blog true "Add blog"
// @Success 200 {object} model.Blog
// @Router /api/blogs/publish [post]
func (c *Controller) PublishBlog(ctx *gin.Context) {
	b := model.Blog{}
	if err := ctx.ShouldBindJSON(&b); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, SetResponseReason("内容含非法字符，请检查"))
		return
	}
	// set blog's author due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
		return
	}
	b.AuthorID = info.ID
	database.AddBlog(&b)
	ctx.JSON(http.StatusOK, &b)
}

// FindBlogsByUser godoc
// @Summary find Blog
// @Description add by json
// @Tags blogs
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Blog
// @Router /api/blogs/find/byUser [get]
func (c *Controller) FindBlogsByUser(ctx *gin.Context) {
	// get user info due to cookie information
	info, ext := mo2utils.GetUserInfo(ctx)
	if !ext {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, SetResponseReason("权限不足，请先登录"))
		return
	}

	blogs := database.FindBlogs(info)
	ctx.JSON(http.StatusOK, blogs)
}
