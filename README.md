[![codecov](https://codecov.io/gh/ThCompiler/vk_quests_test/graph/badge.svg?token=0XHCNFY6DJ)](https://codecov.io/gh/ThCompiler/vk_quests_test)

# Тестовое задание вакансии "Программист-разработчик" в VK

## Задание:

Создать небольшой REST API сервис, засчитывающий задания для пользователя.
Код можно писать на `Go/Php`.

Хранить данные можно в любой реляционной СУБД. Сущность юзера и задания можно создать минимальные

1) Пользователь - ***User*** (**id**, **name**, **balance**);
2) Задание - ***Quest*** (**id**, **name**, **cost**).

Пользователь может выполнить каждое задание **только один раз**.

Сервис должен иметь минимально 4 апи-метода:

* метод создания пользователя,
* метод создания задания,
* метод, который сигнализирует сервису, что произошло некоторое событие, пользователь выполнил
  условие и задание можно посчитать выполненным и начислить награду пользователю (в параметрах `user_id` и `quest_id`,  
  который выполнил пользователь),
* метод, который возвращает историю выполненных пользователем заданий и его баланс.

Плюсом будут unit-тесты, какая-либо кастомизация заданий (например, с несколькими шагами до выдачи награды, 
многоразовые задания) чистая архитектура, запуск в Docker контейнере.

Проект разместить в открытом git-репозитории и предоставить ссылку, обязателен файл Readme с инструкцией как все запустить.

## Инструкция по запуску:

Задание выполнено на `Golang`. Данные хранятся в СУБД `PostgreSQL`.

В качестве дополнения был добавлен тип задач "Random", который с вероятностью 0,5 засчитывает пользователю задачу.
Также расширена сущность Задачи и в историю добавлено время выполнения задачи. Полную API можно посмотреть в swagger.yaml в папке docs. 
Или при запуске сервера на соответствующей странице.

### Запуск

#### Конфигурационный файл

В качестве примера конфигурационный файл находится в корне репозитория с название 'config.yaml'.
Его формат выглядит следующим образом:
```yaml
port: 8080 # Порт на котором запускается сервер
postgres:
  url: "host=quests-bd port=5432 user=films password=qwerty dbname=films sslmode=disable" # Строка подключения к базе Postgres
logger:  # Настройки логгера
  app_name: "vk_quests"        # Имя приложения, будет выводиться в лог
  level: 'debug'              # Минимальный уровень вывода информации в лог
  directory: './app-log/'     # Папка куда сохранять логи
  use_std_and_file: true      # Если установлено в true, то лог будет выводиться как в файл так и в stdErr
  allow_show_low_level: true  # Если установлено в true и use_std_and_file тоже true, то в stdErr будет выводиться лог всех уровней
```

#### Сборка контейнера с сервером

Перед запуском необходимо собрать Docker образ:

```cmd
sudo make build-docker
```

#### Запуск всей системы

Для запуска всей системы можно выполнить команду с выводом информации в консоль:

```cmd
sudo make run-verbose
```

Или команду которая запускает docker compose в режиме daemon:

```cmd
sudo make run
```

Система запущена. Сервер доступен на http://localhost:8080/.

Api можно посмотреть и запускать на http://localhost:8080/api/v1/swagger.
