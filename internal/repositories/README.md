# internal/repositories

This directory contains the data access logic for the Scattergories Backend application. Each repository file handles the database interactions for specific entities in the application.

- `answer_repository.go`: Contains the data access logic for managing answers.
- `game_prompt_repository.go`: Contains the data access logic for managing game prompts.
- `game_repository.go`: Contains the data access logic for managing games.
- `gameroom_config_repository.go`: Contains the data access logic for managing game room configurations.
- `gameroom_repository.go`: Contains the data access logic for managing game rooms.
- `player_repository.go`: Contains the data access logic for managing players.
- `prompt_repository.go`: Contains the data access logic for managing prompts.
- `user_repository.go`: Contains the data access logic for managing users.

Each repository interacts with the database using GORM and encapsulates the SQL queries and data manipulation logic, providing a clean and consistent API for the service layer to use.

