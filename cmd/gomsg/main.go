package main

import (
	"go.uber.org/zap"
	"gomsg/pkg/api"
	"gomsg/pkg/db"
)

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	zap.L().Info("Start")
}

func main() {
	// gin.SetMode(gin.ReleaseMode)
	InitLogger()

	newDb, err := db.NewDb()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	newAPI := api.NewApi(newDb)

	router := newAPI.Start()
	err = router.Run(":8080")
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}
