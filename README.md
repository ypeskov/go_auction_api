## Auction API

# Installation
1. Clone the repository git clone git@github.com:ypeskov/go_auction_api.git [folder]
2. cd [folder]
3. Run bash build-and-run.sh (MacOS/Linux). Windows is not supported.
4. The server will be running on localhost:3000

# Run Linter
golangci-lint run

# Run Tests
make test && make coverfunc

Directory Tree

.\
├── Dockerfile\
├── README.md\
├── build-and-run.sh\
├── cmd\
│   └── main.go\
├── db\
│   ├── migrate.sh\
│   ├── migrations\
│   │   ├── 000001_create_users_table.down.sql\
│   │   ├── 000001_create_users_table.up.sql\
│   │   ├── 000002_add_passwordhash.down.sql\
│   │   └── 000002_add_passwordhash.up.sql\
│   └── tmp\
│       └── build-errors.log\
├── docker-compose.yml\
├── docs\
│   ├── docs.go\
│   ├── swagger.json\
│   └── swagger.yaml\
├── go.mod\
├── go.sum\
├── internal\
│   ├── config\
│   │   └── config.go\
│   ├── database\
│   │   └── database.go\
│   ├── errors\
│   │   └── errors.go\
│   └── log\
│       └── logger.go\
├── proj.html\
├── repository\
│   ├── models\
│   │   ├── item.go\
│   │   └── user.go\
│   └── repositories\
│       ├── item-repository.go\
│       └── user-repository.go\
├── server\
│   ├── middleware\
│   │   └── auth-middleware.go\
│   ├── routes\
│   │   ├── items-routes.go\
│   │   ├── routes.go\
│   │   └── users-routes.go\
│   └── server.go\
├── services\
│   ├── items-service.go\
│   └── users-service.go\
└── tmp\
├── build-errors.log\
└── main