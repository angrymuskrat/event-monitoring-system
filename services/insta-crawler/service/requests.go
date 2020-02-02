package service

type idEpRequest struct {
	ID string `json:"id"`
}

type postsEpRequest struct {
	ID     string `json:"id"`
	Offset string `json:"offset"`
	Num    int    `json:"num"`
}
