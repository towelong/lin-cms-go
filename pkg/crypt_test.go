package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePassword(t *testing.T) {
	password := EncodePassword("123456")
	fmt.Println(password)
	verify := VerifyPsw("123456", password)
	assert.Equal(t, true, verify)
}
