# tech-db-forum

Тестовое задание для реализации проекта "Форумы" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).

Суть задания заключается в реализации API к базе данных проекта «Форумы» по документации к этому API.

Таким образом, на входе:

 * документация к API;

На выходе:

 * репозиторий, содержащий все необходимое для разворачивания сервиса в Docker-контейнере.

## Документация к API
Документация к API предоставлена в виде спецификации OpenAPI (https://ru.wikipedia.org/wiki/OpenAPI_%28%D1%81%D0%BF%D0%B5%D1%86%D0%B8%D1%84%D0%B8%D0%BA%D0%B0%D1%86%D0%B8%D1%8F%29): swagger.yml

Документацию можно читать как собственно в файле swagger.yml, так и через Swagger UI: https://bozaro.github.io/tech-db-forum/

## Требования к проекту
Проект должен включать в себя все необходимое для разворачивания сервиса в Docker-контейнере.

При этом:

 * файл для сборки Docker-контейнера должен называться Dockerfile и располагаться в корне репозитория;
 * реализуемое API должно быть доступно на 5000-ом порту по протоколу http.

Контейнер будет собираться из запускаться командами вида:
```
docker build -t a.navrotskiy https://github.com/bozaro/tech-db-forum-server.git
docker run -p 5000:5000 --name a.navrotskiy -t a.navrotskiy
```

## Функциональное тестирование
Корректность API будет проверяться при помощи автоматического функционального тестирования.

Методика тестирования:

 * собирается Docker-контейнер из репозитория;
 * запускается Docker-контейнер;
 * запускается скрипт на Go, который будет проводить тестирование;
 * останавливается Docker-контейнер.

Для локальной сборки Go-скрипта достаточно выполнить команду:
```
go build github.com/bozaro/tech-db-forum
```
После этого в текущем каталоге будет создан исполняемый файл `tech-db-forum`.
