package response

type NewsClientResponse struct {
	Status       string
	TotalResults int
	Articles     []ArticleResponse
}

type ArticleResponse struct {
	Source      SourceResponse
	Id          string
	Name        string
	Author      string
	Title       string
	Description string
	UrlToImage  string
	PublishedAt string
	Content     string
}

type SourceResponse struct {
	Id   *string
	Name string
}
