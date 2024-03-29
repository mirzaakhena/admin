package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/mirzaakhena/admin/model"
	log "github.com/mirzaakhena/common/logger"
	"github.com/mirzaakhena/common/utils"
)

// IAdminService is
type IAdminService interface {
	IUserService

	GetOneSpace(sc model.ServiceContext) *model.Space
	GetAllUserSpace(sc model.ServiceContext, req model.GetAllBasicRequest) *model.GetAllSpaceResponse
	CreateSpace(sc model.ServiceContext, req model.CreateSpaceRequest) (*model.CreateSpaceResponse, error)

	// IsAdmin(sc model.ServiceContext, req model.IsAdminRequest) bool

	// GenerateInvitationAccount(sc model.ServiceContext, req model.GenerateInvitationAccountRequest) (*model.GenerateInvitationAccountResponse, error)
	// UpdateAccountStatus(sc model.ServiceContext, req model.UpdateStatusRequest) (*model.UpdateStatusResponse, error)
	// RemoveAccount(sc model.ServiceContext, req model.RemoveAccountRequest) (*model.RemoveAccountResponse, error)

	// RemoveWaitingAccount(sc model.ServiceContext, req model.RemoveWaitingAccountRequest) (*model.RemoveWaitingAccountResponse, error)

	// GetAllUserRolePermission(sc model.ServiceContext, req model.GetAllBasicRequest) ([]model.GetAllUserRolePermissionResponse, uint)
	// CreateUserRolePermission(sc model.ServiceContext, req model.CreateUserRolePermissionRequest) (*model.CreateUserRolePermissionResponse, error)
	// UpdateUserRolePermission(sc model.ServiceContext, req model.UpdateUserRolePermissionRequest) (*model.UpdateUserRolePermissionResponse, error)
	// DeleteUserRolePermission(sc model.ServiceContext, req model.DeleteUserRolePermissionRequest) (*model.DeleteUserRolePermissionResponse, error)

	// GetAllAccountUserRole(sc model.ServiceContext, req model.GetAllBasicRequest) ([]model.GetAllAccountUserRoleResponse, uint)
	// UpdateAccountUserRole(sc model.ServiceContext, req model.UpdateAccountUserRoleRequest) (*model.UpdateAccountUserRoleResponse, error)
}

// AdminService is
type AdminService struct {
	UserService
}

// GetOneSpace is
func (o *AdminService) GetOneSpace(sc model.ServiceContext) *model.Space {
	spaceID := sc["spaceId"].(string)
	return o.Space.GetOne(o.Trx.GetDB(false), spaceID)
}

// GetAllUserSpace is
func (o *AdminService) GetAllUserSpace(sc model.ServiceContext, req model.GetAllBasicRequest) *model.GetAllSpaceResponse {
	var ss []model.Space
	us, count := o.UserSpace.GetAll(o.Trx.GetDB(false), req)
	for _, s := range us {
		ss = append(ss, *s.Space)
	}
	result := model.GetAllSpaceResponse{
		TotalCount: count,
		Spaces:     ss,
	}
	return &result
}

// CreateSpace is
func (o *AdminService) CreateSpace(sc model.ServiceContext, req model.CreateSpaceRequest) (*model.CreateSpaceResponse, error) {

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if name == "" {
		return nil, fmt.Errorf("space name must not empty")
	}

	userID, logInfo := o.ExtractServiceContext(sc)

	tx := o.Trx.GetDB(true)

	if o.Space.IsExistName(tx, name, userID) {
		log.GetLog().Error(logInfo, "space with name %s is exist", name)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorUnExistingEmailAddress, "space with name %s is exist. ", name)
	}

	var ws model.Space
	{
		ws.ID = utils.GenID()
		ws.Name = name
		ws.Description = description
		ws.MaxUser = 5
		ws.TotalCurrentUser = 1
		ws.Expired = time.Now().Add(time.Hour * 24 * 100000)
		o.Space.Create(tx, ws)
	}

	var wsa model.UserSpace
	{
		wsa.ID = utils.GenID()
		wsa.UserID = userID
		wsa.SpaceID = ws.ID
		wsa.Type = "ADMIN"
		wsa.Status = "ACTIVE"
		o.UserSpace.Create(tx, wsa)
	}

	o.Trx.CommitTransaction(tx)

	response := model.CreateSpaceResponse{}

	return &response, nil
}

