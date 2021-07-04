# avito-rest-advert

[Тестовое задание](https://github.com/avito-tech/adv-backend-trainee-assignment/blob/main/README.md) для backend-стажёра в команду Advertising

## Содержание

1. [Общее описание](#ООписание)
1. [API сервиса](#API)
1. [Стек](#Стек)
1. [Установка](#Установка)
1. [Документация](#Документация)
1. [Unit-тесты](#Unit-тесты)

## Описание: <a name="ООписание"></a>

Реализован REST API для хранения и подачи объявлений. Объявления храняться в базе данных. Сервис предоставляет API, работающее поверх HTTP в формате JSON.

Сервис реализует следующие 3 метода:

- `POST /create` Метод создания объявления.
- `GET /get/:id` Метод получения конкретного объявления
- `GET /list` Метод получения списка объявлений

Реализованы следующие усложнения:

- Написаны юнит тесты для уровней приложения handler, service, repository с покрытием больше 70%
- Возможность запуска приложения командой docker-compose up;
- Архитектура сервиса описана в виде [диаграммы]() и текста
- Настроена swagger документация: есть структурированное описание методов сервиса.

## API сервиса: <a name="API"></a>

### POST /create

- Метод создания объявления
- Принимает:
  - name (название,не больше 200 символов, обязательное поле, type - string)
  - description (описание,не больше 1000 символов, type - string)
  - price (цена, положительное число, type - float)
  - pictures (несколько ссылок на фотографии, не больше 3 ссылок на фото, type - string).  
    Фотографии идут без пробела, разделяются запятыми.  
    Валидное поле "avito/files/ad1,avito/files/ad2,avito/files/ad3".  
    Невалидное поле - "avito/files/ad1,avito/files/ad2,avito/files/ad3,avito/files/ad4".

```
{
	"name": "name-test",
  "description": "desc-test",
  "price": 1000,
  "pictures": "avito/files/ad1,avito/files/ad2,avito/files/ad3",
}
```

- Возвращает JSON c id созданного объявления или сообщение об ошибки, а также код результата

- Возвращает в случае успеха статус код 200 и JSON c id пользователя

Запрос:

```
curl --request POST "http://localhost:8080/create"  \
--header 'Content-Type: application/json'  \
--data-raw '{
	  "name": "name-test",
		"description": "desc-test",
		"price": 1000,
		"pictures": "avito/files/ad1,avito/files/ad2,avito/files/ad3"
}'
```

Ответ:

```
{
  {"id":1}
}
```

- Возвращает в случае неправильных входных данных статус код 400 и JSON c описанием ошибки

Запрос:

```
curl --request POST "http://localhost:8080/create"  \
--header 'Content-Type: application/json'  \
--data-raw '{
		"description": "desc-test",
		"price": 1000,
		"pictures": "avito/files/ad1,avito/files/ad2,avito/files/ad3"
}'
```

Ответ:

```
{
  {"error":"invalid input body"}
}
```

- Возвращает в случае ошибки не сервере статус код 500 и JSON c описанием ошибки

Запрос:

```
curl --request POST "http://localhost:8080/create"  \
--header 'Content-Type: application/json'  \
--data-raw '{
		"description": "desc-test",
		"price": 1000,
		"pictures": "avito/files/ad1,avito/files/ad2,avito/files/ad3"
}'
```

Ответ:

```
{
  {"error":"server error"}
}
```

## Стек технологий: <a name="Стек"></a>

golang, gin-gonic, postgres

## Установка: <a name="Установка"></a>

- Склонируйте проект с реппозитория GitHub
  ```
  git clone https://github.com/paramonies/avito-rest-advert.git
  ```
- Перейдите в директорию ./avito-rest-advert
  ```
  cd ./avito-rest-advert
  ```
- Запустите docker-compose
  ```
  docker-compose up
  ```

## Документация: <a name="Документация"></a>

- Документация к API http://localhost:8080/swagger/index.html

## Unit-тесты: <a name="Unit-тесты"></a>

Запуск юнит-тестов:

```
docker-compose run --rm make test
```
