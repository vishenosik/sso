# grafana-loki-promtail

## Разворачивание на локальной машине

1. Настроить .env файл - натравить `PROMTAIL_LOCAL_LOGS` на папку с исследуемыми логами (абсолютный путь). Пример есть в example.env. Остальные части конфига можно целиком копировать из example.env.
2. Создать папку log, в которую закидываются логи на анализ. Либо можно миновать этот подход просто натравив `PROMTAIL_LOCAL_LOGS` на целевую папку. Оба подхода имеют место быть.
3. Запустить docker-compose через утилиту task или любым другим удобным способом.
4. Проверить что все запустилось можно перейдя по localhost:9080 (promtail), localhost:3100 (loki), localhost:3000 (grafana).
5. Поздравляю! Вы прекрасны.

Запуск через докер:

```bash
docker compose -f docker-compose.yaml up -d --build
```

[.env example](docs/.env)

```text
# Примерно так выглядит мое файловое дерево при работе
.
├── log # Папка, в которую я закидываю .log файлы для анализа
├── loki/
│   └── config.yaml
├── promtail/
│   └── config.yaml
├── .env
├── Taskfile.yaml # Файл для работы с утилитой task
└── docker-compose.yaml
```
