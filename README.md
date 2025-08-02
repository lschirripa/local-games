# Real-time Mini-Game Platform

This project is a web-based game hosting platform where users can create and join lobbies to play mini-games in real-time. It features a SvelteKit frontend and a Go backend.

## Running the Project

### Backend

1.  Navigate to the `backend` directory:
    ```sh
    cd backend
    ```
2.  Run the application:
    ```sh
    go run cmd/api/main.go
    ```
    The backend server will start on `http://localhost:8080`.

### Frontend

1.  Navigate to the `frontend` directory:
    ```sh
    cd frontend
    ```
2.  Install dependencies:
    ```sh
    npm install
    ```
3.  Run the development server:
    ```sh
    npm run dev
    ```
    The frontend will be available at `http://localhost:5173`, but it's configured to proxy API requests to the backend on port 8080. For the full experience, you should open `http://localhost:8080` in your browser.
