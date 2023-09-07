package controller

import (
	"fmt"
	"lessonPlan/config"
	"lessonPlan/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var datePlan []time.Time
var UserNames []string
var StudentNos []int

func State(c int) {
	switch c {
	case 0:
		fmt.Println("Program ekleniyor.")
	case 1:
		fmt.Println("Program basariyla eklendi.")
	case 2:
		fmt.Println("Hata olustu.İslem iptal edildi.")
	case 3:
		fmt.Println("Program siliniyor.")
	case 4:
		fmt.Println("Program basariyla silindi.")
	case 5:
		fmt.Println("Program guncelleniyor.")
	case 6:
		fmt.Println("Program basariyla guncellendi.")
	}
}

func CreatePlan(c echo.Context) error {
	b := new(model.Plan)
	db := config.DB
	State(0)

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)

	}

	if CheckUsers(b.Name) == false {
		return c.String(http.StatusInternalServerError, "Kullanıcı adı bulunamadı")
	}
	for _, s := range datePlan {
		if s.Year() == b.Date.Year() && s.Month() == b.Date.Month() && b.Date.Day() == s.Day() && b.Date.Hour() == s.Hour() {
			for _, v := range UserNames {
				if v == b.Name {
					State(2)
					return c.String(http.StatusInternalServerError, "Aynı zaman aralığında farklı bir plan olduğu için kayıt iptal edildi.")
				}
			}
		}
	}
	datePlan = append(datePlan, b.Date)
	if err := db.Create(&b).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": b,
	}
	State(1)
	return c.JSON(http.StatusOK, response)
}

func UpdatePlan(c echo.Context) error {
	id := c.Param("id")
	b := new(model.Plan)
	db := config.DB
	State(5)

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	if CheckUsers(b.Name) == false {
		return c.String(http.StatusInternalServerError, "Kullanıcı adı bulunamadı")
	}

	existing_Plan := new(model.Plan)

	if err := db.First(&existing_Plan, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusNotFound, data)
	}

	existing_Plan.Name = b.Name
	existing_Plan.Description = b.Description
	existing_Plan.Date = b.Date

	if err := db.Save(&existing_Plan).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_Plan,
	}
	State(6)
	return c.JSON(http.StatusOK, response)
}

func GetPlan(c echo.Context) error {
	id := c.Param("id")
	db := config.DB
	var plans []*model.Plan

	if res := db.Find(&plans, id); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}
		State(2)
		return c.JSON(http.StatusNotFound, data)
	}
	fmt.Println(plans[0])
	response := map[string]interface{}{
		"data": plans[0],
	}

	return c.JSON(http.StatusOK, response)
}

func GetAllPlans(c echo.Context) error {
	plans := []model.Plan{}
	config.DB.Find(&plans)
	return c.JSON(http.StatusOK, plans)
}

func GetNextWeek(c echo.Context) error {
	plans := []model.Plan{}
	now := time.Now()
	next := now.AddDate(0, 0, +7)
	config.DB.Where("DATE BETWEEN ? AND ?", now, next).Find(&plans)
	return c.JSON(http.StatusOK, plans)
}

func GetNextMonth(c echo.Context) error {
	plans := []model.Plan{}
	now := time.Now()
	next := now.AddDate(0, +1, 0)
	config.DB.Where("DATE BETWEEN ? AND ?", now, next).Find(&plans)
	return c.JSON(http.StatusOK, plans)
}

func DeletePlan(c echo.Context) error {
	id := c.Param("id")
	db := config.DB
	plan := new(model.Plan)

	State(3)

	if err := db.First(&plan, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusNotFound, data)
	}
	if CheckUsers(plan.Name) == false {
		return c.String(http.StatusInternalServerError, "Kullanıcı adı bulunamadı")
	}
	plan.State = "canceled"
	db.Save(&plan)

	err := db.Delete(&plan, id).Error
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	//plan.State = "canceled"
	response := map[string]interface{}{
		"message": "Kayit başariyla silindi",
	}
	State(4)
	return c.JSON(http.StatusOK, response)

}

func CheckAndUpdatePlans() {
	now := time.Now()
	var plans []model.Plan

	// Bitiş zamanına gelmiş planları veritabanından alın
	if err := config.DB.Where("DATE <= ? AND state = ?", now, "processing").Find(&plans).Error; err != nil {
		// Hata işleme
		return
	}

	// Bitiş zamanına gelmiş planları güncelle (states = "bitti")
	for _, plan := range plans {
		plan.State = "finished"
		if err := config.DB.Save(&plan).Error; err != nil {
			// Hata işleme
			return
		}
	}
}
func AddUser(c echo.Context) error {
	a := new(model.User)
	db := config.DB
	State(0)

	if err := c.Bind(a); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)

	}

	for _, s := range StudentNos {
		if s == a.StudentNo {
			fmt.Println(s)
			State(2)
			return c.String(http.StatusInternalServerError, "Öğrenci numarası kayıtlıdır.Lütfen tekrar deneyiniz.")
		}
	}
	StudentNos = append(StudentNos, a.StudentNo)
	UserNames = append(UserNames, a.Name)
	fmt.Println(UserNames)
	if err := db.Create(&a).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": a,
	}
	State(1)
	return c.JSON(http.StatusOK, response)
}

func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	d := new(model.User)
	db := config.DB
	State(5)

	// Binding data
	if err := c.Bind(d); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	if CheckUsers(d.Name) == false {
		return c.String(http.StatusInternalServerError, "Kullanıcı adı bulunamadı")
	}

	existing_Plan := new(model.User)

	if err := db.First(&existing_Plan, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusNotFound, data)
	}

	existing_Plan.Name = d.Name

	if err := db.Save(&existing_Plan).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		State(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_Plan,
	}
	State(6)
	return c.JSON(http.StatusOK, response)
}

func CheckUsers(c string) bool {
	sayac := false
	for _, v := range UserNames {
		if v == c {
			sayac = true
		}
	}
	if sayac == false {
		// Hata işleme
		return false
	} else {
		return true
	}
}
