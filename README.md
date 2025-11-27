# Автоматизированные E2E-тесты

Набор end-to-end (E2E) тестов для UI и API на языке Go с использованием:
- [`chromedp`](https://github.com/chromedp/chromedp) — управление браузером через DevTools Protocol,
- [`allure-go`](https://github.com/ozontech/allure-go) — генерация красивых отчётов с шагами, скриншотами и аттачментами,

Тесты покрывают как API-сценарии, так и пользовательские сценарии в браузере.

---

## 🧪 Требования

- **Go 1.20 или новее**
- **Google Chrome** или **Chromium** (должен быть в `PATH`)
- (Опционально) **Allure CLI** — для просмотра отчётов

> 💡 Убедитесь, что Go и Chrome установлены:
> ```bash
> go version
> google-chrome --version  # или chromium --version
> ```

---

## ⚙️ Настройка окружения

1. Склонируйте проект:
   ```bash
   git clone <ваш-репозиторий>
   cd autogo
   
2. Создайте файл .env на основе шаблона:
   ```bash
   cp .env.example .env
   
3. Отредактируйте .env, указав свои параметры:
   ```bash
    UI_BASE_URL=https://demoqa.com
    API_BASE_URL=https://demoqa.com
    USERNAME=useruser
    PASSWORD=P@ssw0rd
    HEADLESS=true
---

## ▶️ Запуск тестов

### Все тесты

    go test ./tests/... -v

### Только API-тесты

    go test ./tests/api -v

### Только UI-тесты

    go test ./tests/ui -v

### Запуск по тегам
    # Только smoke-тесты
    go test ./tests/... -v --tags=smoke
    
    # Smoke + API
    go test ./tests/... -v --tags=smoke,api
    
    # UI-тесты с тегом regression
    go test ./tests/ui -v --tags=regression,ui

## 📊 Просмотр отчётов Allure
```bash 
# macOS (Homebrew)
brew install allure

# Ubuntu/Debian
sudo apt-add-repository ppa:qameta/allure
sudo apt-get update
sudo apt-get install allure

# Сгенерируйте и откройте отчёт
allure serve allure-results


