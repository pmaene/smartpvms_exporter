package smartpvms

type Result struct {
	Success  bool
	FailCode int `json:"failCode"`
	Params   *struct {
		CurrentTime int `json:"currentTime"`
	}
	Message *string
	Data    interface{}
}

type LoginBody struct {
	Username   string `json:"userName"`
	SystemCode string `json:"systemCode"`
}

type LoginResult struct {
	Result
}

type LogoutBody struct {
	XSRFToken string `json:"xsrfToken"`
}

type LogoutResult struct {
	Result
}
