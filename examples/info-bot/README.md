# info-bot

Консольное приложение-пример для демонстрации возможностей библиотеки go-verbosity. Позволяет получать информацию о пользователях, чатах, организациях и отправлять сообщения через командную строку.

## Сборка и запуск

```bash
cd examples/info-bot
go build -o info-bot .
```

## Автоматические сборки

Проект настроен на автоматическую сборку при каждом изменении кода:

- **Linux x64 статический бинарник** собирается автоматически при каждом push/PR
- **Готовые релизы** создаются автоматически при тегах версий
- **Скачать готовые бинарники** можно из раздела [Releases](https://github.com/ivmaks/go-verbosity/releases)

## Локальная сборка

Для удобства разработки создан Makefile с готовыми командами:

```bash
# Показать все доступные команды
make help

# Собрать библиотеку и примеры
make build
make examples

# Собрать статический бинарник для Linux x64
make build-static

# Собрать для всех платформ
make build-multiplatform

# Запустить все проверки как в CI
make ci

# Установить инструменты разработки
make install-tools
```

## Использование

```bash
# Показать справку
./info-bot -help

# Показать версию
./info-bot -version

# Список всех чатов
./info-bot -list-chats -token YOUR_TOKEN

# Список всех организаций
./info-bot -list-orgs -token YOUR_TOKEN

# Информация о пользователе
./info-bot -user-id 123 -token YOUR_TOKEN

# Информация о чате
./info-bot -chat-id 456 -token YOUR_TOKEN

# Информация об организации
./info-bot -org-id 789 -token YOUR_TOKEN

# Найти чат по названию
./info-bot -chat-title "General" -token YOUR_TOKEN

# Статистика чата
./info-bot -chat-stats 456 -token YOUR_TOKEN

# Топ-5 чатов по количеству участников
./info-bot -top-chats-members 5 -token YOUR_TOKEN

# Обновить сообщение в чате
./info-bot -update-chat 123 -update-post 456 -message "Updated text" -token YOUR_TOKEN

# Использование кастомных URL
./info-bot -api-url https://custom-api.example.com -token YOUR_TOKEN

# Вывод в формате JSON
./info-bot -list-chats -token YOUR_TOKEN -output json
```

## Установка info-bot

### Из готовых релизов (рекомендуется)

```bash
# Скачать последнюю версию для Linux x64
wget https://github.com/ivmaks/go-verbosity/releases/latest/download/info-bot-linux-amd64.tar.gz
tar -xzf info-bot-linux-amd64.tar.gz
chmod +x info-bot
./info-bot -help
```

### Из исходного кода

```bash
# Установить последнюю версию через go install
go install github.com/ivmaks/go-verbosity/examples/info-bot@latest

# Или собрать локально
git clone https://github.com/ivmaks/go-verbosity.git
cd go-verbosity/examples/info-bot
go build -o info-bot .
```

## Переменные окружения

```bash
export VERBOSITY_API_TOKEN="your-bot-token"
export VERBOSITY_API_URL="https://api.verbosity.io"
export VERBOSITY_FILE_URL="https://file.verbosity.io"
export VERBOSITY_OUTPUT_MODE="text"  # text, json, json-pretty
```

## Примеры команд

### Получение информации

```bash
# Получить информацию о всех чатах
./info-bot -list-chats -token YOUR_TOKEN

# Получить топ-10 чатов по количеству постов
./info-bot -top-chats-posts 10 -token YOUR_TOKEN

# Найти организацию по названию
./info-bot -org-title "Engineering" -token YOUR_TOKEN

# Получить информацию о пользователе по уникальному имени
./info-bot -user-unique-name "john_doe" -token YOUR_TOKEN
```

### Отправка сообщений

```bash
# Отправить сообщение в чат
./info-bot -send-public 123 -message "Hello from info-bot!" -token YOUR_TOKEN

# Отправить личное сообщение пользователю
./info-bot -send-private 456 -message "Hello!" -token YOUR_TOKEN

# Ответить на сообщение
./info-bot -send-reply 123 789 -message "This is a reply" -token YOUR_TOKEN

# Рассылка в несколько чатов
./info-bot -send-broadcast 123,456,789 -message "Broadcast message" -token YOUR_TOKEN
```

### Обновление сообщений

```bash
# Обновить сообщение в чате
./info-bot -update-chat 123 -update-post 456 -message "Updated message" -token YOUR_TOKEN

# Обновить сообщение с E2E шифрованием
./info-bot -update-chat 123 -update-post 456 -update-e2e -message "E2E updated message" -token YOUR_TOKEN

# Обновить сообщение с ответом
./info-bot -update-chat 123 -update-post 456 -update-reply 789 -message "Updated with reply" -token YOUR_TOKEN

# Обновить сообщение с вложениями
./info-bot -update-chat 123 -update-post 456 -update-attachments "guid1,guid2,guid3" -message "Updated with files" -token YOUR_TOKEN

# Обновить сообщение с цитатой
./info-bot -update-chat 123 -update-post 456 -update-quote "Original message text" -message "Updated with quote" -token YOUR_TOKEN
```

### Настройка вывода

```bash
# Вывод в формате JSON
./info-bot -list-chats -token YOUR_TOKEN -output json

# Использование кастомного API URL
./info-bot -list-chats -token YOUR_TOKEN -api-url https://custom-api.example.com
```