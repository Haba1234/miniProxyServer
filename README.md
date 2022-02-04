Небольшой тестовый HTTP-сервер для проксирования HTTP-запросов к сторонним сервисам.

## Описание
Ожидает HTTP-запрос от клиента. В теле запроса сообщение в формате json вида:
```json
{
    "method": "GET",
    "url": "http://google.com",
    "headers": {
        "Authentication": "Basic bG9naW46cGFzc3dvcmQ=",
        ....
    }
}
```
Результат возвращается клиенту в формате json вида:
```json
{
  "data": {
    "id": "<generated unique id>",
    "status": <HTTP status of 3rd-party service response>,
    "headers": <headers array from 3rd-party service response>,
    "length": <content length of 3rd-party service response>
  },
  "error": {
    "message": "error text"
  }
}
```

### Сборка и запуск
```
git clone https://github.com/Haba1234/miniProxyServer.git
cd miniProxyServer
make build
make run
```
