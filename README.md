# gophkeeper

# Спецификация

Спецификация проекта находится в файле [SPECIFICATION.md](SPECIFICATION.md)

# Текущая реализация

Т.к. времени остается мало, будет сделан"базовый" функционал. Различные "плюшки", например паджинация при подтягивании
данных с сервера, шифрованиеданных, хранящихся на сервере, и т.д. не будет.

Моменты, которые хотелось сделать, но не были сделаны, будут постепенно добавляться в [TODO](TODO.md) файл.

Т.к. потратил время на прикручивание grpc-gateway, то Reverse Proxy c REST API будет в каком виде. Сейчас оно есть для
регистрации или логина пользователей. Хочеться сделать и для работы с записами (пара логин/пароль, card, text data,
binary data).

# TODO

TODO лист находится в файле [TODO.md](TODO.md)