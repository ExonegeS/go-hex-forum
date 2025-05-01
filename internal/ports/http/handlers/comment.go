package handlers

type CommentService interface {
	// methods here ...
}

type CommentHandler struct {
	service CommentService
}

func NewCommentHandler(service CommentService) *CommentHandler {
	return &CommentHandler{service}
}
