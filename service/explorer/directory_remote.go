package explorer

import (
	"context"
	"errors"
	model "github.com/cloudreve/Cloudreve/v3/models"

	"github.com/cloudreve/Cloudreve/v3/pkg/filesystem"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
)

// DirectoryRemoteService 创建新目录服务
type DirectoryRemoteService struct {
}

// ListDirectory 列出目录内容
func (service *DirectoryRemoteService) ListDirectory(c context.Context) (*serializer.ObjectList, error) {
	email := c.Value("email").(string)
	path := c.Value("path").(string)

	user, err := model.GetUserByEmail(email)
	fs, err := filesystem.NewFileSystem(&user)

	if err != nil {
		return nil, err
	}
	defer fs.Recycle()
	// 上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取子项目
	objects, err := fs.List(ctx, path, nil)

	if err != nil {
		return nil, err
	}

	var parentID uint
	if len(fs.DirTarget) > 0 {
		parentID = fs.DirTarget[0].ID
	}
	list := serializer.BuildObjectList(parentID, objects, fs.Policy)
	return &list, nil
}

// CreateDirectory 创建目录
func (service *DirectoryRemoteService) CreateDirectory(c context.Context) (*model.Folder, error) {
	email := c.Value("email").(string)
	path := c.Value("path").(string)

	user, err := model.GetUserByEmail(email)
	// 创建文件系统
	fs, err := filesystem.NewFileSystem(&user)
	if err != nil {
		return nil, errors.New("create filesystem error")
	}
	defer fs.Recycle()

	// 上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建目录
	folder, err := fs.CreateDirectory(ctx, path)

	if err != nil {
		return nil, err
	}

	return folder, nil
}
