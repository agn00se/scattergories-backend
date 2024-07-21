# internal/repositories

This directory contains the data access logic for the Scattergories Backend application. Each repository file handles the database interactions for specific entities in the application.

- `answerRepository.go`: Contains the data access logic for managing answers.
- `gamePromptRepository.go`: Contains the data access logic for managing game prompts.
- `gameRepository.go`: Contains the data access logic for managing games.
- `gameRoomConfigRepository.go`: Contains the data access logic for managing game room configurations.
- `gameRoomRepository.go`: Contains the data access logic for managing game rooms.
- `playerRepository.go`: Contains the data access logic for managing players.
- `promptRepository.go`: Contains the data access logic for managing prompts.
- `userRepository.go`: Contains the data access logic for managing users.

Each repository interacts with the database using GORM and encapsulates the SQL queries and data manipulation logic, providing a clean and consistent API for the service layer to use.

