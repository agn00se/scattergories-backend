# internal/api

This directory contains client-specific logic for handling HTTP requests and WebSocket connections in the Scattergories Backend application.

## Subdirectories

- **handlers**: Contains handlers for HTTP requests.
- **websocket**: Contains logic for handling WebSocket connections.
- **routes**: Contains route definitions.
- **middleware**: Contains middleware functions, such as JWT authentication.

## Handlers

The handlers in this directory are responsible for handling HTTP REST API requests. These APIs are used for user management and pre-game setup.

- **user_handler.go**: Contains handlers for user-related endpoints.
- **auth_handler.go**: Contains handlers for auth-related endpoints.
- **gameroom_join_handler.go**: Contains handlers for joining and leaving game rooms.
- **gameroom__handler.go**: Contains handlers for game room-related endpoints.

## WebSocket

The WebSocket implementation allows real-time communication between the server and connected clients. It includes features like submitting answers and updating game configurations.

- **hub.go**: Manages WebSocket connections and message broadcasting.
- **handlers.go**: Contains functions for handling WebSocket events.
- **client.go**: Defines the WebSocket client structure.

## Routes

The `routes.go` file sets up the routes for HTTP and Websocket requests, linking them to the appropriate handlers.

- **user_routes.go**: Defines routes for user-related endpoints.
- **gameroom_routes.go**: Defines routes for game room-related endpoints.
- **routes.go**: Registers all routes with the Gin router.

## Middlewares

- **jwt_middleware.go**: Contains JWT authentication middleware to secure routes.