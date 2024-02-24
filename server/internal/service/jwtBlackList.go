package service

import (
	"errors"
	"fqhWeb/internal/model"

	"gorm.io/gorm"
)

type JwtService struct {
}

var (
	insJwt = JwtService{}
)

func Jwt() *JwtService {
	return &insJwt
}

// JwtAddBlacklist
// @description: 拉黑jwt
// @param: jwtList model.JwtBlacklist
// @return: err error
func (s *JwtService) JwtAddBlacklist(jwtList *model.JwtBlacklist) (err error) {
	err = model.DB.Create(jwtList).Error
	return
}

// @function: IsBlacklist
// @description: 判断JWT是否在黑名单内部
// @param: auth string
// @return: bool
func (s *JwtService) IsBlacklist(jwt string) bool {
	err := model.DB.Where("jwt = ?", jwt).First(&model.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}
