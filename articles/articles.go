package articles

import (
	"strconv"
)

type ArticleMetadata struct {
	ID              uint64
	Name            string
	URL             string
	Description     string
	ReadTimeMinutes int
}

var (
	PlaceholderID                          uint64 = 1418336861478195200
	ZeroTrustHomelabID                     uint64 = 1417231583613554688
	ZeroTrustHomelabV2ID                   uint64 = 1428029051347406848
	AnsiblePlusTailsclaleEqualGreatComboID uint64 = 1428744843063988224
)

var Metadata = []ArticleMetadata{
	{
		ID:          AnsiblePlusTailsclaleEqualGreatComboID,
		URL:         "/p/public/articles/" + strconv.FormatUint(AnsiblePlusTailsclaleEqualGreatComboID, 10),
		Name:        "Ansible + Tailscale = 🎉 ",
		Description: "Manage VMs in a private network with Ansible and Tailscale.",
		ReadTimeMinutes: 4,
	},
	{
		ID:          ZeroTrustHomelabV2ID,
		URL:         "/p/public/articles/" + strconv.FormatUint(ZeroTrustHomelabV2ID, 10),
		Name:        "Zero trust homelab V2",
		Description: "An updated version of my homelab setup using FluxCD, Doppler and my own Terraform provider.",
		ReadTimeMinutes: 6,
	},
	{
		ID:          ZeroTrustHomelabID,
		URL:         "/p/public/articles/" + strconv.FormatUint(ZeroTrustHomelabID, 10),
		Name:        "Zero trust homelab",
		Description: "My homelab setup using Terraform, Helm, Cloudflare, Tailscale and more...",
		ReadTimeMinutes: 9,
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
