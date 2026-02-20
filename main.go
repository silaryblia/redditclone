package main

import (
	"log"
	"net/http"
)

func main() {

	storage := NewStorage()

	RegisterAuthHandlers(storage)
	RegisterPostHandlers(storage)
	RegisterUserHandlers(storage)

	log.Println("Server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

/*
POST /api/register - регистрация
POST /api/login - логин
GET /api/posts/ - список всех постов
POST /api/posts/ - добавление поста - обратите внимание - есть с урлом, а есть с текстом
GET /api/posts/{CATEGORY_NAME} - список постов конкретной категории
GET /api/post/{POST_ID} - детали поста с комментами
POST /api/post/{POST_ID} - добавление коммента
DELETE /api/post/{POST_ID}/{COMMENT_ID} - удаление коммента
GET /api/post/{POST_ID}/upvote - рейтинг поста вверх
GET /api/post/{POST_ID}/downvote - рейтинг поста вни
GET /api/post/{POST_ID}/unvote - отмена ( удаление ) своего голоса в рейтинге
DELETE /api/post/{POST_ID} - удаление поста
GET /api/user/{USER_LOGIN} - получение всех постов конкретного пользователя
*/
