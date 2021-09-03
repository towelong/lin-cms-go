package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../../config")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 3306, config.MySQL.Port)
	assert.Equal(t, 3306, viper.GetInt("mysql.port"))
	assert.NotEqual(t, int64(0), config.Lin.CMS.TokenAccessExpire)
	assert.NotEqual(t, int64(0), config.Lin.CMS.TokenRefreshExpire)
	assert.NotEqual(t, "", config.Lin.CMS.TokenSecret)
}

