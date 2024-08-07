# Scattergories Backend

Scattergories Backend is a web application that provides the backend services for a Scattergories game. The application is built using Go and the Gin framework, and it includes functionalities for managing users, game rooms, and game configurations. PostgreSQL is used as the database for storing data, and WebSocket is used for real-time game communication.

## Project Structure

The project is organized into the following directories:

- **cmd/scattergories**: Contains the entry point of the application.
- **config**: Contains configuration files and scripts.
- **internal**: Contains the core application logic divided into several subdirectories:
  - **api**: Contains client-specific logic, including controllers, routes, and WebSocket handling.
  - **domain**: Defines the data models used in the application.
  - **repositories**: Contains the data access logic and repository patterns for interacting with the database.
  - **services**: Contains business logic and services used by the controllers.
- **pkg**: Contains utility packages and custom validators.
- **test**: Contains test files for unit and integration tests.

## Usage

### API Endpoints

The application exposes several API endpoints to manage users, game rooms, and games. Here are some of the key endpoints:

- **User Management**:  
  - `POST /guests`: Join as a guest user.
  - `GET /users`: Retrieve all users.
  - `GET /users/:id`: Retrieve a single user by ID.
  - `POST /users`: Create a new account.
  - `DELETE /users/:id`: Delete a user by ID.
  - `POST /login`: Login.
  - `POST /logout`: Logout.
  - `POST /refresh-token`: Refresh token.


- **Game Room Management**:
  - `GET /game-rooms`: Retrieve all game rooms.
  - `GET /game-rooms/:room_id`: Retrieve a single game room by ID.
  - `POST /game-rooms`: Create a new game room.
  - `DELETE /game-rooms/:room_id`: Delete a game room by ID.
  - `PUT /game-rooms/:room_id/join`: Join a game room.
  - `PUT /game-rooms/:room_id/leave`: Leave a game room.

### WebSocket Endpoints

The application also supports WebSocket connections for real-time communication:

- `/ws/:room_id`: Establish a WebSocket connection for a game room. This connection allows clients to send and receive real-time updates and requests related to the game. The server will broadcast responses to all connected clients in the room.
    - `start_game_request`: Initiate the start of a game.
    - `end_game_request`: End the current game.
    - `submit_answer_request`: Submit an answer for a game prompt.
    - `update_game_config_request`: Update the configuration settings of the game.
    - `countdown_finish_response`: Return game-related information when the countdown finishes.

This setup ensures a responsive and interactive gaming experience for all connected clients by leveraging WebSocket connections for real-time communication and updates.