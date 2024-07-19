# internal/client

This directory contains client-specific logic for handling HTTP requests and WebSocket connections in the Scattergories Backend application.

## Subdirectories

- **controllers**: Contains controllers for handling HTTP requests.
- **ws**: Contains logic for handling WebSocket connections.
- **routes**: Contains route definitions.

## Controllers

The controllers in this directory are responsible for handling HTTP REST API requests. These APIs are used for user management and pre-game setup.

- **userController.go**: Contains handlers for user-related endpoints.
- **gameRoomJoinController.go**: Contains handlers for joining and leaving game rooms.
- **gameRoomController.go**: Contains handlers for game room-related endpoints.
- **helpers.go**: Contains helper functions used by the controllers.
- **gameController.go**: Contains handlers for game-related endpoints.

## WebSocket

The WebSocket implementation allows real-time communication between the server and connected clients. It includes features like submitting answers and updating game configurations.

- **hub.go**: Manages WebSocket connections and message broadcasting.
- **handlers.go**: Contains functions for handling WebSocket events.
- **client.go**: Defines the WebSocket client structure.
- **helpers.go**: Contains helper functions for WebSocket operations.

## Routes

The `routes.go` file sets up the routes for HTTP requests, linking them to the appropriate controllers.

- **userRoutes.go**: Defines routes for user-related endpoints.
- **gameRoomRoutes.go**: Defines routes for game room-related endpoints.
- **routes.go**: Registers all routes with the Gin router.
