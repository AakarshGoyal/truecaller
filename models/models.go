package models

type Response struct {
	RequestID   string `json:"requestID"`
	AccessToken string `json:"accessToken"`
	Endpoint    string `json:"endpoint"`
}

type Rupifi_FE struct {
	AccessToken string `json:"accessToken"`
}
