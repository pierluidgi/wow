# words-of-wisdom

## Server

Сервер принимает запросы на выдачу цитат со словами мудрости, и в зависимости от текущей нагрузки или сразу возвращает успешный ответ или выдаёт задание на подбор хеша по алгоритму Hashcash. При отправке задания сервер добавляет к нему метку времени и подпись к отправляемым данным. Если при получении ответа с решённым заданием временная метка считается истёкшей или же не совпадает подпись, сервер сразу возвращает ошибку. Если данные корректны, сервер валидирует правильность подбора хеша по заданному значению нулевых бит и проверяет, что данный ответ ранее не использовался на сервере. Если все проверки прошли успешно, сервер записывает данные из задачи в свой кеш и возвращает клиенту ответ с цитатой.

Настройки для сервера хранятся в папке *resources* в файле *server-config.yaml*, где:

- **quotes-filename** – путь до файла с цитатами
- **ddos-rate** – пороговое значение rps, при достижении которого включается механизм защиты от DDoS
- **target-bits** – число нулевых бит, которые необходимо получить при расчёте хеша
- **challenge-ttl** – сколько секунд выданное задание считается валидным
- **cache-ttl** – сколько секунд храним обработанные задания в кеше (нет смысла задавать значение больше чем challenge-ttl)
- **rate-interval** – за какой интервал в секундах считаем запросы 
- **rate-size** – количество интервалов для расчёта среднего rps

Алгоритм Hashcash в качестве выполняемой на клиенте работы был выбран как наиболее известный и хорошо описанный для таких задач. Его легко реализовать и при необходимости можно легко менять сложность подбора хеша в зависимости от нагрузки.

Запустить сервер локально можно командой: 

```sh
go run server.go -config ./resources/server-config.yaml
```

## Client

Клиент отправляет запросы на сервер, при необходимости рассчитывает хеш по полученному от сервера заданию, и при успешном ответе выводит в консоль полученную от сервера цитату.

Настройки для клиента хранятся в папке *resources* в файле *client-config.yaml*, где:

- **parallel-requests** – сколько параллельных запросов будет отправлять клиент
- **next-quote-delay-ms** – максимальная задержка в миллисекундах между запросами

Запустить клиент локально можно командой:

```sh
go run client.go -config ./resources/client-config.yaml
```

## Запуск через Docker

Для сервера и клиента созданы отдельные докер-файлы **server.Dockerfile** и **client.Dockerfile**, а также имеется возможность запустить одновременно и сервер и клиент через Docker Compose.