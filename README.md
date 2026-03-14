## Custom Log Linter

Это приложение - линтер для программ на языке Go. Работает как плагин для golangci-lint.
Линтер помогает проверять лог записи в коде на соответствие определенным требованиям.

# Технологии: 
- Язык: Go 1.25
- CI/CD: Github Actions 

# Поддерживаемые правила для сообщений логов:
Линтер поддерживает следующие логгеры: `log/slog`, `go.uber.org/zap`
1. Сообщение должно начинаться со строчной буквы
2. Сообщения должны содержать только английский 
3. Сообщения не должны содержать спецсимволы
4. Сообщения не должны содержать эмодзи
5. Сообщения не должны содержать чувствительные данные 

# Установка и использование: 
Проект можно запустить как самостоятельно, так и в качестве плагина для golangci-lint 
1. Запуск через `golangci-lint` (Рекомендуется):
- Необходимо установить `golangci-lint`
- Клонировать репозиторий: 
```
git clone https://github.com/Anna-Moiseeva-3341/custom_linter.git
cd custom_linter
```
- Сборка: 
```
golangci-lint custom
```
- Запустить проверку необходимого проекта: 
```
./custom-gcl run ./...
```
2. Самостоятельный запуск 
- Клонировать репозиторий: 
```
git clone https://github.com/Anna-Moiseeva-3341/custom_linter.git
cd custom_linter
```
- Собрать линтер:
```
go build -o loglint ./cmd/loglint
```
- Запустить проверку:
```
./loglint ./...
```

# Структура проекта: 
- cmd/loglint - точка входа для запуска линтера в самостоятельном режиме (работает с конфигурацией по умолчанию из analyzer.go)
- pkg/loglint/ - ядро линтера, где реализована проверка сообщений логов
- pkg/loglint/testdata/ - содержит файлы с тестами
- plugin/ - обертка для интеграции линтера с golangci-lint
- .github/workflows - CI для автоматической сборки и тестирования

# Тестирование: 
Для проверки соответствия правилам сообщений логов реализованы Unit-тесты с использованием `analysistest`
Для запуска необходимо ввести следующую команду: 
`go test -v ./...`
При тестировании используется конфигурация по умолчанию, описанная в `analyzer.go`

# Конфигурация: 
Линтер можно настроить через конфигурационный файл .`/golangci.yml`. Это работает только для запуска в качестве плагина.
Можно изменять слова, которые считаются чувствительными данными (`forbidden_words`).
Можно изменять символы, которые считаются запрещенными (одиночные символы - `forbidden_symbols`, сочетание символов - `forbidden_patterns`).
Можно изменять правила для проверки (`lowercase`, `language`, `emoji`, `symbols`, `sensitive`). Для использования правила значение с соответствующей переменной должно быть true.
Пример настройки:
```
version: "2"

linters: 
  enable: 
    - customloglint 
  settings: 
    custom: 
      customloglint: 
        type: "module"
        description: "checks log messages"
        settings:
          # Слова, которые считаются чувствительными данными
          forbidden_words: ["password", "token", "secret"]

          # Символы, которые считаются запрещенными
          # одиночные символы
          forbidden_symbols: ["!", "?", ";"]
          # сочетание символов
          forbidden_patterns: ["..."]

          # Правила, которые включены в проверку
          enabled_checks:
            lowercase: true
            language: true
            emoji: true 
            symbols: true
            sensitive: true
```

# Примеры использования: 
Использование плагина:
```
ann@annpc:~/GoFiles/my_linter_project$ ./custom-gcl run ./...
test_run.go:9:12: log messages must start with lowercase letter (customloglint)
	slog.Info("Bad start")
	          ^
test_run.go:10:12: log messages must contain only english (customloglint)
	slog.Warn("русский язык")
	          ^
test_run.go:11:13: log messages must not contain special symbols (customloglint)
	slog.Debug("!!!error")
	           ^
test_run.go:12:12: log messages must not contain emoji (customloglint)
	slog.Info("🔥")
	          ^
test_run.go:13:47: log messages must not contain sensitive data (customloglint)
	slog.Error("my pass is: ", "password_value", pass)
	                                             ^
5 issues:
* customloglint: 5
``` 

Самостоятельное использование:
```
ann@annpc:~/GoFiles/my_linter_project$ go run ./cmd/loglint/ ./...
/home/ann/GoFiles/my_linter_project/test_run.go:9:12: log messages must start with lowercase letter
/home/ann/GoFiles/my_linter_project/test_run.go:10:12: log messages must contain only english
/home/ann/GoFiles/my_linter_project/test_run.go:11:13: log messages must not contain special symbols
/home/ann/GoFiles/my_linter_project/test_run.go:12:12: log messages must not contain emoji
exit status 3
ann@annpc:~/GoFiles/my_linter_project$ 
```