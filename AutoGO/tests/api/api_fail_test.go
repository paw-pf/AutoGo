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

func TestAPI_fail_Transport_Create(t *testing.T) {

	testenv.SkipIfNotTagged(t, "smoke")
	testenv.SkipIfNotTagged(t, "api")

	testenv.RunTestAPI(
		t,
		"Создание ТС через API",
		allure.CRITICAL,
		"TransportAPI",
		func(t *testing.T, step testenv.StepFunc, lastBody *[]byte) {
			apiClient := api.NewClient(t, config.APIBaseURL)
			apiClient = apiClient.WithToken(config.Token)

			var body []byte
			var resp *http.Response

			step("Отправка POST-запроса на создание ТС", func() {
				body, resp = apiClient.POST("/v3/companies/9697/cars",
					`{
						"sts": "9982214973",
						"grz": "Е032АУ250",
						"unit_id": "cf9931cf-98c8-4640-b7fa-6b216ce9f030",
						"name": "audi",
						"vin": "WAUZZZ8V5LA024808",
						"product_ids": ["all_toll_roads_monitor"],
						"no_default_products": true,
						"sub_auto_upgrade": true
					}`)
				*lastBody = body
			})

			step("Проверка HTTP-статуса", func() {
				assert.Equal(t, 200, resp.StatusCode, "Ожидался статус 200")
			})
			step("Проверка VIN в ответе", func() {
				var response map[string]interface{}
				assert.NoError(t, json.Unmarshal(body, &response))
				assert.Equal(t, "WAU11ZZZ8V5LA024808", response["vin"])
			})
		},
	)
}
