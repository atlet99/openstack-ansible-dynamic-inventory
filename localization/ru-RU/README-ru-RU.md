# Динамический инвентарь для OpenStack

***Вдохновлен проектом [openstack_ansible_dynamic_inventory](https://github.com/MaksimRudakov/openstack_ansible_dynamic_inventory.git) автора [Maxim Rudakov](https://github.com/MaksimRudakov).***

Этот проект — реализация динамического инвентаря для OpenStack на Go, который позволяет организовать и фильтровать серверы по группам на основе метаданных. Подходит для автоматизации управления инфраструктурой, особенно в инструментах вроде Ansible, где удобны динамические инвентари.

## Содержание

- [Динамический инвентарь для OpenStack](#динамический-инвентарь-для-openstack)
  - [Содержание](#содержание)
  - [Функционал](#функционал)
  - [Установка](#установка)
  - [Конфигурация](#конфигурация)
    - [Файл .env](#файл-env)
  - [Использование](#использование)
  - [Переменные окружения](#переменные-окружения)
  - [Структура проекта](#структура-проекта)
  - [Локализации README](#локализации-readme)
  - [Лицензия](#лицензия)

## Функционал

* `Динамическое создание инвентаря:` автоматически получает информацию о серверах из OpenStack и структурирует её в виде динамического инвентаря.
* `Фильтрация по метаданным:` позволяет отбирать серверы на основе указанных метаданных и выбирать конкретные окружения.
* `Группировка по метаданным:` создает группы на основе метаданных, организуя серверы по ролям, окружениям и другим атрибутам.
* `Легкая интеграция:` вывод данных в формате JSON, совместимом с Ansible и аналогичными инструментами.

## Установка

1. Клонируйте репозиторий:
```shell
git clone https://github.com/atlet99/openstack-ansible-dynamic-inventory.git; cd openstack-ansible-dynamic-inventory
```
2. Установите зависимости:
```shell
go mod tidy
```

## Конфигурация

Конфигурация управляется переменными окружения с возможностью использовать файл `.env`. Эти переменные включают как учетные данные для подключения к OpenStack, так и параметры для фильтрации инвентаря. Подробнее об используемых переменных читайте в разделе [Переменные окружения](#environment-variables) ниже.

### Файл .env

Создайте `.env` файл в корне проекта, чтобы сохранить переменные окружения (также, вы можете легко скопировать `.env.template`):
```shell
OS_AUTH_URL=https://your-openstack-url:5000/v3
OS_USERNAME=your-username
OS_PASSWORD=your-password
OS_PROJECT_DOMAIN_ID="default"
OS_REGION_NAME=your-region
OS_DOMAIN_NAME="default"
ENVIRONMENT_TAG=environment
ENVIRONMENT_VALUE=test
BASE_GROUP_NAME=prod-servers
```
***Примечание:*** убедитесь, что `.env` добавлен в `.gitignore`,  чтобы не допустить случайного коммита конфиденциальной информации.

## Использование

1. Чтобы вывести весь инвентарь:
```shell
go run cmd/main.go --list
```
2. Чтобы запросить информацию о конкретном хосте (возвращает пустой JSON для совместимости):
```shell
go run cmd/main.go --host <your-hostname>
```

## Переменные окружения

Каждая переменная окружения служит для подключения к OpenStack и управления динамическим инвентарем.

```plaintext
# URL для аутентификации в OpenStack, предоставленный вашим OpenStack окружением.
# Включает в себя URL для доступа к Keystone API, используемого для аутентификации.
OS_AUTH_URL=https://your-openstack-url:5000/v3

# Имя пользователя для доступа к API OpenStack.
# Убедитесь, что у пользователя есть необходимые права для управления ресурсами.
OS_USERNAME=your-username

# Пароль указанного пользователя OpenStack.
OS_PASSWORD=your-password

# ID домена для проекта в OpenStack, обычно "default", если не указано иное.
# Этот домен связан с именем проекта, в котором управляются ресурсы.
OS_PROJECT_DOMAIN_ID="default"

# Имя региона в OpenStack, указывающее на региональный дата-центр для подключения.
# Полезно для много-региональных сетапов для обеспечения соединения с нужным дата-центром.
OS_REGION_NAME=your-region

# Имя домена для пользователя OpenStack, обычно "default", если не требуется конкретный домен.
OS_DOMAIN_NAME="default"

# Ключ метаданных, используемый для фильтрации серверов в инвентаре.
# Скрипт выберет только те серверы, у которых есть эта пара ключ-значение в метаданных.
ENVIRONMENT_TAG=environment

# Значение метаданных, которое должно совпадать для включения сервера в инвентарь.
# Вместе с ENVIRONMENT_TAG определяет окружение (например, "test", "prod"), к которому принадлежат сервера.
ENVIRONMENT_VALUE=test

# Основное имя группы, в которую попадают все выбранные серверы в инвентаре.
# Это имя организует серверы в логическую группу для удобства управления.
BASE_GROUP_NAME=prod-servers
```

## Структура проекта
```
openstack-ansible-dynamic-inventory
├── .env.template   # Файл с переменными окружения (опционально)
├── .gitignore      # Файл игнорирования временных и конфиденциальных файлов
├── LICENSE         # Лицензия проекта
├── README.md       # Документация проекта
├── cmd
│   └── main.go     # Точка входа; загружает конфигурацию и инициализирует инвентарь
├── pkg/
├── go.mod
├── go.sum
├── localization
│   ├── kz-KZ
│   │   └── README-kz-KZ.md     # Казахская версия документации
│   └── ru-RU
│       └── README-ru-RU.md     # Русская версия документации
└── pkg
    ├── inventory
    │   ├── inventory.go    # Основная логика создания инвентаря и группировки
    │   └── openstack.go    # Подключение к OpenStack и получение данных
    └── utils
        └── json.go     # Утилита для форматирования JSON
```

## Локализации README 

| Языки                  |
| -------------------------- |
| [English](../../README.md)|
| [Kazakh](localization/kz-KZ/README-kz-KZ.md) |

## Лицензия

Этот проект распространяется под лицензией MIT. См. [LICENSE](LICENSE) для получения подробной информации.