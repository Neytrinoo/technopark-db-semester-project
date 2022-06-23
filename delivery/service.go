package delivery

import (
	"context"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"technopark-db-semester-project/domain"
)

type ServiceHandler struct {
	serviceRepo domain.ServiceRepo
}

func MakeServiceHandler(serviceRepo domain.ServiceRepo) ServiceHandler {
	return ServiceHandler{serviceRepo: serviceRepo}
}

// GET service/status
func (a *ServiceHandler) GetInfo(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)
	result, _ := a.serviceRepo.GetInfo(uctx)

	body, _ := json.Marshal(result)
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}

// POST service/clear
func (a *ServiceHandler) Clear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	uctx := ctx.UserValue("ctx").(context.Context)

	_ = a.serviceRepo.Clear(uctx)
	
	ctx.SetStatusCode(fasthttp.StatusOK)

	return
}
