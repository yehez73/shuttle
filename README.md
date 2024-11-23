# REST API For Shuttle System

[WIP]
This project is an API built using Golang programming language with fiber framework.

## Installation

### From Source

#### Requirements

- [Golang](https://go.dev/doc/install)
- [MongoDB](https://www.mongodb.com/try/download/community-edition)

#### Building

##### Manually

```sh
git clone https://github.com/yehez73/shuttleapps.git
cd shuttleapps
go run .\cli\main.go
```

## Usage
Base URL = http://:8080

### Run this first

```sh
cd shuttleapps
go run .\databases\seeders\seeders.go
```

It will created dummy user for starting access

### Then

/login (user_email, password) (required all)
