package service

import (
	"github.com/mirzaakhena/admin/dao"
	model "github.com/mirzaakhena/admin/model"
	log "github.com/mirzaakhena/common/logger"
	"github.com/mirzaakhena/common/utils"
)

// IUserService is
type IUserService interface {
	IGuestService
	IsAccessable(sc model.ServiceContext, req model.IsAccessableRequest) bool
	GetBasicUserInfo(sc model.ServiceContext, req model.GetBasicUserInfoRequest) *model.GetBasicUserInfoResponse
	UpdateBasicUserInfo(sc model.ServiceContext, req model.UpdateBasicUserInfoRequest) (*model.UpdateBasicUserInfoResponse, error)
	UpdatePassword(sc model.ServiceContext, req model.UpdatePasswordRequest) (*model.UpdatePasswordResponse, error)
	GetAllPermission(sc model.ServiceContext, req model.GetAllBasicRequest) (*model.GetAllPermissionResponse, uint64)
	CreateAdminUserIfNotExist(sc model.ServiceContext)
	ExtractServiceContext(sc model.ServiceContext) (string, interface{})
}

// UserService is
type UserService struct {
	GuestService
	Space          dao.ISpaceDao
	UserSpace      dao.IUserSpaceDao
	UserPermission dao.IUserPermissionDao
}

// CreateAdminUserIfNotExist is
func (o *UserService) CreateAdminUserIfNotExist(sc model.ServiceContext) {
	tx := o.Trx.GetDB(true)
	_, count := o.User.GetAll(tx, model.GetAllBasicRequest{})
	if count > 0 {
		return
	}

	userID := utils.GenID()

	var us model.User
	us.ID = userID
	us.Name = "Admin"
	us.Email = "admin@mail.com"
	us.Password = "12345"
	us.Status = "ACTIVE"
	us.LoginToken = ""

	o.User.Create(tx, &us)

	o.Trx.CommitTransaction(tx)
}

// IsAccessable is
func (o *UserService) IsAccessable(sc model.ServiceContext, req model.IsAccessableRequest) bool {
	userID, _ := o.ExtractServiceContext(sc)

	up := o.UserPermission.GetUserPermission(o.Trx.GetDB(false), userID, req.SpaceID, req.MethodEndpoint)

	return up != nil
}

// GetBasicUserInfo is
func (o *UserService) GetBasicUserInfo(sc model.ServiceContext, req model.GetBasicUserInfoRequest) *model.GetBasicUserInfoResponse {

	userID, _ := o.ExtractServiceContext(sc)

	bu := o.User.GetOneByID(o.Trx.GetDB(false), userID)

	response := model.GetBasicUserInfoResponse{
		User: bu,
	}

	return &response
}

// UpdateBasicUserInfo is
func (o *UserService) UpdateBasicUserInfo(sc model.ServiceContext, req model.UpdateBasicUserInfoRequest) (*model.UpdateBasicUserInfoResponse, error) {
	userID, _ := o.ExtractServiceContext(sc)

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByID(tx, userID)

	bu.Address = req.Address
	bu.Name = req.Name
	bu.Phone = req.Phone

	o.User.Update(tx, bu)

	response := model.UpdateBasicUserInfoResponse{}

	o.Trx.CommitTransaction(tx)

	return &response, nil
}

// UpdatePassword is
func (o *UserService) UpdatePassword(sc model.ServiceContext, req model.UpdatePasswordRequest) (*model.UpdatePasswordResponse, error) {

	userID, logInfo := o.ExtractServiceContext(sc)

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByID(tx, userID)

	if !o.Password.IsValidPassword(req.OldPassword, bu.Password) {
		log.GetLog().Error(logInfo, "Old Password does not match")
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorUnExistingEmailAddress, "Old Password does not match ")
	}

	bu.Password = o.Password.GenerateHashPassword(req.NewPassword)

	o.User.Update(tx, bu)

	o.Trx.CommitTransaction(tx)

	return nil, nil
}

// GetAllPermission is
func (o *UserService) GetAllPermission(sc model.ServiceContext, req model.GetAllBasicRequest) (*model.GetAllPermissionResponse, uint64) {

	// userID, logInfo := o.ExtractServiceContext(sc)

	return nil, 0
}

// ExtractServiceContext is
func (o *UserService) ExtractServiceContext(sc model.ServiceContext) (string, interface{}) {
	logInfo := sc["logInfo"]
	userIDInterface, ok := sc["userId"]
	if !ok {
		log.GetLog().Panic(logInfo, "user.id not setup yet")
	}

	userID := userIDInterface.(string)
	return userID, logInfo
}
