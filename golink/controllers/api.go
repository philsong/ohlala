package controllers

import (
	"fmt"
	"github.com/QLeelulu/goku"
	"github.com/philsong/ohlala/golink/filters"
	"github.com/philsong/ohlala/golink/forms"
	"github.com/philsong/ohlala/golink/models"
)

var _ = goku.Controller("api").
	//
	Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

	return ctx.View(nil)
}).
	/**
	 * 获取一个链接的信息
	 */
	Get("link_info", func(ctx *goku.HttpContext) goku.ActionResulter {
	return ctx.View(nil)
}).

	/**
	 * 提交一个链接并保存到数据库
	 */
	Post("link_submit", func(ctx *goku.HttpContext) goku.ActionResulter {

	f := forms.CreateLinkSubmitForm()
	f.FillByRequest(ctx.Request)

	var resubmit bool
	success, linkId, errorMsgs, _ := models.Link_SaveForm(f, (ctx.Data["user"].(*models.User)).Id, resubmit)

	if success {
		//go addLinkForSearch(0, m, linkId, user.Name) //contextType:0: url, 1:文本   TODO:
		return ctx.Redirect(fmt.Sprintf("/link/%d", linkId))
	} else {
		ctx.ViewData["Errors"] = errorMsgs
		ctx.ViewData["Values"] = f.Values()
	}
	r := map[string]interface{}{
		"success": success,
		"errors":  errorMsgs,
	}
	return ctx.Json(r)

}).Filters(filters.NewRequireLoginFilter()).

	/**
	 * 添加评论
	 */
	Post("link_comment", func(ctx *goku.HttpContext) goku.ActionResulter {

	return ctx.View(nil)

}).Filters(filters.NewRequireLoginFilter())
