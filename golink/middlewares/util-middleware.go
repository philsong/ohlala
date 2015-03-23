package middlewares

import (
	"github.com/QLeelulu/goku"
	"github.com/philsong/ohlala/golink/models"
)

// 一些基本的处理
// 例如检查用户是否登陆，如果登陆则获取登陆用户信息，并添加 ctx.Data 中
type UtilMiddleware struct {
}

func (tmd *UtilMiddleware) OnBeginRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	return nil, nil
}

func (tmd *UtilMiddleware) OnBeginMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	user := getUser(ctx)
	getTopNavTopics(ctx, user)
	return checkUserStatus(user, ctx)
}
func (tmd *UtilMiddleware) OnEndMvcHandle(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	return nil, nil
}

func (tmd *UtilMiddleware) OnEndRequest(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	return nil, nil
}

func getUser(ctx *goku.HttpContext) *models.User {
	c, err := ctx.Request.Cookie("_glut")
	if err == nil {
		user, _ := models.User_GetByTicket(c.Value)
		if user != nil {
			ctx.Data["user"] = user
			// 暂时先设置到ViewData里面吧，应该需要一个更好的办法？
			ctx.ViewData["user"] = user
			return user
		}
	}
	return nil
}

// 检查用户的状态（禁言、封号等）
func checkUserStatus(user *models.User, ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {
	if user != nil && user.IsBaned() {
		if ctx.Request.Method == "POST" {
			if ctx.IsAjax() {
				ar = ctx.Json(map[string]interface{}{
					"success":   false,
					"needLogin": false,
					"errors":    "你已经被禁言。",
				})
			} else {
				ctx.ViewData["errorMsg"] = "你已经被禁言"
				ar = ctx.Render("error", nil)
			}
		}
	}
	return
}

// 顶部导航栏的话题列表。
// 如果用户已经登陆，则获取用户关注的话题，
// 如果未登陆则获取全站的最流行话题列表.
func getTopNavTopics(ctx *goku.HttpContext, user *models.User) {
	var topics []models.Topic
	topicLen := 35
	if user == nil {
		topics, _ = models.Topic_GetTops(1, topicLen)
	} else {
		tuser, _ := models.User_GetFollowTopics(user.Id, 1, topicLen, "link_count desc")
		if len(tuser) < topicLen {
			// 不够 topicLen 条，则合并
			tall, _ := models.Topic_GetTops(1, topicLen-len(tuser))
			topics = make([]models.Topic, 0, len(tall))
			tmp := map[string]bool{}
			for _, v := range tuser {
				tmp[v.Name] = true
			}
			topics = append(topics, tuser...)
			for _, v := range tall {
				if _, ok := tmp[v.Name]; !ok {
					topics = append(topics, v)
				}
			}
		} else {
			topics = tuser
		}
	}
	ctx.ViewData["TopNavTopics"] = topics
}
