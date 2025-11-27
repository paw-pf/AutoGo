package api

import (
	"fmt"
	"net/http"
	"testing"

	"autogo/api"
	"autogo/config"
	"autogo/testenv"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Transport_Create(t *testing.T) {

	testenv.SkipIfNotTagged(t, "smoke")
	testenv.SkipIfNotTagged(t, "api")

	testenv.RunTestAPI(t,
		"Создание ТС через API",
		allure.CRITICAL,
		"TransportAPI",
		func(t *testing.T, step testenv.StepFunc, lastBody *[]byte) {
			apiClient := api.AuthenticatedClient(t, config.APIBaseURL)
			userid, _ := apiClient.InfoUser()
			fmt.Println(userid)

			var body []byte
			var resp *http.Response

			step("Отправка POST-запроса на создание ТС", func() {
				body, resp = apiClient.POST("/BookStore/v1/Books",
					`{
						  "userId": "`+"1231231"+`",
						  "collectionOfIsbns": [
							{
							  "isbn": "9781449325862"
							}
						  ]
						}`)
				*lastBody = body
			})

			step("Проверка HTTP-статуса", func() {
				assert.Equal(t, 200, resp.StatusCode, "Ожидался статус 200")
			})
		},
	)
}
