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
	PlaceholderID        uint64 = 1418336861478195200
	ZeroTrustHomelabID   uint64 = 1417231583613554688
	ZeroTrustHomelabV2ID uint64 = 1428029051347406848
)

var Metadata = []ArticleMetadata{
	{
		ID:          ZeroTrustHomelabV2ID,
		URL:         "/p/public/articles/" + strconv.FormatUint(ZeroTrustHomelabV2ID, 10),
		Name:        "Zero trust homelab V2",
		Description: "An updated version of my homelab setup using FluxCD, Doppler and my own Terraform provider.",
	},
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
