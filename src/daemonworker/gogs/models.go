package gogs

type gogsHandler struct {
	baseURL  string
	username string
	token    string
}

type AccessToken struct {
	Name string `json:"name"`
	Sha1 string `json:"sha1"`
}

type SignIn struct {
	UserName string `binding:"Required;MaxSize(254)"`
	Password string `binding:"Required;MaxSize(255)"`
	Remember bool
}

type createAccessTokenOption struct {
	Name string `json:"name" binding:"Required"`
}

type createHookOption struct {
	Type   string            `json:"type" binding:"Required"`
	Config map[string]string `json:"config" binding:"Required"`
	Events []string          `json:"events"`
	Active bool              `json:"active"`
}
