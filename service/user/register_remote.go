package user

import (
	"context"
	"errors"
	model "github.com/cloudreve/Cloudreve/v3/models"
	"strings"
)

// UserRegisterRemoteService 管理用户注册的服务
type UserRegisterRemoteService struct {
	DefaultGroup int
}

// Register 新用户注册
func (service *UserRegisterRemoteService) Register(ctx context.Context) (*model.User, error) {
	defaultGroup := service.DefaultGroup
	email := ctx.Value("email").(string)
	password := ctx.Value("password").(string)
	nickname := ctx.Value("nickname").(string)
	avatar := ctx.Value("avatar").(string)

	if email == "" {
		return nil, errors.New("email is empty")
	}

	if nickname == "" {
		nickname = strings.Split(email, "@")[0]
	}

	// 创建新的用户对象
	user := model.NewUser()
	user.Email = email
	user.Nick = nickname
	user.Avatar = avatar
	user.Status = model.Active
	user.GroupID = uint(defaultGroup)

	user.SetPassword(password)

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		//检查已存在使用者是否尚未激活
		expectedUser, err := model.GetUserByEmail(email)
		if expectedUser.Status == model.NotActivicated {
			user = expectedUser
		} else {
			return nil, err
		}
	}

	return &user, nil
}
