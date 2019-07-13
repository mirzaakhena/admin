package model

// GetAllBasicRequest is
type GetAllBasicRequest struct {
	PageNumber int               `` //
	PageSize   int               `` //
	SortBy     string            `` //
	SortDir    string            `` //
	Filter     map[string]string `` //
}

// RegisterRequest is
type RegisterRequest struct {
	Name     string `json:"name"`     //
	Email    string `json:"email"`    //
	Password string `json:"password"` //
}

// LoginRequest is
type LoginRequest struct {
	Email    string `json:"email"`    //
	Password string `json:"password"` //
}

// ActivateRequest is
type ActivateRequest struct {
	RegisterToken string `json:"registerToken"` //
}

// ForgotPasswordInitRequest is
type ForgotPasswordInitRequest struct {
	Email string `json:"email"` //
}

// ForgotPasswordResetRequest is
type ForgotPasswordResetRequest struct {
	ResetPasswordToken string `json:"resetPasswordToken"` //
	NewPassword        string `json:"newPassword"`        //
}

// UpdateStatusRequest is
type UpdateStatusRequest struct {
	SpaceID string `json:"spaceId"` //
	UserID  string `json:"userId"`  //
	Status  string `json:"status"`  //
}

// GenerateInvitationAccountRequest is
type GenerateInvitationAccountRequest struct {
	SpaceID string `json:"spaceId"`
}

// UpdatePasswordRequest is
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"` //
	NewPassword string `json:"newPassword"` //
}

// CreateSpaceRequest is
type CreateSpaceRequest struct {
	Name        string `json:"name"`        //
	Description string `json:"description"` //
	SpacePlanID string `json:"spacePlanID"` //
}

// IsAdminRequest is
type IsAdminRequest struct {
	SpaceID string `json:"spaceID"` //
}

// GetBasicUserInfoRequest is
type GetBasicUserInfoRequest struct{}

// UpdateBasicUserInfoRequest is
type UpdateBasicUserInfoRequest struct {
	Name    string `json:"name"`    //
	Phone   string `json:"phone"`   //
	Address string `json:"address"` //
}

// RemoveAccountRequest is
type RemoveAccountRequest struct {
	SpaceID string `json:"spaceId"` //
	UserID  string `json:"userId"`  //
}

// RemoveWaitingAccountRequest is
type RemoveWaitingAccountRequest struct{}

// GetAllUserRolePermissionRequest is
type GetAllUserRolePermissionRequest struct{}

// CreateUserRolePermissionRequest is
type CreateUserRolePermissionRequest struct{}

// UpdateUserRolePermissionRequest is
type UpdateUserRolePermissionRequest struct{}

// DeleteUserRolePermissionRequest is
type DeleteUserRolePermissionRequest struct{}

// GetAllAccountUserRoleRequest is
type GetAllAccountUserRoleRequest struct{}

// UpdateAccountUserRoleRequest is
type UpdateAccountUserRoleRequest struct{}

// ApplyForPermissionRequest is
type ApplyForPermissionRequest struct{}

// GrantAppliedPermissionRequest is
type GrantAppliedPermissionRequest struct{}

// RefuseAppliedPermissionRequest is
type RefuseAppliedPermissionRequest struct{}

// IsAccessableRequest is
type IsAccessableRequest struct {
	MethodEndpoint string
	SpaceID        string
}
