# internal/services

This directory contains the business logic and services used by the controllers in the Scattergories Backend application.

- `answer_service.go`: Contains the business logic for managing answers.
- `auth_service.go`: Contains the business logic for authenticating users.
- `game_prompt_service.go`: Contains the business logic for managing game prompts.
- `game_service.go`: Contains the business logic for managing games.
- `game_room_config_service.go`: Contains the business logic for managing game room configurations.
- `game_room_data_service.go`: Manages the retrieval of data for game rooms.
- `game_room_join_service.go`: Contains the business logic for joining and leaving game rooms.
- `game_room_service.go`: Contains the business logic for managing game rooms.
- `permission_service.go`: Manages permissions for various actions within the application.
- `player_service.go`: Contains the business logic for managing players.
- `prompt_service.go`: Contains the business logic for managing prompts.
- `token_service.go`: Manages token generation, refresh, validation, and invalidation.
- `user_service.go`: Contains the business logic for managing users.
- `user_registration_service.go`: Handles user registration processes.

Each service interacts with the data models and encapsulates the business rules and logic of the application.
