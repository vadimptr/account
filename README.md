# Инструкция пользователя

## Запуск приложения

1. Клонируем репозиторий
```
git@github.com:vadimptr/account-sync.git
```

2. Собираем проект
```
go build
```

3. Запускаем приложение
```
./account-sync.exe
```

Должны увидеть
```
Connection to amqp... [success]
Connecting to postgres... [success]
Press Ctrl+C to exit
```

## Настройка приложения

Приложение по умолчанию подключается к

RabbitMQ
```
amqp://zwdkijew:FgA3Ilyct6--rfo1zfXIMFRlNpO6OC5j@whale.rmq.cloudamqp.com/zwdkijew
```

Postgres
```
postgres://mbympgbxovcaec:5c254085dca2140af8553b3c941abe44b47f7569e63d782c8db52b3e40970205@ec2-54-228-246-214.eu-west-1.compute.amazonaws.com:5432/d5idksj6ro3iuo
```

Это поведение можно изменить задав соответственно две переменные среды: CLOUDAMQP_URL и DATABASE_URL

## Проведение тестов

Описание с учетом дефолтных внешних сервисов

1. Запускаем в не менее чем двух экземплярах наш воркер

2. Открываем менедж консоль RabbitMQ
```
https://whale.rmq.cloudamqp.com/#/exchanges/zwdkijew/input_balance_change_exchange
login: zwdkijew
password: FgA3Ilyct6--rfo1zfXIMFRlNpO6OC5j
```

3. В админке находим блок спойлера "Publish message"

Routing key: *input_balance_change_queue*

Payload:
```
{
    "single_user": {
        "user": "oleg",
        "amount": 100
    }
}
```

4. Жмем кнопку "Publish message"

5. Один из воркеров должен обработать сообщение

```
Recieved new message.
   Message valid. {oleg 100}
   Now ballance: {oleg 200}
   Message processed.
   Pushed response: success
```

## Структура сообщения

1. Сообщение об изменении баланса одного пользователя

Структуру сообщения можно было сделать множеством способов, я выбрал такой вариант, тк он очевидно описывает то что мы хотим сделать (хотя название single_user мне не очень нравится)
```
{
    "single_user": {
        "user": "<имя пользователя>",
        "amount": <сумма снятия или зачисления>
    }
}
```
Поля user и amount обязательные.
Минимальная длинна user 1 символ.
amount не может быть равен 0

2. Сообщение трансфера денег
```
{
    "transfer": {
        "from_user": "<имя пользователя>",
        "to_user": "<имя пользователя>",
        "amount": <положительная сумма>
    }
}
```
Поля from_user, to_user и amount обязательные.
Минимальная длинна from_user и to_user 1 символ.
amount исключительно положительное число

## Какие ошибки можно получить

1. Ошибка обновления пользователя (то ради чего все затевалось)
```
Recieved new message.
   Message valid. {test 100}

(pq: could not serialize access due to concurrent update)
[2019-07-14 21:54:17]
   Pushed response: error pq: could not serialize access due to concurrent update
```

2. Пришло кривое сообщение (не валидный json)
```
Recieved new message.
   Pushed response: error unexpected EOF
```

3. Недостаточно средств при списании у одного пользователя
```
Recieved new message.
   Message valid. {test -10000}
   Pushed response: error not enough amount. exist 600 try to subsctruct -10000
```

4. Не прошло валидацию входное сообщение
```
Recieved new message.
   Pushed response: error transfer.amount: Must be greater than or equal to 1; transfer.from_user: String length must be greater than or equal to 1;
```

## Нюансы и особенности работы

1. Валидация пришедшего сообщения осуществляется json схемой которая генерируется из объектной модели

2. Если пользователю начисляется некая сумма причем пользователя нет в базе, то пользователь создается автоматически (для упрощения тестирования)

3. Реализовал только один способ из нижеперечисленных, но код позволяет реализовать и другие способы.

## Способы безопасных конкурентных изменений баланса пользователя

1. Способ на основе транзакции с уровнем изоляции serializable. 

2. Способ "уход от read-write-update" когда мы конструктцию UPDATE заменияем на относительное изменение 
```
UPDATE accounts SET amount = amount + 100 WHERE name = "test"
```
нам не подходит тк приводит к отрицательным балансам.

3. Способ FOR UPDATE (подходит) и он также прост в реализации.

4. Способ optimistic concurency control, когда мы проставляем версию каждой строке данных. При каждом обновлении версию увеличиваем. Также каждое обновление ищет запись с учетом известной на данный момент версии.

## Замечания и возможные улучшения

1. Можно было сделать внешний конфигурационный файл, вместо использования переменных среды
2. Ошибки можно было бы затюнить более человекочитаемыми (но пришлось бы расширять скажем json schema validator)
3. Оградиться от rabbitmq и postgres интерфейсами, чтобы не зависеть ни от очереди ни от базы данных
