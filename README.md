# Database Systems Final Project

## Backend of a basic online shop

### Install dependencies
```shell
go get ./...
```

### Configure settings with `.env`
Create `.env`
```shell
touch .env
```

(Optional) Create sqlite database
```shell
touch test.db
```

Set database configs in `.env`
```
DB_DRIVER=sqlite
DB_SOURCE=./test.db
```

Set JWT secret in `.env`
```
IMGUR_API_CLIENT_ID="deff1952d59f883ece260e8683fed21ab0ad9a53323eca4f"
```

[Imgur API docs](https://apidocs.imgur.com)
[Register application to get imgur API client ID](https://api.imgur.com/oauth2/addclient)


Set Imgur client ID in `.env`
```
IMGUR_CLIENT_ID="008B8008B848B8"
```


### Start server
```shell
go run ./main.go
```
