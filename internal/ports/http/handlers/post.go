package handlers

type PostService interface {
	// methods here ...
}

type PostHandler struct {
	service  PostService
	sessions SessionService
	// templates *template.Template
}

func NewPostHandler(service PostService, sessions SessionService) *PostHandler {
	return &PostHandler{service, sessions}
}
