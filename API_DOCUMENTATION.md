# Дополнительная документация API go-verbosity

## Обзор

Этот файл содержит дополнительную документацию API, которая не включена в основной README.md. Здесь описаны расширенные возможности и экспериментальные функции.


## Удаление сообщений

### DELETE /msg/post/{chat_id}/{post_no}


## Обновление сообщений

### PUT /msg/post/{chat_id}/{post_no}

Обновляет ранее отправленное сообщение в чате.

#### Параметры запроса

- `chat_id` (integer, required) - ID чата
- `post_no` (integer, required) - номер поста для обновления

#### Тело запроса

```json
{
  "uuid": "<string> | null",           // string|null - UUID сообщения
  "e2e": <boolean>,                    // boolean - флаг E2E шифрования
  "text": "<string>",                  // string, required - текст сообщения
  "reply_no": "<integer> | null",      // integer|null - номер поста для ответа
  "quote": "<string 0...32000> | null",// string|null - цитата (макс. 32000 символов)
  "attachments": ["<string>"] | null   // array|null - массив GUID вложений
}
```

#### Ответ

```json
{
  "uuid": "<string>",          // string, required - UUID обновленного сообщения
  "chat_id": "<ChatId>",       // integer, required - ID чата
  "post_no": "<PostNo>",       // integer, required - номер поста
  "ver": "<integer> | null"    // integer|null, required - версия сообщения
}
```

#### Пример использования

```go
// Обновление текста сообщения
updateRequest := &UpdateMessageRequest{
    Text: "Обновленный текст сообщения",
    E2E: false,
}

response, err := client.UpdateMessage(123, 456, updateRequest)
if err != nil {
    log.Fatalf("Ошибка обновления сообщения: %v", err)
}

fmt.Printf("Сообщение обновлено: UUID=%s, версия=%d\n", response.UUID, response.Version)
```
