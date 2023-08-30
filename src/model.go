package src

type (
	httpResponseGeneral struct {
		content []byte
		code    int
	}

	UserData struct {
		ID           string   `json:"_id"`
		ActiveStatus bool     `json:"isActive"`
		Balance      string   `json:"balance"`
		Tags         []string `json:"tags"`
	}
)
