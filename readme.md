# Api proxy

Проксирует запросы на удаленные сервера

## Запуск:
```bash
apiproxy /path/to/config.json
```

## Пример конфига:

```json
{
  "Port": ":8000"
}
```
### Параметры конфига:
* Port - порт для прослушивания

## Использование

http://localhost:8000/?url=http://ya.ru
