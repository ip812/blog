package articles

import (
	"strconv"
)

type ArticleMetadata struct {
	ID          uint64
	Name        string
	URL         string
	Description string
}

var (
	PlaceholderID      uint64 = 1418336861478195200
	ZeroTrustHomelabID uint64 = 1417231583613554688
)

var Metadata = []ArticleMetadata{
	{
		ID:          ZeroTrustHomelabID,
		URL:         "/p/public/articles/" + strconv.FormatUint(ZeroTrustHomelabID, 10),
		Name:        "Zero trust homelab",
		Description: "My homelab setup using Terraform, Helm, Cloudflare, Tailscale and more...",
	},
}

func GetByID(id uint64) *ArticleMetadata {
	for _, m := range Metadata {
		if m.ID == id {
			return &m
		}
	}

	return nil
}
