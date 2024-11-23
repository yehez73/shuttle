# REST API For Shuttle System

[WIP]
This project is an REST API built using Golang programming language and Fiber framework.

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

##### Run it with automatic recompilation when any Project files are changed
```sh
cd shuttleapps
air init
```
It will create a toml file, open it and change the cmd to
```
cmd = "go build -o ./tmp/main.exe ./cli/main.go"
```
Then type this in command prompt
```
air
```

## Usage
Base URL = http://:8080

### Run this first

```sh
cd shuttleapps
go run .\databases\seeders\seeders.go
```

It will create a dummy user for starting access

### Then

/login (user_email, password) (required all)
