package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/herodot"
	"github.com/ory/hydra/policy"
)

func newPolicyHandler(c *config.Config, router *httprouter.Router) *policy.Handler {
	ctx := c.Context()
	h := &policy.Handler{
		H:       &herodot.JSON{},
		W:       ctx.Warden,
		Manager: ctx.LadonManager,
	}
	h.SetRoutes(router)
	return h
}
