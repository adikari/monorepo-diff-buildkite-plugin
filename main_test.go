package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogger(t *testing.T) {
	setupLogger("debug")
	assert.Equal(t, log.GetLevel(), log.DebugLevel)
	setupLogger("weird level")
	assert.Equal(t, log.GetLevel(), log.InfoLevel)
}
