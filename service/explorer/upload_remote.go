package explorer

import (
	"context"
	"errors"
	model "github.com/cloudreve/Cloudreve/v3/models"
	"github.com/cloudreve/Cloudreve/v3/pkg/filesystem"
	"github.com/cloudreve/Cloudreve/v3/pkg/filesystem/fsctx"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
	"io/ioutil"
	"strings"
	"time"
)

// CreateUploadSessionRemoteService 获取上传凭证服务
type CreateUploadSessionRemoteService struct {
}

// Create 创建新的上传会话
func (service *CreateUploadSessionRemoteService) Create(ctx context.Context) (*serializer.UploadCredential, error) {
	email := ctx.Value("email").(string)
	size := ctx.Value("size").(uint64)
	name := ctx.Value("name").(string)
	path := ctx.Value("path").(string)
	lastModified := ctx.Value("last_modified").(int64)

	user, err := model.GetUserByEmail(email)

	if err != nil {
		return nil, errors.New("user not found")
	}

	fs, err := filesystem.NewFileSystem(&user)

	if err != nil {
		return nil, errors.New("filesystem error")
	}

	file := &fsctx.FileStream{
		Size:        size,
		Name:        name,
		VirtualPath: path,
		File:        ioutil.NopCloser(strings.NewReader("")),
	}

	if lastModified > 0 {
		lastModified := time.UnixMilli(lastModified)
		file.LastModified = &lastModified
	}

	credential, err := fs.CreateUploadSession(ctx, file)

	if err != nil {
		return nil, err
	}

	return credential, nil
}
