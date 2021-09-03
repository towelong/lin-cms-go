package jwt

import "github.com/google/wire"

var Set = wire.NewSet(NewJWTMaker)
