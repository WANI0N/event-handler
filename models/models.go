package models

// @Description If not provided, `videoQuality` & `audioQuality` default to `["720p"]` & `["Low"]`, respectively.
// @Description If provided, first item in the list is event's default quality.
type EventData struct {
	Id string `json:"-"`
	//allowed chars: A-Za-z0-9 _-
	Name string `json:"name" example:"A event-Name3_x" binding:"required,min=1,max=255,checkEventName"`
	//YYYY-MM-DDTHH:MM:SSZ
	Timestamp    string   `json:"date" example:"2006-01-02T15:04:05Z" binding:"required,checkTimeFieldFormat"`
	Languages    []string `json:"languages" example:"English,French" binding:"required,min=1,unique"`
	VideoQuality []string `json:"videoQuality" example:"720p,1080p,1440p,2160p" binding:"checkVideoQuality,unique"`
	AudioQuality []string `json:"audioQuality" example:"Low,Mid,High" binding:"checkAudioQuality,unique"`
	Invitees     []string `json:"invitees" example:"example@mail.com" binding:"required,min=1,max=100,unique,checkEmail"`
	Description  string   `json:"description"  binding:"max=512"`
}

type EventResponseData struct {
	Id string `json:"id" example:"db6bed50-7172-4051-86ab-d1e90705c692"`
	EventData
}

type JsonHealthCheckStatus struct {
	Result     string `json:"result"`
	DeployDate string `json:"deployDate"`
	Version    string `json:"version"`
}
