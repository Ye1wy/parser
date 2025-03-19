# Scraper for web site
#### This program, at the moment, can only scrap samokat.ru
## Installation
### 1. Install chromedriver latest version
### 2. For scrap needed .env file like example.env on root directory of this project
```
cp example.env .env
```
### 3. For proxy rotation needed create proxy_list.json on *config* directory
```
cp config/example.proxy_list.json config/proxy_list.json
```
### 4. After creation, fill out the form
### Example proxy_list:
```
[
    {
        "host": "your.proxy.address",
        "port": your.proxy.port,
        "username": "login_on_proxy_server",
        "password": "password_on_proxy_server",
        "type": "your_proxy_protocol_type" #sock5, http, httpss
    },
    {
        "host": "your.proxy.address",
        "port": your.proxy.port,
        "username": "login_on_proxy_server",
        "password": "password_on_proxy_server",
        "type": "your_proxy_protocol_type" #sock5, http, httpss
    },
    ...
]
```

## Run
```
go run cmd/main.go
```