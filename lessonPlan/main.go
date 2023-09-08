package main

import (
	"lessonPlan/config"
	"lessonPlan/controller"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	config.DatabaseInit()
	gorm := config.DB

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()

	planRoute := e.Group("/plan")
	planRoute.POST("/", controller.CreatePlan)
	planRoute.POST("/user", controller.CreateUser)
	planRoute.GET("/:id", controller.GetPlan)
	planRoute.GET("/", controller.GetAllPlans)
	planRoute.GET("/week", controller.GetNextWeek)
	planRoute.GET("/month", controller.GetNextMonth)
	planRoute.PUT("/:id", controller.UpdatePlan)
	planRoute.PUT("/user/:id", controller.UpdateUser)
	planRoute.DELETE("/:id", controller.DeletePlan)

	// Zamanlanmış işlemi başlatın
	ticker := time.NewTicker(10 * time.Second) // Örnek: 1 saat aralıklarla kontrol edin
	go func() {
		for {
			select {
			case <-ticker.C:
				// Her zamanlayıcı tetiklendiğinde, bitiş zamanına gelmiş planları kontrol edin ve güncelleyin
				controller.CheckAndUpdatePlans()
			}
		}
	}()

	e.Logger.Fatal(e.Start(":3343"))
}
