package testenv

import (
	"autogo/browser"
	"autogo/driver"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
)

const allureResultsDir = "allure-results"

var initOnce sync.Once

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filename))
}

func initAllure() {
	initOnce.Do(func() {
		projectRoot := getProjectRoot()
		allurePath := filepath.Join(projectRoot, allureResultsDir)
		os.RemoveAll(allurePath)
		os.MkdirAll(allurePath, 0755)
		os.Setenv("ALLURE_OUTPUT_PATH", projectRoot)
	})
}

var runTags = flag.String("tags", "", "Список тегов через запятую: smoke,regression,ui,api")

// SkipIfNotTagged проверка существующий тэг
func SkipIfNotTagged(t *testing.T, requiredTag string) {
	if *runTags == "" {
		return
	}

	tags := strings.Split(*runTags, ",")
	for _, tag := range tags {
		if strings.TrimSpace(tag) == requiredTag {
			return
		}
	}

	t.Skipf("Пропущен: требуемый тег '%s' не включён в --tags=%s", requiredTag, *runTags)
}

// StepFunc — функция для шагов в тесте
type StepFunc func(name string, f func())

// RunTestUI — запускает UI-тест с поддержкой шагов и Allure-отчёта
func RunTestUI(
	t *testing.T,
	description string,
	severity allure.SeverityType,
	feature string,
	testFunc func(drv *driver.DriverAction, step StepFunc),
) {
	t.Helper()
	initAllure()

	result := allure.NewResult(t.Name(), t.Name())
	result.Description = description
	result.AddLabel(
		allure.EpicLabel("UI"),
		allure.FeatureLabel(feature),
		allure.SeverityLabel(severity),
	)
	result.Begin()

	ctx, cancel, err := browser.NewContext()
	if err != nil {
		t.Fatalf("Не удалось запустить браузер: %v", err)
	}
	defer cancel()
	drv := driver.NewDriverAction(ctx, t)

	var currentStepName string

	doStep := func(name string, f func()) {
		currentStepName = name

		step := &allure.Step{
			Name:   name,
			Status: "passed",
		}

		defer func() {
			if r := recover(); r != nil {
				step.Status = "failed"
				result.Status = "failed"
				result.StatusDetails = allure.StatusDetail{
					Message: "Тест упал на шаге: \"" + name + "\"",
					Trace:   fmt.Sprintf("Panic: %v", r),
				}
			} else if t.Failed() {
				step.Status = "failed"
				result.Status = "failed"
				result.StatusDetails = allure.StatusDetail{
					Message: "Тест упал на шаге: \"" + name + "\"",
					Trace:   "\"UI-шаг не выполнен: возможно, элемент не загрузился, исчез или страница изменилась.\"",
				}
			}
			result.Steps = append(result.Steps, step)
		}()

		f()
	}

	defer func() {
		if r := recover(); r != nil {
			if screenshot, err := drv.MakeScreenshot(); err == nil {
				attachment := allure.NewAttachment("Скриншот при падении", allure.Png, screenshot)
				result.Attachments = append(result.Attachments, attachment)
			}
			result.Status = "failed"
			result.StatusDetails = allure.StatusDetail{
				Message: "Тест упал на шаге: \"" + currentStepName + "\"",
				Trace:   fmt.Sprintf("Panic в основном потоке: %v", r),
			}
		} else {
			if t.Failed() {
				if screenshot, err := drv.MakeScreenshot(); err == nil {
					attachment := allure.NewAttachment("Скриншот при падении", allure.Png, screenshot)
					result.Attachments = append(result.Attachments, attachment)
				}
				result.Status = "failed"
			} else {
				result.Status = "passed"
			}
		}

		result.Finish()
		if err := result.Print(); err != nil {
			t.Logf("Ошибка при сохранении Allure результата: %v", err)
		}
	}()

	testFunc(drv, doStep)
}

// RunTestAPI — запускает API-тест с поддержкой шагов и Allure-отчёта
func RunTestAPI(
	t *testing.T,
	description string,
	severity allure.SeverityType,
	feature string,
	testFunc func(t *testing.T, step StepFunc, lastBody *[]byte),
) {
	t.Helper()
	initAllure()

	result := allure.NewResult(t.Name(), t.Name())
	result.Description = description
	result.AddLabel(
		allure.EpicLabel("API"),
		allure.FeatureLabel(feature),
		allure.SeverityLabel(severity),
	)
	result.Begin()

	var lastBody []byte

	doStep := func(name string, f func()) {
		failedBefore := t.Failed()
		f()
		stepStatus := "passed"
		if t.Failed() && !failedBefore {
			stepStatus = "failed"
		}
		result.Steps = append(result.Steps, &allure.Step{
			Name:   name,
			Status: allure.Status(stepStatus),
		})
	}
	testFunc(t, doStep, &lastBody)

	defer func() {
		if t.Failed() {
			result.Status = "failed"

			if len(lastBody) > 0 {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, lastBody, "", "  "); err == nil {
					attachment := allure.NewAttachment("Response Body", "application/json", prettyJSON.Bytes())
					result.Attachments = append(result.Attachments, attachment)

					bodyStr := prettyJSON.String()
					if len(bodyStr) > 800 {
						bodyStr = bodyStr[:800] + "\n...\n(обрезано)"
					}
					result.StatusDetails = allure.StatusDetail{
						Message: "Тест завершился с ошибкой.\n\nResponse Body:\n" + bodyStr,
						Trace:   "Полный ответ см. в аттачменте 'Response Body'.",
					}
				} else {
					attachment := allure.NewAttachment("Response Body", "text/plain", lastBody)
					result.Attachments = append(result.Attachments, attachment)
					bodyStr := string(lastBody)
					if len(bodyStr) > 800 {
						bodyStr = bodyStr[:800] + "\n...\n(обрезано)"
					}
					result.StatusDetails = allure.StatusDetail{
						Message: "Тест завершился с ошибкой.\n\nResponse Body:\n" + bodyStr,
						Trace:   "Полный ответ см. в аттачменте 'Response Body'.",
					}
				}
			} else {
				result.StatusDetails = allure.StatusDetail{
					Message: "Тест завершился с ошибкой. Response Body отсутствует.",
					Trace:   "Проверьте логи и шаги теста.",
				}
			}
		} else {
			result.Status = "passed"
		}
		result.Finish()
		if err := result.Print(); err != nil {
			t.Logf("Ошибка при сохранении Allure результата: %v", err)
		}
	}()
}
