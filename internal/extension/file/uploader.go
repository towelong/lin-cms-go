package file

import (
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
	"github.com/towelong/lin-cms-go/pkg"
	"github.com/towelong/lin-cms-go/pkg/response"
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

func (d DefaultUploader) IsValid(files []*multipart.FileHeader) error {
	nums := viper.GetInt("lin.file.nums")
	if len(files) > nums {
		return response.NewResponse(10121)
	}
	for _, file := range files {
		singleLimit := viper.GetInt64("lin.file.singleLimit") * 1024 * 1024
		if file.Size > singleLimit {
			return response.NewResponse(10110)
		}
		extension := "." + strings.Split(file.Filename, ".")[1]
		include := viper.GetStringSlice("lin.file.include")
		flag := false
		for _, i := range include {
			if i == extension {
				flag = true
			}
		}
		if !flag {
			return response.NewResponse(10130)
		}

	}

	return nil
}
