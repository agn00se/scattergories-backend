package controllers

// func GetGamesByRoomID(c *gin.Context) {
// 	roomID, err := GetIDParam(c, "room_id")
// 	if err != nil {
// 		HandleError(c, http.StatusBadRequest, "Invalid room ID")
// 		return
// 	}

// 	games, err := services.GetGamesByRoomID(roomID)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			HandleError(c, http.StatusNotFound, "Room not found")
// 		} else {
// 			HandleError(c, http.StatusInternalServerError, "Failed to retrieve games")
// 		}
// 		return
// 	}

// 	var response []responses.GameResponse
// 	for _, game := range games {
// 		response = append(response, responses.ToGameResponse(game))
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// func GetGame(c *gin.Context) {
// 	roomID, err := GetIDParam(c, "room_id")
// 	if err != nil {
// 		HandleError(c, http.StatusBadRequest, "Invalid room ID")
// 		return
// 	}

// 	gameID, err := GetIDParam(c, "game_id")
// 	if err != nil {
// 		HandleError(c, http.StatusBadRequest, "Invalid game ID")
// 		return
// 	}

// 	game, err := services.GetGameByID(roomID, gameID)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			HandleError(c, http.StatusNotFound, "Game not found")
// 		} else {
// 			HandleError(c, http.StatusInternalServerError, "Failed to get game")
// 		}
// 		return
// 	}

// 	response := responses.ToGameResponse(game)
// 	c.JSON(http.StatusOK, response)
// }

// func CreateGame(c *gin.Context) {
// 	roomID, err := GetIDParam(c, "room_id")
// 	if err != nil {
// 		HandleError(c, http.StatusBadRequest, "Invalid room ID")
// 		return
// 	}

// 	game, err := services.CreateGame(roomID)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			HandleError(c, http.StatusNotFound, "Room not found")
// 		} else if err == services.ErrActiveGameExists {
// 			HandleError(c, http.StatusConflict, err.Error())
// 		} else {
// 			HandleError(c, http.StatusInternalServerError, "Failed to create game")
// 		}
// 		return
// 	}

// 	response := responses.ToGameResponse(game)
// 	c.JSON(http.StatusOK, response)
// }
