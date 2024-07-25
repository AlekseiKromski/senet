# API documentation

## Auth

### Authorize in application
After successful auth, server will return you `jwt` token given to one hour
```http request
POST /api/auth
Content-Type: application/json

{
  "username": "{{username}}",
  "password": "{{password}}"
}

>>> Response
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Set-Cookie: user=${json user here}; Path=/; Domain=localhost; Max-Age=3600; Secure
Date: Thu, 25 Jul 2024 18:09:03 GMT
Content-Length: 200

{
  "token": "${jwt token here..}"
}
```

## Chat

### Get all chats
```http request
GET /api/chat/get
Authorization: Bearer {{token}}
Content-Type: application/json

>>> Response
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 25 Jul 2024 18:15:51 GMT
Transfer-Encoding: chunked

[
  {
    "id": "8bde6125-d0c9-46ca-bc99-3592cb0ed547",
    "name": "user1 <-> user2",
    "chat_type": "PRIVATE",
    "security_level": "SERVER_PRIVATE_KEY",
    "last_message": {
      "id": "9ca8ccdc-e3db-4bae-a02e-d94d5a9be44a",
      "chat_id": "8bde6125-d0c9-46ca-bc99-3592cb0ed547",
      "sender_id": "c8658ba2-e7d6-4246-ac8b-84fdf5caf68a",
      "message": "some message",
      "created_at": "2024-07-20T22:22:29.875042Z",
      "updated_at": "2024-07-20T22:22:29.875042Z",
      "deleted_at": null
    },
    "created_at": "2024-07-06T16:10:44.916727Z",
    "updated_at": "2024-07-06T16:10:44.916727Z",
    "deleted_at": null,
    "users": [
      {
        "id": "c8658ba2-e7d6-4246-ac8b-84fdf5caf68a",
        "username": "user1",
        "first_name": "user1",
        "second_name": "user1",
        "image": "default_user.png",
        "created_at": "2024-07-06T11:00:16.732985Z",
        "updated_at": "2024-07-06T11:00:16.732985Z",
        "deleted_at": null
      },
      {
        "id": "1dccc1d6-e95d-44d0-bccc-5e73856f989a",
        "username": "user2",
        "first_name": "user2",
        "second_name": "user2",
        "image": "default_user.png",
        "created_at": "2024-07-06T12:07:49.729041Z",
        "updated_at": "2024-07-06T12:07:49.729041Z",
        "deleted_at": null
      }
    ]
  }
]
```

