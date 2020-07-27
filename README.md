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

### Send notification
open postman or another programm to send request
use your USERNAME and PASSWORD from previous step
```
POST localhost:8080/user/login

{
    "email": USERNAME,
    "password": PASSWORD
}
```
then copy the Cookie header to another request
```
POST localhost:8080/notification/send

{
    "toUser": 0,
    "msg": "OH MY GOD"
}
```
if "toUser" == 0 - broadcast to all online users
else send notification to user with id == "toUser"

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

