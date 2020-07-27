### Require
- docker
- golang

### Start
```
docker-compose up -d
go run main.go
```

### Start Client
```
// open new terminal
cd client
go run clent.go -user=USERNAME -pass=PASSWORD
```

### Client commands
```
send [USER_ID] [MSG]
if USER_ID == 0 - broadcast to all online users

users
get count of online users

exit
```


## API

#### Auth
```
POST

/user/new - create new user
/user/login - login and get cookie

{
    "email": string,
    "password": string
}
```
#### WS connection
```
GET

/notification/subscribe - get ws connet with server
```
#### Send notification
```
POST

/notification/send

{
    "toUser": uint,
    "msg": string
}

toUser == 0 - broadcast to all online users
toUser > 0 - send notification to user
```
#### Online users count
```
GET

/notification/users

{
    "count": uint32,
}
```
