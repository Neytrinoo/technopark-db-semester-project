package delivery

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"technopark-db-semester-project/domain"
)

type ServiceHandler struct {
	serviceRepo domain.ServiceRepo
}

func MakeServiceHandler(serviceRepo domain.ServiceRepo) ServiceHandler {
	return ServiceHandler{serviceRepo: serviceRepo}
}

// GET service/status
func (a *ServiceHandler) GetInfo(c echo.Context) error {
	result, _ := a.serviceRepo.GetInfo()
	return c.JSON(http.StatusOK, result)
}

// POST service/clear
func (a *ServiceHandler) Clear(c echo.Context) error {
	_ = a.serviceRepo.Clear()
	c.Response().Status = http.StatusOK

	return nil
}
