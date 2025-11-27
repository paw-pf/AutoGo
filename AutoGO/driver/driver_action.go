package driver

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Locator — аналог словаря {"name": "...", "XPath": "..."}
type Locator struct {
	Name  string
	XPath string
}

// DriverAction — центральный класс для действий с браузером
type DriverAction struct {
	ctx context.Context
	t   *testing.T
}

func NewDriverAction(ctx context.Context, t *testing.T) *DriverAction {
	return &DriverAction{
		ctx: ctx,
		t:   t,
	}
}

// MakeScreenshot делает скриншот и возвращает PNG-данные
func (d *DriverAction) MakeScreenshot() ([]byte, error) {
	var buf []byte
	err := chromedp.Run(d.ctx,
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось сделать скриншот: %w", err)
	}
	return buf, nil
}

// retry оборачивает функцию с повторными попытками (макс. 5 раз)
func (d *DriverAction) retry(fn func() error) error {
	var err error
	for i := 0; i < 5; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}

// ScrollTo — прокрутка к элементу
func (d *DriverAction) ScrollTo(loc Locator) {
	_ = d.retry(func() error {
		return chromedp.Run(d.ctx,
			chromedp.ScrollIntoView(loc.XPath, chromedp.BySearch),
		)
	})
}

// ClickButton — клик по элементу с ожиданием до 15 секунд
func (d *DriverAction) ClickButton(loc Locator) {
	log.Printf("🖱️ Нажатие на кнопку '%s'...", loc.Name)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.ScrollIntoView(loc.XPath, chromedp.BySearch),
		chromedp.WaitVisible(loc.XPath, chromedp.BySearch),
		chromedp.WaitEnabled(loc.XPath, chromedp.BySearch),
		chromedp.Click(loc.XPath, chromedp.BySearch),
	)
	if err != nil {
		d.t.Fatalf("❌ Не удалось кликнуть по кнопке '%s': %v", loc.Name, err)
	}

	log.Printf("✅ Успешно нажата кнопка '%s'", loc.Name)
}

// FillField — заполнение поля с ожиданием до 15 секунд
func (d *DriverAction) FillField(loc Locator, value string) {
	log.Printf("✏️ Заполнение поля '%s' значением '%s'...", loc.Name, value)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.ScrollIntoView(loc.XPath, chromedp.BySearch),
		chromedp.WaitVisible(loc.XPath, chromedp.BySearch),
		chromedp.Clear(loc.XPath, chromedp.BySearch),
		chromedp.SendKeys(loc.XPath, value, chromedp.BySearch),
		chromedp.KeyEvent("\t"),
	)
	if err != nil {
		d.t.Fatalf("❌ Не удалось заполнить поле '%s': %v", loc.Name, err)
	}

	log.Printf("✅ Поле '%s' успешно заполнено", loc.Name)
}

// FillFieldEnter — заполнение поля с ожиданием до 15 секунд и нажатие Enter
func (d *DriverAction) FillFieldEnter(loc Locator, value string) {
	log.Printf("✏️ Заполнение поля '%s' значением '%s'...", loc.Name, value)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.ScrollIntoView(loc.XPath, chromedp.BySearch),
		chromedp.WaitVisible(loc.XPath, chromedp.BySearch),
		chromedp.Clear(loc.XPath, chromedp.BySearch),
		chromedp.SendKeys(loc.XPath, value, chromedp.BySearch),
		chromedp.KeyEvent("\r"),
	)
	if err != nil {
		d.t.Fatalf("❌ Не удалось заполнить поле '%s': %v", loc.Name, err)
	}

	log.Printf("✅ Поле '%s' успешно заполнено", loc.Name)
}

// GetElement — проверка существования и видимости элемента
func (d *DriverAction) GetElement(loc Locator) {
	log.Printf("🔍 Поиск элемента '%s'...", loc.Name)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(loc.XPath, chromedp.BySearch),
	); err != nil {
		d.t.Fatalf("❌ Элемент '%s' не найден за 15 секунд: %v", loc.Name, err)
	}

	log.Printf("✅ Элемент '%s' найден", loc.Name)
}

