package dao

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mirzaakhena/admin/model"
)

// ISpaceDao is
type ISpaceDao interface {
	IsExistName(dc model.DaoContext, spaceName, userID string) bool
	Create(dc model.DaoContext, bu model.Space) error
	GetOne(dc model.DaoContext, ID string) *model.Space
	GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.Space, string)
	Delete(dc model.DaoContext, ID string) error
	Update(dc model.DaoContext, ID string, obj *model.Space) error
}

// SpaceDao is
type SpaceDao struct {
}

// IsExistName is
func (s *SpaceDao) IsExistName(dc model.DaoContext, spaceName, userID string) bool {
	var userSpaces []model.UserSpace
	dc.(*gorm.DB).Preload("Space").Find(&userSpaces, "user_id = ?", userID)
	for _, us := range userSpaces {
		if strings.Compare(strings.ToLower(us.Space.Name), strings.ToLower(spaceName)) == 0 {
			return true
		}
	}
	return false
}

// Create is
func (s *SpaceDao) Create(dc model.DaoContext, bu model.Space) error {
	return dc.(*gorm.DB).Create(&bu).Error
}

// GetOne is
func (s *SpaceDao) GetOne(dc model.DaoContext, ID string) *model.Space {
	return nil
}

// GetAll is
func (s *SpaceDao) GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.Space, string) {
	return nil, ""
}

// Delete is
func (s *SpaceDao) Delete(dc model.DaoContext, ID string) error {
	return nil
}

// Update is
func (s *SpaceDao) Update(dc model.DaoContext, ID string, obj *model.Space) error {
	return nil
}

// NewSpaceDao is
func NewSpaceDao() *SpaceDao {
	return &SpaceDao{}
}
