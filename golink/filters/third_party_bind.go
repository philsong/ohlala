package filters

import (
	"github.com/QLeelulu/goku"
	"github.com/philsong/ohlala/golink/config"
	"github.com/philsong/ohlala/golink/models"
	"github.com/philsong/ohlala/golink/utils"
)

type ThirdPartyBindFilter struct{}

func (f *ThirdPartyBindFilter) OnActionExecuting(ctx *goku.HttpContext) (ar goku.ActionResulter, err error) {

	sessionIdBase, err := ctx.Request.Cookie(config.ThirdPartyCookieKey)
	if err != nil || len(sessionIdBase.Value) == 0 {
		ar = ctx.NotFound("no user binding context found.")
		return
	}
	ctx.Data["thirdPartySessionIdBase"] = sessionIdBase.Value

	profileSessionId := models.ThirdParty_GetThirdPartyProfileSessionId(sessionIdBase.Value)
	profile := models.ThirdParty_GetThirdPartyProfileFromSession(profileSessionId)

	if profile == nil {
		ar = ctx.NotFound("no user binding context found.")
		return
	}

	ctx.ViewData["profile"] = profile
	if len(profile.Email) > 0 {
		sensitiveInfoRemovedEmail := utils.GetSensitiveInfoRemovedEmail(profile.Email)
		ctx.ViewData["directCreateEmail"] = sensitiveInfoRemovedEmail
	}

	var profileShow struct {
		Avatar bool
		Link   bool
		Name   string
	}
	profileShow.Avatar = (len(profile.AvatarUrl) > 0)
	profileShow.Link = (len(profile.Link) > 0)
	profileShow.Name = profile.GetDisplayName()
	ctx.ViewData["profileShow"] = profileShow

	return
}

func (f *ThirdPartyBindFilter) OnActionExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	if _, ok := ctx.ViewData["bindValues"].(map[string]string); !ok {
		profile := ctx.ViewData["profile"].(*models.ThirdPartyUserProfile)
		m := make(map[string]string)
		m["name"] = profile.GetDisplayName()
		ctx.ViewData["bindValues"] = m
	}
	return nil, nil
}

func (f *ThirdPartyBindFilter) OnResultExecuting(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	return nil, nil
}

func (f *ThirdPartyBindFilter) OnResultExecuted(ctx *goku.HttpContext) (goku.ActionResulter, error) {
	return nil, nil
}

func NewThirdPartyBindFilter() *ThirdPartyBindFilter {
	return &ThirdPartyBindFilter{}
}
