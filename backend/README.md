# Local Games Backend

A Go REST API server for the local games platform.

## API Endpoints

### Games
- `GET /api/games` - Get all games
- `POST /api/games` - Create a new game
- `GET /api/games/{id}` - Get a specific game

### Players
- `GET /api/players` - Get all players
- `POST /api/players` - Create a new player

## Running the Server

```bash
# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The server will start on `http://localhost:8080`.

## CORS Configuration

The server is configured to allow requests from `http://localhost:5173` (SvelteKit dev server).