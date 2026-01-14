# go-verbosity

Go модуль для работы с Verbosity Bot API.

## Установка

```bash
go get github.com/ivmaks/go-verbosity/verbosity
```

## Быстрый старт

```go
package main

import (
    "fmt"
    "github.com/ivmaks/go-verbosity/verbosity"
)

func main() {
    // Создаем клиент с настройками из переменных окружения
    client := verbosity.NewClientFromEnv()

    // Получаем список чатов
    chats, err := client.GetChatIDs()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d chats\n", len(chats.Chats))
}
```

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `VERBOSITY_API_URL` | API URL | `https://api.verbosity.io` |
| `VERBOSITY_FILE_URL` | URL для загрузки файлов | `https://file.verbosity.io` |
| `VERBOSITY_API_TOKEN` | Токен бота | (обязательно) |

### Кастомная конфигурация

```go
config := &verbosity.Config{
    APIURL:   "https://custom-api.example.com",
    FileURL:  "https://custom-file.example.com",
    APIToken: "your-bot-token",
}

client := verbosity.NewClient(config)
```

## API Методы

### Пользователи

```go
// Получить пользователей по ID
users, err := client.GetUsersByIDs([]int64{1, 2, 3})

// Получить пользователей по уникальным именам
users, err := client.GetUsersByUniqueNames([]string{"user1", "user2"})

// Получить одного пользователя по ID
user, err := client.GetUserByID(123)

// Получить одного пользователя по уникальному имени
user, err := client.GetUserByUniqueName("john")
```

### Чаты

```go
// Получить все ID чатов
chats, err := client.GetChatIDs()

// Получить информацию о чатах по ID
chats, err := client.GetChatsByIDs([]int64{1, 2, 3})

// Получить информацию о конкретном чате
chat, err := client.GetChatByID(456)

// Получить все чаты
allChats, err := client.GetAllChats()

// Найти чат по названию
chat, err := client.FindChatByTitle("General")

// Получить список участников чата
members, err := client.ChatMemberIDs(chatID)

// Получить список админов чата
admins, err := client.ChatAdminIDs(chatID)

// Проверить является ли пользователь участником/админом
isMember, err := client.IsChatMember(chatID, userID)
isAdmin, err := client.IsChatAdmin(chatID, userID)

// Получить чаты, где бот является участником
myChats, err := client.GetMyChats()

// Получить избранные чаты
favorites, err := client.GetFavoriteChats()

// Топ чатов по количеству участников
topChats, err := client.GetTopChatsByMembers(10)

// Топ чатов по количеству постов
topChats, err := client.GetTopChatsByPosts(10)
```

### Организации (команды)

```go
// Получить все ID организаций
orgs, err := client.GetOrganizationIDs()

// Получить информацию об организациях по ID
orgs, err := client.GetOrganizationsByIDs([]int64{1, 2, 3})

// Получить информацию о конкретной организации
org, err := client.GetOrganizationByID(789)

// Получить все организации
allOrgs, err := client.GetAllOrganizations()

// Найти организацию по названию
org, err := client.FindOrganizationByTitle("Engineering")

// Найти организацию по slug
org, err := client.FindOrganizationBySlug("engineering")

// Получить список участников организации
members, err := client.OrganizationMembers(orgID)

// Получить список админов организации
admins, err := client.OrganizationAdmins(orgID)

// Получить организации, где бот является участником
myOrgs, err := client.GetMyOrganizations()

// Получить организации, где бот является админом
adminOrgs, err := client.GetAdminOrganizations()

// Топ организаций по количеству пользователей
topOrgs, err := client.GetTopOrgsByUsers(10)
```

### Отправка сообщений

```go
// Отправить сообщение в чат
response, err := client.SendMessage(chatID, "Hello, world!", nil)

// Отправить ответ на сообщение
response, err := client.SendReply(chatID, postNo, "This is a reply")

// Отправить личное сообщение по ID пользователя
response, err := client.SendPrivateMessageByID(userID, "Hello!", nil)

// Отправить личное сообщение по email
response, err := client.SendPrivateMessageByEmail("user@example.com", "Hello!", nil)

// Отправить личное сообщение по уникальному имени
response, err := client.SendPrivateMessageByUniqueName("username", "Hello!", nil)

// Рассылка сообщений в несколько чатов
responses, err := client.BroadcastMessage([]int64{chatID1, chatID2}, "Hello all!")

// Отправить сообщение во все чаты бота
responses, err := client.SendMessageToAllMyChats("Hello everyone!")
```

### Загрузка файлов

```go
// Загрузить файл
response, err := client.UploadFile(chatID, "/path/to/file.txt")

// Загрузить файл из байтов
response, err := client.UploadFileFromBytes(chatID, data, "filename.txt")

// Загрузить текстовый файл
response, err := client.UploadTextFile(chatID, "file content", "file.txt")

// Загрузить изображение
response, err := client.UploadImage(chatID, "/path/to/image.png")

// Загрузить документ
response, err := client.UploadDocument(chatID, "/path/to/document.pdf")

// Загрузить аудио/видео
response, err := client.UploadAudio(chatID, "/path/to/audio.mp3")
response, err := client.UploadVideo(chatID, "/path/to/video.mp4")
```

## Пример: info-bot

В директории `examples/info-bot` находится пример консольного приложения, которое демонстрирует использование всех методов API для получения информации.

### Сборка и запуск

```bash
cd examples/info-bot
go build -o info-bot .
```

### Использование

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

# Использование кастомных URL
./info-bot -api-url https://custom-api.example.com -token YOUR_TOKEN

# Вывод в формате JSON
./info-bot -list-chats -token YOUR_TOKEN -output json
```

### Переменные окружения для info-bot

```bash
export VERBOSITY_API_TOKEN="your-bot-token"
export VERBOSITY_API_URL="https://api.verbosity.io"
export VERBOSITY_FILE_URL="https://file.verbosity.io"
export VERBOSITY_OUTPUT_MODE="text"  # text, json, json-pretty
```

## Документация API

Полная документация API доступна на сайте: https://verbosity.ru/pages/docs/bots.html

## Лицензия

MIT