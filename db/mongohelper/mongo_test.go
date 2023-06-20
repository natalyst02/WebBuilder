package mongohelper

import (
	"path/filepath"
	"testing"

	"appota/web-builder/config"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	err := config.GetConfig()
	if err != nil {
		log.Warn(".env file not found.")
	}

	h := NewHelper()
	err = h.Connect(filepath.Join(config.GetMongoURI(), "test"))
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
}
