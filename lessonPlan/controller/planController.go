package controller

import (
	"fmt"
	"lessonPlan/config"
	"lessonPlan/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var users []model.User
var plans []model.Plan

// Plan oluşturma işlemi yapılır
func CreatePlan(c echo.Context) error {
	b := new(model.Plan)
	db := config.DB
	Status(0)
	GetAllDates()

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)

	}

	// Kullanıcı mevcut diye kontrol edilir
	if !CheckUsers(b.Name) {
		Status(2)
		return c.String(http.StatusInternalServerError, "Kullanıcı adı bulunamadı")
	}
	// Aynı tarihte ve saatte plan mevcut mu diye kontrol yapılır
	for _, s := range plans {
		if s.Date.Year() == b.Date.Year() && s.Date.Month() == b.Date.Month() && b.Date.Day() == s.Date.Day() && b.Date.Hour() == s.Date.Hour() && s.Name == b.Name {
			Status(2)
			return c.String(http.StatusInternalServerError, "Aynı zaman aralığında farklı bir plan olduğu için kayıt iptal edildi.")
		}
	}
	// Herhangi bir hata yoksa plan oluşturulur
	if err := db.Create(&b).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": b,
	}
	Status(1)
	return c.JSON(http.StatusOK, response)
}

// Plan güncelleme işlemi yapılır
func UpdatePlan(c echo.Context) error {
	id := c.Param("id")
	b := new(model.Plan)
	db := config.DB
	Status(0)

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Kullanıcı mevcut diye kontrol edilir
	if !CheckUsers(b.Name) {
		Status(2)
		return c.String(http.StatusInternalServerError, "Kullanıcı adı bulunamadı")
	}

	existing_Plan := new(model.Plan)

	// Güncellenmek istenen id değeri bulunur
	if err := db.First(&existing_Plan, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusNotFound, data)
	}
	// Değerler güncellenir
	existing_Plan.Name = b.Name
	existing_Plan.Description = b.Description
	existing_Plan.Date = b.Date

	// Değişiklik yapılır ve databasede değerler güncellenir
	if err := db.Save(&existing_Plan).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_Plan,
	}
	Status(1)
	return c.JSON(http.StatusOK, response)
}

// Plan silme işlemi yapılır
func DeletePlan(c echo.Context) error {
	id := c.Param("id")
	db := config.DB
	plan := new(model.Plan)

	Status(0)

	// Silmek istenilen planın id değeri kontrol edilir
	if err := db.First(&plan, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusNotFound, data)
	}
	plan.State = "canceled"
	db.Save(&plan)

	// Plan silinir
	err := db.Delete(&plan, id).Error
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "Kayit basariyla silindi",
	}
	Status(1)
	return c.JSON(http.StatusOK, response)

}

// İstenen planın değeri listelenir
func GetPlan(c echo.Context) error {
	id := c.Param("id")
	db := config.DB
	var plans []*model.Plan
	Status(0)

	// Görüntülenmek istenen id değeri bulunur
	if res := db.Find(&plans, id); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}
		Status(2)
		return c.JSON(http.StatusNotFound, data)
	}
	response := map[string]interface{}{
		"data": plans[0],
	}
	Status(1)
	return c.JSON(http.StatusOK, response)
}

// Yeni bir kullanıcı oluşturur
func CreateUser(c echo.Context) error {
	a := new(model.User)
	db := config.DB
	Status(0)

	// Binding data
	if err := c.Bind(a); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)

	}

	// Öğrenci numarası sistemde mevcut diye mi kontrol edilir
	for _, s := range GetAllNOs() {
		if s == a.StudentNo {
			Status(2)
			return c.String(http.StatusInternalServerError, "Öğrenci numarası kayıtlıdır.Lütfen tekrar deneyiniz.")
		}
	}
	// Kullanıcı oluşturulur
	if err := db.Create(&a).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": a,
	}
	Status(1)
	return c.JSON(http.StatusOK, response)
}

