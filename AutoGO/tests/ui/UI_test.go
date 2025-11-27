package ui

import (
	"autogo/config"
	"autogo/driver"
	"autogo/locators"
	"autogo/testenv"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
)

func TestUI_Practice_Form(t *testing.T) {
	testenv.SkipIfNotTagged(t, "smoke")
	testenv.SkipIfNotTagged(t, "api")

	testenv.RunTestUI(t,
		"Заполнение формы",
		allure.CRITICAL,
		"Practice Form",
		func(drv *driver.DriverAction, step testenv.StepFunc) {
			step("Открытие страницы входа", func() {
				drv.GoToURL(config.UIBaseURL + "/forms")
			})

			step("Переход на Practice Form", func() {
				drv.ClickButton(locators.PracticeForm.Practice)
			})

			step("Заполнение полей name", func() {
				drv.FillField(locators.PracticeForm.UserFormField("firstName"), "Roman")
				drv.FillField(locators.PracticeForm.UserFormField("lastName"), "Roman")
				drv.FillField(locators.PracticeForm.UserFormField("userEmail"), "name@example.ru")
			})

			step("Чек-бокс Gender", func() {
				drv.ClickButton(locators.PracticeForm.Gender)
			})

			step("Заполнение номера телефона", func() {
				drv.FillField(locators.PracticeForm.Phone, "88005553535")
			})

			step("Заполнение Subjects", func() {
				drv.FillField(locators.PracticeForm.Subjects, "Computer Science")
			})

			//Загрузка файл
			//step("Загрузка файла", func() {
			//	drv.UploadFile(locators.PracticeForm.File, "./tests/file/Autogo.txt")
			//})

			step("Запись State and City", func() {
				drv.FillField(locators.PracticeForm.State, "Haryana")
				drv.FillField(locators.PracticeForm.City, "Panipat")
			})

			step("Нажатие кнопки Submit", func() {
				drv.ClickButton(locators.PracticeForm.Submit)
			})
			step("Проверка отправки", func() {
				drv.GetElement(locators.PracticeForm.Submitting)
			})

		},
	)
}