// IsAdmin is
func (o *AdminService) IsAdmin(sc model.ServiceContext, req model.IsAdminRequest) bool {

	userID, _ := o.ExtractServiceContext(sc)

	tx := o.Trx.GetDB(false)
	us := o.UserSpace.GetOne(tx, req.SpaceID, userID)

	return us.ID != "" && us.Type == "ADMIN" && us.Status == "ACTIVE"
}

// GenerateInvitationAccount is
func (o *AdminService) GenerateInvitationAccount(sc model.ServiceContext, req model.GenerateInvitationAccountRequest) (*model.GenerateInvitationAccountResponse, error) {
	data := map[string]string{
		"SpaceId": req.SpaceID,
	}
	token := o.Token.GenerateToken("INVITATION", "APPS", "NEWUSER", data, 24)

	response := model.GenerateInvitationAccountResponse{
		SpaceInvitationToken: token,
	}

	return &response, nil
}

// UpdateAccountStatus is
func (o *AdminService) UpdateAccountStatus(sc model.ServiceContext, req model.UpdateStatusRequest) (*model.UpdateStatusResponse, error) {

	// userID, logInfo := o.ExtractServiceContext(sc)

	tx := o.Trx.GetDB(true)

	wsa := o.UserSpace.GetOne(tx, req.SpaceID, req.UserID)
	wsa.Status = req.Status
	o.UserSpace.Update(tx, wsa.ID, wsa)

	o.Trx.CommitTransaction(tx)

	response := model.UpdateStatusResponse{}

	return &response, nil
}

// RemoveAccount is
func (o *AdminService) RemoveAccount(sc model.ServiceContext, req model.RemoveAccountRequest) (*model.RemoveAccountResponse, error) {
	tx := o.Trx.GetDB(true)

	wsa := o.UserSpace.GetOne(tx, req.SpaceID, req.UserID)
	o.UserSpace.Delete(tx, wsa.ID)

	o.Trx.CommitTransaction(tx)

	response := model.RemoveAccountResponse{}

	return &response, nil
}

// GetAllAppliedPermission is
func (o *AdminService) GetAllAppliedPermission(sc model.ServiceContext, req model.GetAllBasicRequest) ([]model.GetAllAppliedPermissionResponse, uint, error) {

	return nil, 0, nil
}

// GrantAppliedPermission is
func (o *AdminService) GrantAppliedPermission(sc model.ServiceContext, req model.GrantAppliedPermissionRequest) (*model.GrantAppliedPermissionResponse, error) {
	return nil, nil
}

// RefuseAppliedPermission is
func (o *AdminService) RefuseAppliedPermission(sc model.ServiceContext, req model.RefuseAppliedPermissionRequest) (*model.RefuseAppliedPermissionResponse, error) {
	return nil, nil
}

// GetAllUserRolePermission is
func (o *AdminService) GetAllUserRolePermission(sc model.ServiceContext, req model.GetAllUserRolePermissionRequest) (*model.GetAllUserRolePermissionResponse, uint, error) {
	return nil, 0, nil
}

// CreateUserRolePermission is
func (o *AdminService) CreateUserRolePermission(sc model.ServiceContext, req model.CreateUserRolePermissionRequest) (*model.CreateUserRolePermissionResponse, error) {
	return nil, nil
}

// UpdateUserRolePermission is
func (o *AdminService) UpdateUserRolePermission(sc model.ServiceContext, req model.UpdateUserRolePermissionRequest) (*model.UpdateUserRolePermissionResponse, error) {
	return nil, nil
}

// DeleteUserRolePermission is
func (o *AdminService) DeleteUserRolePermission(sc model.ServiceContext, req model.DeleteUserRolePermissionRequest) (*model.DeleteUserRolePermissionResponse, error) {
	return nil, nil
}

// GetAllAccountUserRole is
func (o *AdminService) GetAllAccountUserRole(sc model.ServiceContext, req model.GetAllAccountUserRoleRequest) (*model.GetAllAccountUserRoleResponse, uint, error) {
	return nil, 0, nil
}

// UpdateAccountUserRole is
func (o *AdminService) UpdateAccountUserRole(sc model.ServiceContext, req model.UpdateAccountUserRoleRequest) (*model.UpdateAccountUserRoleResponse, error) {
	return nil, nil
}
