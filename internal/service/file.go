package service

type Uploader interface {
	upload()
}

type LocalUploader struct {

}
