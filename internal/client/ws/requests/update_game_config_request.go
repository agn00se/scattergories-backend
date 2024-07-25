package requests

type UpdateGameConfigRequest struct {
	Type            string `json:"type"`
	HostID          uint   `json:"host_id" validate:"required"`
	TimeLimit       int    `json:"time_limit" validate:"required,gte=10,lte=300"`
	NumberOfPrompts int    `json:"number_of_prompts" validate:"required,gte=1,lte=30"`
	Letter          string `json:"letter" validate:"required,len=1,alpha"`
}