package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// to minimize test log
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