// Kullanıcının değerleri güncellemesi yapılır
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	d := new(model.User)
	db := config.DB
	Status(0)

	// Binding data
	if err := c.Bind(d); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}
	// Öğrenci numarası sistemde mevcut diye mi kontrol edilir
	if CheckNumbers(d.StudentNo) {
		return c.String(http.StatusInternalServerError, "Ogrenci numarasi sistemde mevcut.")
	}

	existing_Plan := new(model.User)

	// Güncellenmek istenilen planın id değeri kontrol edilir
	if err := db.First(&existing_Plan, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusNotFound, data)
	}

	// Değerler güncellenir
	existing_Plan.Name = d.Name
	existing_Plan.StudentNo = d.StudentNo

	// Değişiklik yapılır ve databasede değerler güncellenir
	if err := db.Save(&existing_Plan).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}
		Status(2)
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_Plan,
	}
	Status(1)
	return c.JSON(http.StatusOK, response)
}

// Tüm planlar listelenir
func GetAllPlans(c echo.Context) error {
	plans := []model.Plan{}
	config.DB.Find(&plans)
	return c.JSON(http.StatusOK, plans)
}

// Haftalık planlar listelenir
func GetNextWeek(c echo.Context) error {
	plans := []model.Plan{}
	now := time.Now()
	next := now.AddDate(0, 0, +7)
	config.DB.Where("DATE BETWEEN ? AND ?", now, next).Find(&plans)
	return c.JSON(http.StatusOK, plans)
}

// Aylık planlar listelenir
func GetNextMonth(c echo.Context) error {
	plans := []model.Plan{}
	now := time.Now()
	next := now.AddDate(0, +1, 0)
	config.DB.Where("DATE BETWEEN ? AND ?", now, next).Find(&plans)
	return c.JSON(http.StatusOK, plans)
}

// State kontrolü yapar
func CheckAndUpdatePlans() {
	now := time.Now()
	var plans []model.Plan

	// Bitiş zamanına gelmiş planları veritabanından alır
	if err := config.DB.Where("DATE <= ? AND state = ?", now, "processing").Find(&plans).Error; err != nil {
		return
	}

	// Bitiş zamanına gelmiş planların state verilerini günceller
	for _, plan := range plans {
		plan.State = "finished"
		if err := config.DB.Save(&plan).Error; err != nil {
			return
		}
	}
}

// Databasedeki tüm isimleri alır
func GetAllNames() []string {
	db := config.DB
	db.Find(&users)
	var UserNames []string
	for _, v := range users {
		UserNames = append(UserNames, v.Name)
	}
	return UserNames
}

// Databasedeki tüm date bilgilerini alır
func GetAllDates() []time.Time {
	db := config.DB
	db.Find(&plans)
	var dates []time.Time
	for _, v := range plans {
		dates = append(dates, v.Date)
	}
	return dates
}

// Databasedeki tüm okul numarası bilgilerini alır
func GetAllNOs() []int {
	db := config.DB
	db.Find(&users)
	var studentNumbers []int
	for _, v := range users {
		studentNumbers = append(studentNumbers, v.StudentNo)
	}
	return studentNumbers
}

// Kullanıcı kontrolü yapar
func CheckUsers(c string) bool {
	sayac := false
	for _, v := range GetAllNames() {
		if v == c {
			sayac = true
		}
	}
	return sayac
}

// Okul numarası kontrolü yapar
func CheckNumbers(c int) bool {
	sayac := false
	for _, v := range GetAllNOs() {
		if v == c {
			sayac = true
		}
	}
	return sayac
}

// Konsol ekranına durum bilgilendirmesi yapılır
func Status(c int) {
	switch c {
	case 0:
		fmt.Println("İşlem yapılıyor.")
	case 1:
		fmt.Println("İşlem başarılı.")
	case 2:
		fmt.Println("Hata olustu.İslem iptal edildi.")
	}
}
