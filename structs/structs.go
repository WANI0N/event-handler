package structs

type EventData struct {
	Id           string   `json:"id"`
	Name         string   `json:"name" binding:"required,min=1,max=255,checkEventName"`
	Timestamp    string   `json:"date" binding:"required,checkTimeFieldFormat"`
	Languages    []string `json:"languages" binding:"required,min=1,unique"`
	VideoQuality []string `json:"videoQuality" default:"720p" binding:"checkVideoQuality"`
	AudioQuality []string `json:"audioQuality" default:"Low" binding:"checkAudioQuality"`
	Invitees     []string `json:"invitees" binding:"required,min=1,max=100,unique,checkEmail"`
	Description  string   `json:"description"  binding:"max=512"`
}

type JsonHealthCheckStatus struct {
	Result     string `json:"result"`
	DeployDate string `json:"deployDate"`
	Version    string `json:"version"`
}
