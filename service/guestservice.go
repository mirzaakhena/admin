package service

import (
	"fmt"
	"time"

	"github.com/mirzaakhena/admin/dao"
	"github.com/mirzaakhena/admin/model"
	myConfig "github.com/mirzaakhena/common/config"
	myEmail "github.com/mirzaakhena/common/email"
	log "github.com/mirzaakhena/common/logger"
	myPassword "github.com/mirzaakhena/common/password"
	myToken "github.com/mirzaakhena/common/token"
	"github.com/mirzaakhena/common/transaction"
	"github.com/mirzaakhena/common/utils"
)

// IGuestService is
type IGuestService interface {
	Register(sc model.ServiceContext, req model.RegisterRequest, bypassActivation bool) (*model.RegisterResponse, error)
	Activate(sc model.ServiceContext, req model.ActivateRequest) (*model.ActivateResponse, error)
	Login(sc model.ServiceContext, req model.LoginRequest) (*model.LoginResponse, error)
	ForgotPasswordInit(sc model.ServiceContext, req model.ForgotPasswordInitRequest) (*model.ForgotPasswordInitResponse, error)
	ForgotPasswordReset(sc model.ServiceContext, req model.ForgotPasswordResetRequest) (*model.ForgotPasswordResetResponse, error)
}

// GuestService is
type GuestService struct {
	Trx      transaction.ITransaction
	Email    myEmail.IEmail
	User     dao.IUserDao
	Token    myToken.IToken
	Password myPassword.IPassword
	Config   myConfig.IConfig
}

// Register is
func (o *GuestService) Register(sc model.ServiceContext, req model.RegisterRequest, bypassActivation bool) (*model.RegisterResponse, error) {

	logInfo := sc["logInfo"]

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByEmail(tx, req.Email)
	if bu != nil {

		log.GetLog().Error(logInfo, "user status is %v", bu.Status)

		if bu.Status == "ACTIVE" || bu.Status == "SUSPENDED" {
			o.Trx.RollbackTransaction(tx)
			return nil, utils.PrintError(model.ConstErrorExistingEmailAddress, "Email %s already exist", req.Email)
		}

		if bu.Status == "WAITING" {
			o.Trx.RollbackTransaction(tx)
			return nil, utils.PrintError(model.ConstErrorExistingEmailAddressButNotActiveYet, "Email %s is waiting for activation", req.Email)
		}
	}

	hashedPassword := o.Password.GenerateHashPassword(req.Password)

	userID := utils.GenID()

	data := map[string]string{
		"userId": userID,
	}

	registerToken := o.Token.GenerateToken("ACTIVATION", "APPS", "GUEST", data, 1)

	var us model.User
	us.ID = userID
	us.Name = req.Name
	us.Email = req.Email
	us.Password = hashedPassword

	if bypassActivation {
		us.Status = "ACTIVE"
		us.LoginToken = ""
	} else {
		us.Status = "WAITING"
		us.LoginToken = registerToken
	}

	o.User.Create(tx, &us)

	if !bypassActivation {

		host, ok := sc["host"]
		if !ok {
			log.GetLog().Panic(logInfo, "host not setup yet")
		}

		protocol := o.Config.GetString("server.protocol", "http")
		path := o.Config.GetString("link.activation", "/activate")
		activationLink := fmt.Sprintf("%s://%s%s?token=%s", protocol, host, path, registerToken)
		message := fmt.Sprintf("please click this %s to activate the account", activationLink)

		if err := o.Email.Send(o.Config.GetString("email.from", "app@mail.com"), req.Email, "Activation Account", message); err != nil {
			log.GetLog().Error(logInfo, "problem with sending email to %v", us.Email)
			o.Trx.RollbackTransaction(tx)
			return nil, utils.PrintError(model.ConstErrorWhenSendingEmail, "Server problem in sending email. %s", err.Error())
		}
	}

	o.Trx.CommitTransaction(tx)

	response := model.RegisterResponse{
		RegisterToken: registerToken,
	}

	return &response, nil
}

// Activate is
func (o *GuestService) Activate(sc model.ServiceContext, req model.ActivateRequest) (*model.ActivateResponse, error) {

	logInfo := sc["logInfo"]

	jwtObject := o.Token.ValidateToken("ACTIVATION", req.RegisterToken)

	if jwtObject == nil {
		// no token found
		log.GetLog().Error(logInfo, "token not found for %v", req.RegisterToken)
		return nil, utils.PrintError(model.ConstErrorTokenNotFound, "invalid token")
	}

	data, ok := jwtObject.ExtendData.(map[string]interface{})
	if !ok {
		log.GetLog().Error(logInfo, "token data is inval1d")
		return nil, utils.PrintError(model.ConstErrorTokenNotFound, "invalid token")
	}

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByID(tx, data["userId"].(string))

	if bu == nil {
		// user not exist
		log.GetLog().Error(logInfo, "invalid token")
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorUnExistingEmailAddress, "invalid token")
	}

	if bu.Status != "WAITING" {
		// user status is not WAITING
		log.GetLog().Error(logInfo, "user status is %v and not WAITING", bu.Status)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorInvalidExpectedStatus, "invalid token")
	}

	bu.Status = "ACTIVE"
	bu.LoginToken = ""
	o.User.Update(tx, bu)

	o.Trx.CommitTransaction(tx)

	response := model.ActivateResponse{}

	return &response, nil
}

