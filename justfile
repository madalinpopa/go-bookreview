#!/usr/bin/env just --justfile

# Variables
db_path := "db.sqlite"
migrations_dir := "migrations"
css_input := "ui/assets/input.css"
css_output := "ui/static/css/output.css"
dev_port := "4000"
browser_sync_port := "4001"

# Check and install required tools
[private]
ensure-tools:
    #!/usr/bin/env sh
    if ! command -v goose >/dev/null 2>&1; then
        echo "Installing goose..."
        go install github.com/pressly/goose/v3/cmd/goose@latest
    fi
    if ! command -v air >/dev/null 2>&1; then
        echo "Installing air..."
        go install github.com/air-verse/air@latest
    fi
    if ! command -v tailwindcss >/dev/null 2>&1; then
      echo "tailwindcss cli is not installed.."
      exit 1
    fi
    if ! command -v browser-sync >/dev/null 2>&1; then
        echo "browser-sync cli is not installed.."
        exit 1
    fi

# Update Go dependencies
update:
    go get -u
    go mod tidy -v

# Run development server with live reload
dev: ensure-tools
    air & \
    tailwindcss -i {{css_input}} -o {{css_output}} --watch & \
    browser-sync start \
      --files 'ui/html/**/*' \
      --port {{browser_sync_port}} \
      --proxy 'localhost:{{dev_port}}' \
      --no-ui \
      --no-notify \
      --no-open \
      --no-ghost-mode \
      --reload-delay 1000 \
      --no-inject-changes \
      --watchEvents 'change add' \
      --middleware 'function(req, res, next) { \
        res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
        return next(); \
      }'

# Build production CSS
build:
    tailwindcss -i {{css_input}} -o {{css_output}} --minify

# Run database migrations
migrate command="up":
    goose sqlite3 {{db_path}} {{command}} --dir={{migrations_dir}}

# Create new migration
makemigrations name:
    goose sqlite3 {{db_path}} create {{name}} sql --dir={{migrations_dir}}

# Database seed
seed:
    go run ./cmd/seed/

# Run tests
test:
    go test ./internal...

# Build docker image
docker-build:
    docker build . -t coderustle/bookreview:latest --platform linux/amd64

# Run docker container
docker-run:
    docker run -d \
              --rm \
              --name bookreview \
              -v bookreview_data:/app/data \
              -v bookreview_uploads:/app/uploads \
              -p 4000:4000 \
              coderustle/bookreview:latest