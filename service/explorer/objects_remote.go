package explorer

import (
	"context"
	"errors"
	model "github.com/cloudreve/Cloudreve/v3/models"
	"github.com/cloudreve/Cloudreve/v3/pkg/filesystem"
	"github.com/cloudreve/Cloudreve/v3/pkg/hashid"
)

// ItemRemoteService 处理多文件/目录相关服务
type ItemRemoteService struct {
}

// Delete 删除对象
func (service *ItemRemoteService) Delete(c context.Context) (interface{}, error) {
	email := c.Value("email").(string)

	dirs := c.Value("dirs").([]string)
	items := c.Value("items").([]string)

	Dirs := make([]uint, 0, len(dirs))
	Items := make([]uint, 0, len(items))

	for _, folder := range dirs {
		id, err := hashid.DecodeHashID(folder, hashid.FolderID)
		if err == nil {
			Dirs = append(Dirs, id)
		}
	}
	for _, file := range items {
		id, err := hashid.DecodeHashID(file, hashid.FileID)
		if err == nil {
			Items = append(Items, id)
		}
	}

	// 获取用户
	user, err := model.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	// 创建文件系统
	fs, err := filesystem.NewFileSystem(&user)

	if err != nil {
		return nil, errors.New("create filesystem error")
	}
	defer fs.Recycle()

	force, unlink := false, false

	//if fs.User.Group.OptionsSerialized.AdvanceDelete {
	//	force = service.Force
	//	unlink = service.UnlinkOnly
	//}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = fs.Delete(ctx, Dirs, Items, force, unlink)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Rename 重命名对象
func (service *ItemRemoteService) Rename(c context.Context) (interface{}, error) {
	email := c.Value("email").(string)
	dirs := c.Value("dirs").([]string)
	items := c.Value("items").([]string)
	newName := c.Value("new_name").(string)

	Dirs := make([]uint, 0, len(dirs))
	Items := make([]uint, 0, len(items))

	for _, folder := range dirs {
		id, err := hashid.DecodeHashID(folder, hashid.FolderID)
		if err == nil {
			Dirs = append(Dirs, id)
		}
	}
	for _, file := range items {
		id, err := hashid.DecodeHashID(file, hashid.FileID)
		if err == nil {
			Items = append(Items, id)
		}
	}

	// 获取用户
	user, err := model.GetUserByEmail(email)

	if err != nil {
		return nil, errors.New("user not found")
	}

	// 创建文件系统
	fs, err := filesystem.NewFileSystem(&user)

	if err != nil {
		return nil, errors.New("create filesystem error")
	}

	defer fs.Recycle()

	// 重命名作只能对一个目录或文件对象进行操作
	if len(Items)+len(Dirs) > 1 {
		return nil, errors.New("can only rename one object")
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	// 重命名对象
	err = fs.Rename(ctx, Dirs, Items, newName)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
