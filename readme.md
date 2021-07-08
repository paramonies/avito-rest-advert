# avito-rest-advert

[Тестовое задание](https://github.com/avito-tech/adv-backend-trainee-assignment/blob/main/README.md) для backend-стажёра в команду Advertising

## Содержание

1. [Общее описание](#ООписание)
1. [Стек технологий](#Стек)
1. [Установка](#Установка)
1. [Документация](#Документация)
1. [Архитектура](#Архитектура)
3. [Тестирование](#тесты)

## Описание: <a name="ООписание"></a>

Реализован REST API для хранения и подачи объявлений. Объявления храняться в базе данных. Сервис предоставляет API, работающее поверх HTTP в формате JSON.

Сервис реализует следующие 3 метода:

- `POST /create` Метод создания объявления, поля создаваемого объявления передаются в теле запроса в json формате и являются обязательными:
  - name: название, type - string, валидация: не больше 200 символов
  - description: описание объявления, type - string, валидация: не больше 1000 символов
  - price: цена, type - int, валидация: положительное число
  - pictures: ссылки на фотографии, type - string, валидация: не больше 3 ссылок на фото(ссылки на фото разделяются запятыми,идут без пробелов)  
    
    
- `GET /get/:id?fields=description,pictures` Метод получения конкретного объявления
  - id - идентификатор объявление, обязательный параметр
  - fields - список дополнительных полей в ответе, принимает одно из значении {"description", "pictures", "description,pictures", ""}, по умолчанию ""

- `GET /list?page=2&order_by=createdat_desc` Метод получения списка объявлений
  - page - номер страницы, 1 по умолчанию
  - order_by - сортировка по цене (возрастание/убывание) или по дате создания (возрастание/убывание), по умолчанию "createdat_desc", 
    принимает одно из значений {"price_desc", "price_asc", "createdat_desc", "createdat_asc"}

Реализованы следующие усложнения:

- Написаны юнит тесты для уровней приложения handler, service, repository с покрытием больше 70%
- Возможность запуска приложения командой docker-compose up;
- Архитектура сервиса описана в виде [диаграммы](/docs/media/service-arc.png) и текста
- Настроена swagger документация: есть структурированное описание методов сервиса.

## Стек технологий: <a name="Стек"></a>

Golang, Gin-gonic, Postgres

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
- После установки сервис доступен по http на 8080 порту: http://localhost:8080/create
  

## Документация: <a name="Документация"></a>

Документация к API http://localhost:8080/swagger/index.html
![docs](/docs/media/swagger_doc.png)

## Архитектура сервиса: <a name="Архитектура"></a>
![arc](/docs/media/service-arc.png)

## Тестирование: <a name="тесты"></a>

Запуск юнит-тестов:

```
docker exec -it api-server make test
```
Рассчет покрытия:

```
docker exec -it api-server make cover
```

