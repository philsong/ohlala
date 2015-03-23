package controllers

import (
	"github.com/QLeelulu/goku"
	"github.com/philsong/ohlala/golink/filters"
	// "github.com/philsong/ohlala/golink/models"
)

var _ = goku.Controller("host").
	/**
	 * 查看一个host下的链接
	 */
	Get("show", func(ctx *goku.HttpContext) goku.ActionResulter {

	return ctx.View(nil)
}).Filters(filters.NewRequireLoginFilter())
