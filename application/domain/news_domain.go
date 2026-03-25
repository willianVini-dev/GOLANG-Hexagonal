package domain

type NewsDomain struct {
	Status       string
	TotalResults int
	Articles     []Article
}

type Article struct {
	Source      SourceResponse
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

type NewsRequestDomain struct {
	Subject string
	From    string
}
