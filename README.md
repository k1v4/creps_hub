# Creps hub

creps_hub - это веб-форум или же веб-блог, который будет создан для того, чтобы люди, которые интресуются миром обуви или же являются его активными ценителями (сникерхэдами) могли узнавать что-то новое 
или же делиться этим с миром.

## Функциональность

- Данная веб-сайт будет позволять пользователям регистрироваться на платформе.
- В личном кабинете пользователи смогут создавать свою частную коллекцию обуви. Для этого надо будет загрузить фото пары и её название.
- Также в личном кабинете будет возможность писать свои статьи о мире обуви.
- Одна из главных функций - это возможность чтения и выделения статей (можно поставить отметку нравиться статье).

## Техническая часть

- Приложение разбито на сервисы, которые взаимодействуют, используя или http, или gRPC. 
- Работает это всё с СУБД PostgreSQL. 
- В данном приложении используется несколько адаптированная чистая архитектура. 
- Также для проекта написан docker-compose файл.

