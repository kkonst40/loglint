# loglint

Кастомный линтер для Go, предназначенный для проверки стиля и безопасности лог-сообщений.

## Возможности

Линтер проверяет соответствие следующим правилам:
1. **Регистр**: Сообщения должны начинаться со строчной буквы.
2. **Язык**: Сообщения должны быть только на английском языке.
3. **Символы**: Запрещено использование спецсимволов и эмодзи.
4. **Безопасность**: Поиск потенциальных утечек чувствительных данных (password, token, secret и т.д.) через анализ имен переменных.

## Использование с golangci-lint

Этот линтер поддерживает автоматическую интеграцию в качестве модуля.

### Установка golanci-lint

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Настройка `.custom-gcl.yml`

Создать файл .custom-gcl.yml в корневой папке проекта.

```yaml
version: v2.11.3

name: loglint

plugins:
  - module: github.com/kkonst40/loglint
    version: v0.0.8
```

### Настройка `.golangci.yml`

Создать файл .golangci.yml в корневой папке проекта, либо добавить указанные поля в уже существующий файл.

```yaml
version: "2"

linters:
  enable:
    - loglint
  settings:
    custom:
      loglint:
        type: "module"
        settings:
          check_first_char: true
          check_nonenglish_chars: true
          check_special_chars: true
          check_sensitive_words: true
          sensitive_words: [password, pass, token, apikey, api_key]
```

Значения каждого из параметров check_first_char, check_nonenglish_chars, check_special_chars, check_sensitive_words могут быть изменены на false для отключения соответствующей проверки.
В массив sensitive_words могут быть добавлены дополнительные слова, которые линтер будет считать как чувствительные.

### Сборка и запуск кастомного линтера
```bash
golangci-lint custom

./loglint run
```

Все команды созданного кастомного линтера совпадают с командами golangci-lint