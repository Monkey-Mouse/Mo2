package controller

import (
	"fmt"
	"net/http"

	"github.com/Monkey-Mouse/mo2/database"
	"github.com/Monkey-Mouse/mo2/mo2utils"
	"github.com/Monkey-Mouse/mo2/server/controller/badresponse"
	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/Monkey-Mouse/mo2/services/loghelper"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var notifyLogClient = loghelper.GetNotifyLogClient()

// GetComment get comments
// @Summary get comments
// @Description get json comments
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path string true "article id"
// @Param page query int false "int 0" 0
// @Param pagesize query int false "int 5" 5
// @Success 200 {object} []model.Comment
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/comment/{id} [get]
func (c *Controller) GetComment(ctx *gin.Context) {
	sid := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	page, pagesize, err := mo2utils.ParsePagination(ctx)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	cs := database.GetComments(id, page, pagesize)
	ctx.JSON(200, cs)
}

// PostComment godoc
// @Summary upsert comments
// @Description upsert json comments
// @Tags comments
// @Accept  json
// @Produce  json
// @Param comment body model.Comment true "comment"
// @Success 200 {object} model.Comment
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/comment [post]
func (c *Controller) PostComment(ctx *gin.Context) {
	var cmt model.Comment
	if ctx.ShouldBindJSON(&cmt) != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	u, _ := mo2utils.GetUserInfo(ctx)
	cmt.Author = u.ID
	database.UpsertComment(&cmt)
	b := model.Blog{}
	database.BlogCol.FindOne(ctx,
		bson.M{"_id": cmt.Article},
		options.FindOne().SetProjection(bson.D{{"author_id", 1}, {"title", 1}})).Decode(&b)
	notifyLogClient.LogInfo(
		loghelper.Log{
			Operator:             u.ID,
			Operation:            loghelper.COMMENT,
			OperationTarget:      cmt.Article,
			OperationTargetOwner: b.AuthorID,
			ExtraMessage:         fmt.Sprintf(`评论了你的文章<a href="/article/%s">%s</a>：%s`, b.ID.Hex(), b.Title, cmt.Content)})
	ctx.JSON(200, &cmt)
}

// PostSubComment post subcomments
// @Summary upsert subcomments
// @Description upsert json comments
// @Tags comments
// @Accept  json
// @Produce  json
// @Param id path string true "comment id"
// @Param comment body model.Subcomment true "subcomment"
// @Success 200 {object} model.Subcomment
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/comment/{id} [post]
func (c *Controller) PostSubComment(ctx *gin.Context) {
	sid := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	var cmt model.Subcomment
	if ctx.ShouldBindJSON(&cmt) != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	u, _ := mo2utils.GetUserInfo(ctx)
	cmt.Aurhor = u.ID
	database.UpsertSubComment(id, &cmt)
	ctx.JSON(200, &cmt)
}

// GetCommentNum godoc
// @Summary count comments
// @Description get article comment num
// @Tags comments
// @Produce  json
// @Param id path string true "article id"
// @Success 200 {object} map[string]int64
// @Failure 422 {object} badresponse.ResponseError
// @Router /api/commentcount/{id} [get]
func (c *Controller) GetCommentNum(ctx *gin.Context) {
	sid := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		badresponse.SetErrResponse(ctx, http.StatusUnprocessableEntity, "格式错误")
		return
	}
	num := database.GetCommentNum(id)
	ctx.JSON(200, gin.H{"count": num})
}
