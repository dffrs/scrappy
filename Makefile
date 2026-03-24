.DEFAULT_GOAL := build

.PHONY: clean fmt vet build

SQLITE_FLAGS=CGO_ENABLED=1 CGO_CFLAGS="-DSQLITE_ENABLE_FTS5" CGO_LDFLAGS="-lm"

### ───────── DATABASE ─────────

upDB:
	$(SQLITE_FLAGS) go run ./cmd/migrate/main.go up

downDB:
	$(SQLITE_FLAGS) go run ./cmd/migrate/main.go down

resetDB: downDB upDB

clean:
				go clean
fmt: clean
				go fmt ./...
vet: fmt
				go vet ./...
build: vet
				go build -o bin/scrappy cmd/scrappy.go
