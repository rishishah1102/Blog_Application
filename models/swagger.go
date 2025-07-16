package models

type BlogPostRequest struct {
	Title       string `json:"title" example:"Blog-1"`
	Description string `json:"description" example:"This is the sample description of blog"`
	Image       string `json:"image" example:"https://sample-image.png"`
}

type CreateBlogPostSuccessResponse struct {
	Message string   `json:"message" example:"successfully created the blog"`
	Blog    BlogPost `json:"blog"`
}

type GetAllBlogPostsSuccessResponse struct {
	Message string     `json:"message" example:"successfully fetched all the blogs"`
	Blog    []BlogPost `json:"blog"`
}

type GetBlogPostByIDSuccessResponse struct {
	Message string   `json:"message" example:"successfully fetched a blog"`
	Blog    BlogPost `json:"blog"`
	BlogID  string   `json:"blogID" example:"6876a210c3abc4d716fb11d8"`
}

type UpdateBlogPostSuccessResponse struct {
	Message string   `json:"message" example:"successfully updated a blog"`
	Blog    BlogPost `json:"blog"`
	BlogID  string   `json:"blogID" example:"6876a210c3abc4d716fb11d8"`
}

type DeleteBlogPostSuccessResponse struct {
	Message string `json:"message" example:"blog post deleted successfully"`
	BlogID  string `json:"blogID" example:"6876a210c3abc4d716fb11d8"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
