package file

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/towelong/lin-cms-go/internal/domain/model"
	"github.com/towelong/lin-cms-go/internal/domain/vo"
	"github.com/towelong/lin-cms-go/internal/pkg/log"
	"github.com/towelong/lin-cms-go/internal/service"
	"github.com/towelong/lin-cms-go/pkg"
	"github.com/towelong/lin-cms-go/pkg/response"
)

type LocalUploader struct {
	FileService service.IFileService
}

func (l LocalUploader) Upload(ctx *gin.Context) ([]vo.FileVo, error) {
	fileVos := make([]vo.FileVo, 0)
	form, err := ctx.MultipartForm()
	if err != nil {
		return fileVos, response.NewResponse(10190)
	}
	files := form.File["file"]
	d := DefaultUploader{}
	if err := d.IsValid(files); err != nil {
		return fileVos, err
	}
	for _, file := range files {
		md5 := pkg.GetFileMd5(file)
		fileInDB, err := l.FileService.GetFileByMD5(md5)
		// 如果存在
		var fileVo vo.FileVo
		if err == nil {
			copier.CopyWithOption(&fileVo, &fileInDB, copier.Option{IgnoreEmpty: true})
			fileVo.ID = fileInDB.ID
			fileVo.URL = viper.GetString("lin.file.domain") + viper.GetString("lin.file.storeDir") + fileVo.Path
			fileVos = append(fileVos, fileVo)
		} else {
			dst, err := pkg.CreateDirAndFileForCurrentTime(viper.GetString("lin.file.storeDir"), "2006/01/02")
			if err != nil {
				log.Logger.Error(err.Error())
			}
			extension := "." + strings.Split(file.Filename, ".")[1]
			path := time.Now().Format("2006/01/02") + "/" + md5 + extension
			url := viper.GetString("lin.file.domain") + viper.GetString("lin.file.storeDir") + path
			file.Filename = md5 + extension
			fileVo = vo.FileVo{
				Path:      path,
				Type:      LOCAL,
				Name:      md5 + extension,
				Extension: extension,
				Size:      int(file.Size),
				Md5:       md5,
				URL:       url,
			}
			var fileModel model.File
			copier.CopyWithOption(&fileModel, &fileVo, copier.Option{IgnoreEmpty: true})
			ctx.SaveUploadedFile(file, dst+"/"+file.Filename)
			f, err := l.FileService.CreateFile(fileModel)
			if err == nil {
				fileVo.ID = f.ID
				fileVos = append(fileVos, fileVo)
			}
		}
	}
	return fileVos, nil
}
