package articles

import "strconv"

type ArticleMetadata struct {
	ID          uint64
	Name        string
	URL         string
	Description string
}

func newArticle(id uint64, name, description string) ArticleMetadata {
	return ArticleMetadata{
		ID:          id,
		URL:         "/p/public/articles/" + strconv.FormatUint(id, 10),
		Name:        name,
		Description: description,
	}
}

var Metadata = []ArticleMetadata{
	newArticle(
		1417231583613554688,
		"Zero trust homelab",
		"My homelad setup using Terraform, Helm, CLoudflare, Tailscale and more...",
	),
}