// Login is
func (o *GuestService) Login(sc model.ServiceContext, req model.LoginRequest) (*model.LoginResponse, error) {

	logInfo := sc["logInfo"]

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByEmail(tx, req.Email)

	if bu == nil {
		// user not exist
		log.GetLog().Error(logInfo, "user with email %v is not exist", req.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorUnExistingEmailAddress, "Unexisting email or incorrect Password")
	}

	if bu.Status != "ACTIVE" {
		// user status is not ACTIVE
		log.GetLog().Error(logInfo, "user with email %v is not in status ACTIVE", req.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorInvalidExpectedStatus, "Unexisting email or incorrect Password")
	}

	log.GetLog().Info(logInfo, "Password %v xxx %v", req.Password, bu.Password)

	if !o.Password.IsValidPassword(req.Password, bu.Password) {
		// invalid password
		log.GetLog().Error(logInfo, "user with email %v entered wrong password", req.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorInvalidPassword, "Inexisting email or incorrect Password")
	}

	jwtObject := o.Token.ValidateToken("LOGIN", bu.LoginToken)

	oneHourAfterNow := time.Now().UTC().Add(time.Duration(1) * time.Hour)

	if jwtObject == nil || !jwtObject.VerifyExpiresAt(oneHourAfterNow.Unix(), true) {
		data := map[string]string{
			"userId": bu.ID,
		}
		bu.LoginToken = o.Token.GenerateToken("LOGIN", "APPS", "GUEST", data, 24)
		o.User.Update(tx, bu)
	}

	o.Trx.CommitTransaction(tx)

	response := model.LoginResponse{
		LoginToken: bu.LoginToken,
	}

	return &response, nil
}

// ForgotPasswordInit is
func (o *GuestService) ForgotPasswordInit(sc model.ServiceContext, req model.ForgotPasswordInitRequest) (*model.ForgotPasswordInitResponse, error) {

	logInfo := sc["logInfo"]

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByEmail(tx, req.Email)

	if bu == nil {
		// user not exist
		log.GetLog().Error(logInfo, "user with email %v is not exist", req.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorUnExistingEmailAddress, "User Is Not Exist")
	}

	if bu.Status != "ACTIVE" {
		// user status is not ACTIVE
		log.GetLog().Error(logInfo, "user with email %v is not in status ACTIVE", req.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorInvalidExpectedStatus, "User Is Not Active")
	}

	jwtObject := o.Token.ValidateToken("RESET_PASSWORD", bu.ResetPasswordToken)

	oneMinuteAfterNow := time.Now().UTC().Add(time.Duration(1) * time.Minute)

	if jwtObject == nil || !jwtObject.VerifyExpiresAt(oneMinuteAfterNow.Unix(), true) {
		data := map[string]string{
			"userId": bu.ID,
		}
		bu.ResetPasswordToken = o.Token.GenerateToken("LOGIN", "APPS", "GUEST", data, 1)
	}

	host, ok := sc["host"]
	if !ok {
		log.GetLog().Panic(logInfo, "host not setup yet")
	}
	protocol := o.Config.GetString("server.protocol", "http")
	path := o.Config.GetString("link.resetpassword", "/resetpassword")
	activationLink := fmt.Sprintf("%s://%s%s?token=%s", protocol, host, path, bu.ResetPasswordToken)
	message := fmt.Sprintf("please click this %s to activate the account", activationLink)

	if err := o.Email.Send(o.Config.GetString("email.from", "app@mail.com"), req.Email, "Reset Password", message); err != nil {
		log.GetLog().Error(logInfo, "problem with sending email to %v", req.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorWhenSendingEmail, "Server problem in sending email. %s", err.Error())
	}

	o.User.Update(tx, bu)

	o.Trx.CommitTransaction(tx)

	response := model.ForgotPasswordInitResponse{
		ResetPasswordToken: bu.ResetPasswordToken,
	}

	return &response, nil
}

// ForgotPasswordReset is
func (o *GuestService) ForgotPasswordReset(sc model.ServiceContext, req model.ForgotPasswordResetRequest) (*model.ForgotPasswordResetResponse, error) {

	logInfo := sc["logInfo"]

	jwtObject := o.Token.ValidateToken("RESET_PASSWORD", req.ResetPasswordToken)

	if jwtObject == nil {
		// no token found
		log.GetLog().Error(logInfo, "invalid token %v", req.ResetPasswordToken)
		return nil, utils.PrintError(model.ConstErrorTokenNotFound, "invalid token")
	}

	data, ok := jwtObject.ExtendData.(map[string]string)
	if !ok {
		log.GetLog().Error(logInfo, "token data is invalid")
		return nil, utils.PrintError(model.ConstErrorTokenNotFound, "invalid token")
	}

	tx := o.Trx.GetDB(true)

	bu := o.User.GetOneByID(tx, data["userId"])

	if bu == nil {
		// user not exist
		log.GetLog().Error(logInfo, "user not exist")
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorUnExistingEmailAddress, "invalid token")
	}

	if bu.Status != "ACTIVE" {
		// user status is not ACTIVE
		log.GetLog().Error(logInfo, "user with email %v is not in status ACTIVE", bu.Email)
		o.Trx.RollbackTransaction(tx)
		return nil, utils.PrintError(model.ConstErrorInvalidExpectedStatus, "invalid token")
	}

	bu.Password = req.NewPassword
	bu.ResetPasswordToken = ""

	o.Trx.CommitTransaction(tx)

	response := model.ForgotPasswordResetResponse{}

	return &response, nil
}
