package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"autogo/api"
	"autogo/config"
	"autogo/testenv"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Book(t *testing.T) {

	testenv.SkipIfNotTagged(t, "smoke")
	testenv.SkipIfNotTagged(t, "api")
	testenv.SkipIfNotTagged(t, "regression")

	testenv.RunTestAPI(t,
		"Получение книг",
		allure.CRITICAL,
		"BookAPI",
		func(t *testing.T, step testenv.StepFunc, lastBody *[]byte) {
			apiClient := api.AuthenticatedClient(t, config.APIBaseURL)

			var body []byte
			var resp *http.Response

			step("Отправка POST-запроса на получения списка тс", func() {
				body, resp = apiClient.GET("/BookStore/v1/Books")
				*lastBody = body
			})

			step("Проверка HTTP-статуса", func() {
				var response map[string]interface{}
				assert.NoError(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, resp.StatusCode, "Ожидался статус 200")
				isbn := response["books"].([]interface{})[0].(map[string]interface{})["isbn"].(string)
				assert.Equal(t, "9781449325862", isbn)
			})
		},
	)
}
