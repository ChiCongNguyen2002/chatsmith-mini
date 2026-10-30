package controller

import (
	"net/http"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/auth"
	"chatsmith-mini/internal/domain/user"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/reqctx"
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	cookiePath      = "/"
	routeHome       = "/"
	queryState      = "state"
	queryCode       = "code"
	cookieExpireNow = -1
)

type Auth struct {
	service *auth.Service
	cfg     configs.AuthConfig
}

func NewAuth(service *auth.Service, cfg configs.AuthConfig) *Auth {
	return &Auth{service: service, cfg: cfg}
}

func (a *Auth) Login(ctx *gin.Context) {
	redirectURL, stateToken, err := a.service.StartLogin(providerParam(ctx))
	if err != nil {
		response.ErrFrom(ctx, err)
		return
	}

	a.setCookie(ctx, auth.CookieState, stateToken, int(a.cfg.StateTTL.Seconds()))
	ctx.Redirect(http.StatusFound, redirectURL)
}

func (a *Auth) Callback(ctx *gin.Context) {
	stateToken, err := ctx.Cookie(auth.CookieState)
	if err != nil {
		response.ErrFrom(ctx, failure.New(failure.CodeUnauthorized))
		return
	}

	sessionToken, err := a.service.FinishLogin(ctx.Request.Context(), providerParam(ctx), stateToken, ctx.Query(queryState), ctx.Query(queryCode))
	if err != nil {
		response.ErrFrom(ctx, err)
		return
	}

	a.setCookie(ctx, auth.CookieState, "", cookieExpireNow)
	a.setCookie(ctx, auth.CookieSession, sessionToken, int(a.cfg.SessionTTL.Seconds()))
	ctx.Redirect(http.StatusFound, routeHome)
}

func (a *Auth) Logout(ctx *gin.Context) {
	if token, err := ctx.Cookie(auth.CookieSession); err == nil {
		if err := a.service.Logout(ctx.Request.Context(), token); err != nil {
			response.ErrFrom(ctx, err)
			return
		}
	}

	a.setCookie(ctx, auth.CookieSession, "", cookieExpireNow)
	response.Success(ctx, nil)
}

func (a *Auth) Me(ctx *gin.Context) {
	account, err := a.service.Me(ctx.Request.Context(), reqctx.UserID(ctx))
	if err != nil {
		response.ErrFrom(ctx, err)
		return
	}

	response.Success(ctx, user.ToResponse(*account))
}

func (a *Auth) setCookie(ctx *gin.Context, name, value string, maxAge int) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(name, value, maxAge, cookiePath, "", a.cfg.CookieSecure, true)
}

func providerParam(ctx *gin.Context) constants.AuthProvider {
	return constants.AuthProvider(ctx.Param(constants.ParamProvider))
}
