package file

import (
	"github.com/gin-gonic/gin"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
	"github.com/towelong/lin-cms-go/pkg"
)

type FileConfig struct {
}

type Uploader interface {
	Upload(ctx *gin.Context) ([]vo.FileVo, error)
}

type DefaultUploader struct {
}

func (d DefaultUploader) GetFileType() string {
	return LOCAL
}

func (d DefaultUploader) GetStorePath(fileName string) string {
	storeDir, err := pkg.CreateDirAndFileForCurrentTime("assets", "2006/01/02")
	if err != nil {
		log.Logger.Error(err.Error())
	}
	return storeDir
}
