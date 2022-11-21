package resp

type User struct {
	Username string   `json:"username"`
	Projects []string `json:"projects"`
}
