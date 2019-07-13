package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirzaakhena/admin/model"
	"github.com/mirzaakhena/admin/service"
	log "github.com/mirzaakhena/common/logger"
	"github.com/mirzaakhena/common/utils"
)

type MyApi struct {
	GuestService service.IGuestService
}

func (m *MyApi) Register(c *gin.Context) {

	logInfo := log.Data{
		ClientIP: c.ClientIP(),
		Session:  utils.GenerateUniqueID(),
		Type:     "API",
	}

	log.GetLog().Info(logInfo, "Called")

	sc := model.ServiceContext{
		"logInfo": logInfo,
		"host":    c.Request.Host,
	}

	var req model.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprintf("Invalid Params. Please check data structure and type"), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Request %v", req)

	res, err := m.GuestService.Register(sc, req)

	if err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error(), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Response. %v", res)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": fmt.Sprintf("Record with ID created"), "data": res.RegisterToken})
}

func (m *MyApi) Activate(c *gin.Context) {

	logInfo := log.Data{
		ClientIP: c.ClientIP(),
		Session:  utils.GenerateUniqueID(),
		Type:     "API",
	}

	log.GetLog().Info(logInfo, "Called")

	sc := model.ServiceContext{
		"logInfo": logInfo,
	}

	var req model.ActivateRequest
	req.RegisterToken = c.DefaultQuery("token", "")

	log.GetLog().Info(logInfo, "Request %v", req)

	res, err := m.GuestService.Activate(sc, req)

	if err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error(), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Response. %v", res)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Success", "data": nil})
}

func (m *MyApi) Login(c *gin.Context) {

	logInfo := log.Data{
		ClientIP: c.ClientIP(),
		Session:  utils.GenerateUniqueID(),
		Type:     "API",
	}

	log.GetLog().Info(logInfo, "Called")

	sc := model.ServiceContext{
		"logInfo": logInfo,
	}

	var req model.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprintf("Invalid Params. Please check data structure and type"), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Request %v", req)

	res, err := m.GuestService.Login(sc, req)

	if err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error(), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Response. %v", res)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Success", "data": res.LoginToken})
}
func (m *MyApi) ForgotPasswordInit(c *gin.Context) {

	logInfo := log.Data{
		ClientIP: c.ClientIP(),
		Session:  utils.GenerateUniqueID(),
		Type:     "API",
	}

	log.GetLog().Info(logInfo, "Called")

	sc := model.ServiceContext{
		"logInfo": logInfo,
	}

	var req model.ForgotPasswordInitRequest
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprintf("Invalid Params. Please check data structure and type"), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Request %v", req)

	res, err := m.GuestService.ForgotPasswordInit(sc, req)

	if err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error(), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Response. %v", res)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Success", "data": nil})
}
func (m *MyApi) ForgotPasswordReset(c *gin.Context) {

	logInfo := log.Data{
		ClientIP: c.ClientIP(),
		Session:  utils.GenerateUniqueID(),
		Type:     "API",
	}

	log.GetLog().Info(logInfo, "Called")

	sc := model.ServiceContext{
		"logInfo": logInfo,
	}

	var req model.ForgotPasswordResetRequest
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprintf("Invalid Params. Please check data structure and type"), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Request %v", req)

	res, err := m.GuestService.ForgotPasswordReset(sc, req)

	if err != nil {
		log.GetLog().Error(logInfo, "%v", err.Error())
		c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error(), "data": nil})
		return
	}

	log.GetLog().Info(logInfo, "Response. %v", res)
	c.JSON(http.StatusCreated, map[string]interface{}{"message": "Success", "data": nil})
}
