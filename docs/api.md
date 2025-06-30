## API спецификация


### Базовый URL
```
https://api.maintenance-advisor.monitoring.cti.ru/v1
```

### Эндпоинты

#### 1. Получение рекомендаций
```
POST /recommendations
```

**Запрос:**
```json
{
  "metrics": {
    "cpu_usage": 95.2,
    "memory_usage": 87.5,
    "disk_io": 1200,
    "network_latency": 45
  },
  "problem_class": "performance_degradation",
  "problem_description": "Система стала медленно отвечать на запросы пользователей",
  "algorithm": "llm" // или "rag" или "mock"
}
```

**Ответ:**
```json
{
  "recommendation_id": "rec_12345",
  "recommendations": [
    {
      "priority": "high",
      "action": "Увеличьте ресурсы CPU",
      "details": "Текущая загрузка CPU составляет 95.2%, что превышает рекомендуемый порог в 80%. Рассмотрите возможность вертикального масштабирования или оптимизации рабочих нагрузок.",
      "references": [
        {
          "title": "CPU Optimization Guide",
          "url": "https://internal-docs.example.com/cpu-optimization"
        }
      ]
    },
    {
      "priority": "medium",
      "action": "Проверьте процессы с высокой загрузкой памяти",
      "details": "Использование памяти составляет 87.5%. Выполните анализ процессов для выявления утечек памяти.",
      "references": [
        {
          "title": "Memory Leak Detection",
          "url": "https://internal-docs.example.com/memory-leaks"
        }
      ]
    }
  ],
  "algorithm_used": "llm",
  "processing_time_ms": 450
}
```

#### 2. Управление классами проблем

##### Создать класс проблемы
```
POST /problem-classes
```
**Запрос:**
```json
{
  "class_name": "database_connection_issue",
  "description": "Проблемы с подключением к БД",
  "default_solution": "Проверить доступность БД и учетные данные"
}
```

##### Получить все классы
```
GET /problem-classes
```

---

#### 3. Привязка метрик к классу

##### Добавить правила для класса
```
POST /problem-classes/{class_id}/metrics-rules
```
**Запрос:**
```json
{
  "rules": [
    {
      "metric_name": "db_connections",
      "condition": "< 5",
      "weight": 0.7
    },
    {
      "metric_name": "query_timeout_errors",
      "condition": "> 10",
      "weight": 0.9
    }
  ]
}
```

---

#### 4. Управление решениями

##### Создать решение
```
POST /solutions
```
**Запрос:**
```json
{
  "title": "Восстановление подключения к БД",
  "steps": [
    "Проверить статус кластера БД",
    "Верифицировать учетные данные",
    "Перезапустить соединения"
  ],
  "documentation_url": "https://kb.example.com/db-connection"
}
```

##### Привязать решение к классу
```
POST /problem-classes/{class_id}/solutions/{solution_id}
```

---

#### 5. Обучение классификатора

##### Запустить обучение
```
POST /training
```
**Запрос:**
```json
{
  "training_data": [
    {
      "metrics": {"db_connections": 3, "query_timeout_errors": 15},
      "problem_class": "database_connection_issue"
    }
  ],
  "algorithm": "random_forest"
}
```
