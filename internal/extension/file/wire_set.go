package file

import "github.com/google/wire"

var Set = wire.NewSet(
	wire.Struct(new(LocalUploader), "*"),
	wire.Bind(new(Uploader), new(*LocalUploader)),
)