// GoToURL — переход по URL
func (d *DriverAction) GoToURL(url string) {
	log.Printf("🌐 Переход на: %s", url)

	err := d.retry(func() error {
		return chromedp.Run(d.ctx,
			chromedp.Navigate(url),
			chromedp.WaitReady("body", chromedp.ByQuery),
		)
	})
	if err != nil {
		log.Printf("❌ Не удалось загрузить страницу %s: %v", url, err)
		d.t.Fatalf("❌ Не удалось перейти на %s: %v", url, err)
	}

	//Нужно для указания локального хранилища если требуется
	//if err := chromedp.Run(d.ctx,
	//	chromedp.Evaluate(`localStorage.setItem('');`, nil),
	//); err != nil {
	//	log.Printf("⚠️ Не удалось установить localStorage: %v", err)
	//} else {
	//	log.Printf("✅ Подсказки отключены")
	//}

	log.Printf("✅ Страница загружена: %s", url)
}

// SwitchFrame — ожидание видимости фрейма до 15 секунд
func (d *DriverAction) SwitchFrame(loc Locator) {
	log.Printf("🖼️ Ожидание фрейма '%s'...", loc.Name)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(loc.XPath, chromedp.BySearch),
	)
	if err != nil {
		d.t.Fatalf("❌ Фрейм '%s' не стал видимым за 15 секунд: %v", loc.Name, err)
	}

	log.Printf("✅ Фрейм '%s' стал видимым", loc.Name)
}

// WaitVisibilityOfAnyElements — ожидание видимости элемента до 15 секунд
func (d *DriverAction) WaitVisibilityOfAnyElements(loc Locator) {
	log.Printf("👁️ Ожидание видимости элемента '%s'...", loc.Name)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(loc.XPath, chromedp.BySearch),
	)
	if err != nil {
		d.t.Fatalf("❌ Элемент '%s' не стал видимым за 15 секунд: %v", loc.Name, err)
	}

	log.Printf("✅ Элемент '%s' стал видимым", loc.Name)
}

// WaitInvisibilityOfElement — ожидание исчезновения элемента до 15 секунд
func (d *DriverAction) WaitInvisibilityOfElement(loc Locator) {
	log.Printf("👻 Ожидание исчезновения элемента '%s'...", loc.Name)

	ctx, cancel := context.WithTimeout(d.ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.WaitNotPresent(loc.XPath, chromedp.BySearch),
	)
	if err != nil {
		d.t.Fatalf("❌ Элемент '%s' не исчез за 15 секунд: %v", loc.Name, err)
	}

	log.Printf("✅ Элемент '%s' исчез", loc.Name)
}

// UploadFile — загрузка файла (работает только для <input type="file">)
func (d *DriverAction) UploadFile(loc Locator, filePath string) {
	log.Printf("📤 Загрузка файла '%s' в поле '%s'", filePath, loc.Name)

	if err := d.retry(func() error {
		return chromedp.Run(d.ctx,
			chromedp.SetUploadFiles(loc.XPath, []string{filePath}, chromedp.BySearch),
		)
	}); err != nil {
		d.t.Fatalf("❌ Не удалось загрузить файл '%s': %v", filePath, err)
	}

	log.Printf("✅ Файл '%s' успешно загружен", filePath)
}

// CheckNotExistElement — проверка отсутствия элемента
func (d *DriverAction) CheckNotExistElement(loc Locator, timeoutSeconds ...int) {
	sec := 3
	if len(timeoutSeconds) > 0 {
		sec = timeoutSeconds[0]
	}

	log.Printf("🚫 Проверка отсутствия элемента '%s' (таймаут: %dс)", loc.Name, sec)

	ctx, cancel := context.WithTimeout(d.ctx, time.Duration(sec)*time.Second)
	defer cancel()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("✅ Элемент '%s' не появился за отведённое время", loc.Name)
			return
		case <-ticker.C:
			var nodes []*cdp.Node
			err := chromedp.Run(d.ctx,
				chromedp.Nodes(loc.XPath, &nodes, chromedp.BySearch),
			)
			if err == nil && len(nodes) > 0 {
				d.t.Fatalf("❌ Ожидалось отсутствие элемента '%s', но он найден", loc.Name)
			}
		}
	}
}
