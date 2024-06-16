## Auction API

# Installation
1. Clone the repository git clone git@github.com:ypeskov/go_auction_api.git [folder]
2. cd [folder]
3. Run bash build-and-run.sh (macOS/Linux). Windows is not supported.
4. The server will be running on localhost:3000

# Run Linter
```bash
golangci-lint run
```

# Run Tests
```bash
make test && make coverfunc
```

# Packages for migration
```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

# Create new migration
```bash
cd db
migrate create -ext sql -dir ./migration [migration_name]
```

# Run migrations
```bash
cd db
./migrate.sh up
```
