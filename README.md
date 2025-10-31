# necronomnomicon

This project arose from a need to store our family recipes in a centralized location, and you didn't think I was going to use someone *else's* recipe keeper did you.

A recipe keeper web application for preserving culinary knowledge that should not be forgotten.

## Overview

Necronomnomicon is a web-based application built with Go. It allows users to create, search, and maintain their recipes online without any nasty adds. Each recipe can be created, edited, deleted, and viewed by those who possess the proper credentials.

## Features

- User authentication via JWT tokens
- Create, read, update, and delete recipes
- Search recipes by name
- Infinite scroll recipe browsing
- Secure recipe ownership and editing permissions

## Tech Stack

- **Backend**: Go with Echo web framework
- **Database**: libsql (Turso)
- **Authentication**: cookie base JWT auth
- **Templating**: Go templates
- **Frontend**: htmx
- **Migrations**: goose
- **Query Generation**: sqlc

## Prerequisites

Before setting up the application, ensure you have the following installed:

- Go 1.24 or higher
- [goose](https://github.com/pressly/goose) - for database migrations
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```
- [sqlc](https://sqlc.dev/) (optional) - if you need to regenerate database queries
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/seandisero/necronomnomicon.git
cd necronomnomicon
```

2. Create a `.env` file with the required configuration:
```
PORT=8080
DB_URL=your_turso_database_url
DB_TOKEN=your_turso_auth_token
```

3. Install dependencies:
```bash
go mod download
```

4. Run database migrations to summon the schema into existence:
```bash
./scripts/migrateup.sh
```

5. Run the application:
```bash
go run cmd/main.go
```

The application will awaken on the port specified in your `.env` file.

## Development

For development with hot reloading, you can use [Air](https://github.com/air-verse/air):
```bash
go install github.com/air-verse/air@latest
air
```

## Database

The application uses libsql with embedded replicas, synchronizing with a remote Turso database every 30 minutes. Local changes are written to the embedded replica, ensuring the knowledge persists even when the connection to the ancient servers is severed.

Database migrations are managed with goose and can be found in `sql/schema/`. To roll back migrations:
```bash
./scripts/migratedown.sh
```

## License

This project contains no eldritch horrors, only recipes.
