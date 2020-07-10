package tests

import (
	"github.com/bradenrayhorn/switchboard-core/config"
	"github.com/bradenrayhorn/switchboard-core/database"
	"github.com/bradenrayhorn/switchboard-core/routing"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	viper.AddConfigPath("../")
	config.LoadConfig()
	database.Setup()
	redisDB := database.MakeRedisClient()

	r = routing.MakeRouter(redisDB)

	return m.Run()
}
