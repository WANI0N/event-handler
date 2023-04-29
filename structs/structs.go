package structs

// @Description If not provided, `videoQuality` & `audioQuality` default to `["720p"]` & `["Low"]`, respectively.
// @Description If provided, first item in the list is default.
type EventData struct {
	Id string `json:"id"`
	//allowed chars: A-Za-z0-9 _-
	Name string `json:"name" binding:"required,min=1,max=255,checkEventName"`
	//YYYY-MM-DDTHH:MM:SSZ
	Timestamp    string   `json:"date" binding:"required,checkTimeFieldFormat"`
	Languages    []string `json:"languages" binding:"required,min=1,unique"`
	VideoQuality []string `json:"videoQuality" binding:"checkVideoQuality"`
	AudioQuality []string `json:"audioQuality" binding:"checkAudioQuality"`
	//["example@mail.com"]
	Invitees    []string `json:"invitees" binding:"required,min=1,max=100,unique,checkEmail"`
	Description string   `json:"description"  binding:"max=512"`
}

type JsonHealthCheckStatus struct {
	Result     string `json:"result"`
	DeployDate string `json:"deployDate"`
	Version    string `json:"version"`
}
