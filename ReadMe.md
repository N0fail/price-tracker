# Трэкер цен

Трэкер цен предполагает запись истории цен на какой то товар и просмотр этой истории

В качестве интерфейса реализованы: 
- Телеграм бо
- GRPC
- REST

## Конфигурация БД
Для запуска Postresql понадобятся 
```
docker https://docs.docker.com/engine/install/ 
docker-compose https://docs.docker.com/compose/install/
```
При возникновении проблем с установкой docker-compose на debian-based linux
```
https://stackoverflow.com/a/49839172
```
БД поднимается с помощью
```
make up_db
```

## Миграции
Для миграций используется утилита goose, она компилируется в папку bin при вызове
```
make .deps
```
Все миграции описаны в папке migrations

Для запуска миграций можно использовать скрипт
```
./migrate.sh
```

Для добавления свей миграции можно использовать
```
make migration NAME=migration_name
```
В результате будет создан файл в папке migrations, в котором нужно описать новую миграцию

## Запуск проекта

1. скачать зависимости
    ```
    go get ./...
    ```
2. задать переменную окружения PriceTrackerApiKey с API ключом бота
   ```
    export PriceTrackerApiKey=your_key
    ```
   1. Ключ можно получить у https://t.me/botfather прислав ему сообщение /newbot
3. запустить 
    ```
    make run_server
    ```

## Доступные комманды телеграм бота

Доступные комманды описаны в internal/pkg/bot/command/ 

- /help выводит краткую информацию о доступных коммандах

- /add <product_code>_<product_name> добавляет продукт для отслеживания  
product_code - строка-ключ продукта, при попытке добавления продуктов с одинаковым product_code возникнет ошибка, о которой бот сообщит  
product_name - строка "читаемое" имя продукта, должна быть не короче, чем config.MinNameLength(4)  
пример: /add qwer1234_Guitar - добавит продукт с кодом qwer1234 и именем Guitar  

- /add_price <product_code>_<date>_<price> добавляет запись о цене продукта  
product_code - строка-ключ продукта, при попытке записи для неотслеживаемого продукта возникнет ошибка, о которой бот сообщит  
date - дата, в которую была записана цена, формат указан в config.DateFormat(2 Jan 2006 15:04)  
price - цена продукта, является целым неотрицательным числом  
пример: /add_price qwer1234_2 Jan 2006 15:04_3000 - добавит для продукта qwer1234 запись о цене 3000 второго января 2006 в 15:04  

- /price_history <product_code> выводит историю цен продукта в хронологическом порядке  
product_code - строка-ключ продукта, при попытке вывести историю для неотслеживаемого продукта возникнет ошибка, о которой бот сообщит  
пример: /price_history qwer1234  

- /list выводит список всех отслеживаемых продуктов с последней ценой(в хронологическом смысле)  

- /delete <product_code> удаляет продукт из отслеживания  
product_code - строка-ключ продукта, при попытке удаления неотслеживаемого продукта возникнет ошибка, о которой бот сообщит  
вместе с продуктом удаляется и его история цен  
пример: /delete qwer1234

## Доступные GRPC процедуры
Доступные GRPC функции описаны в файле api/api.proto

- ProductCreate - добавляет продукт для отслеживания, аналог /add комманды бота

- PriceTimeStampAdd - добавляет запись о цене продукта, аналог /add_price комманды бота

- PriceHistory выводит историю цен продукта в хронологическом порядке, аналог /price_history комманды бота

- ProductList выводит список всех отслеживаемых продуктов с последней ценой(в хронологическом смысле), аналог /list комманды бота

- ProductDelete удаляет продукт из отслеживания, аналог /delete комманды бота

## Пример GRPC клиента

В проект добавлен пример GRPC клиента в client/client.go

Для запуска можно использовать 
```
make run_client
```

Генерацию кода по протофайлу можно выполнить с помощью 
```
make gen
```

Возможно понадобится перекомпилировать модули протобафа, для этого можно выполнить
```
make .deps
```

## SWAGGER

Для проекта также запускается сваггер по адресу http://localhost:8080/swagger/

## Примеры запросов к REST

- Получение списка продуктов с последней ценой
```
curl localhost:8080/v1/product
```
- Добавление продукта
```
curl -X POST localhost:8080/v1/product -d '{"code":"qwer1234", "name":"Guitar"}'
```
- Добавление цены продукта
```
curl -X PUT localhost:8080/v1/price -d '{"code":"1", "price":123, "ts":1234567}'
```
- Получение истории цен продукта
```
curl localhost:8080/v1/price?code=1
```
- Удаление продукта
```
curl -X DELETE localhost:8080/v1/product?code=1
```
