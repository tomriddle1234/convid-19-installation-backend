package main

func initializeRoutes() {
	router.GET("/info", infoHandler)

	apiRoutes := router.Group("/api")
	{
		// POST to http://covidinstal.vfxware.com/api/accumulatesignal
		apiRoutes.POST("/accumulatesignal", accumulateSignalHandler)
		apiRoutes.POST("/reliefsignal", reliefSignalHandler)
		apiRoutes.POST("/resetballon", resetBallonHandler)
		apiRoutes.GET("/statuscheck", statusCheckHandler)
	}
}