package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodePassword(t *testing.T) {
	password := EncodePassword("123456")
	verify := VerifyPsw("123456", password)
	assert.Equal(t, true, verify)
}
