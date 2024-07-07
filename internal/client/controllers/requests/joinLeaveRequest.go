package requests

type JoinLeaveRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
