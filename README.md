# Scattergories Backend

Scattergories Backend is a web application that provides the backend services for a Scattergories game. The application is built using Go and the Gin framework, and it includes functionalities for managing users, game rooms, and game configurations.

## Project Structure

The project is organized into the following directories:

- **cmd/scattergories**: Contains the entry point of the application.
- **config**: Contains configuration files and scripts.
- **internal**: Contains the core application logic divided into several subdirectories:
  - **models**: Defines the data models used in the application.
  - **client**: Contains client-specific logic, including controllers, routes, and WebSocket handling.
  - **services**: Contains business logic and services used by the controllers.
- **pkg**: Contains utility packages and custom validators.
- **test**: Contains test files for unit and integration tests.
- **docs**: Documentation files.

## Setup and Installation

### Steps

1. **Clone the repository**:

    ```sh
    git clone https://github.com/your-username/scattergories-backend.git
    cd scattergories-backend
    ```

2. **Set up the environment variables**:

    Create a `.env` file in the root directory and add the necessary environment variables. Example:

    ```env
    DB_HOST=localhost
    DB_USER=your_db_user
    DB_NAME=your_db_name
    DB_SSLMODE=disable
    DB_PASSWORD=your_db_password
    ```

3. **Install dependencies**:

    ```sh
    go mod download
    ```

4. **Run the application**:

    ```sh
    go run ./cmd/scattergories/main.go
    ```

## Usage

### API Endpoints

The application exposes several API endpoints to manage users, game rooms, and games. Here are some of the key endpoints:

- **User Management**:
  - `GET /users`: Retrieve all users.
  - `GET /users/:id`: Retrieve a single user by ID.
  - `POST /users`: Create a new user.
  - `PUT /users/:id`: Update an existing user.
  - `DELETE /users/:id`: Delete a user by ID.

- **Game Room Management**:
  - `GET /game-rooms`: Retrieve all game rooms.
  - `GET /game-rooms/:room_id`: Retrieve a single game room by ID.
  - `POST /game-rooms`: Create a new game room.
  - `DELETE /game-rooms/:room_id`: Delete a game room by ID.
  - `PUT /game-rooms/:room_id/update-host`: Update the host of a game room.
  - `POST /game-rooms/:room_id/join`: Join a game room.
  - `POST /game-rooms/:room_id/leave`: Leave a game room.

### WebSocket Endpoints

The application also supports WebSocket connections for real-time communication. Here are some key WebSocket endpoints:

- `GET /ws/:room_id`: Handle WebSocket connections for a game room.

## Testing

The project includes unit and integration tests located in the `test` directory. To run the tests, use the following command:

```sh
go test ./...
```