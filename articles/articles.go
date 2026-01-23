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
	DeferDeepDiveID                        uint64 = 1458103253970456576
	SelfManagedObservabilityStackID        uint64 = 1463957572842164224
)

var Metadata = []ArticleMetadata{
	{
		ID:              SelfManagedObservabilityStackID,
		URL:             "/p/public/articles/" + strconv.FormatUint(SelfManagedObservabilityStackID, 10),
		Name:            "From Grafana Cloud to a self-managed observability stack",
		Description:     "Why I decided to manage my own observability stack for my homelab and how I did it.",
		ReadTimeMinutes: 10,
	},
	{
		ID:              DeferDeepDiveID,
		URL:             "/p/public/articles/" + strconv.FormatUint(DeferDeepDiveID, 10),
		Name:            "Defer in Go: Deep Dive",
		Description:     "How defer works in Go, common pitfalls and best practices.",
		ReadTimeMinutes: 8,
	},
	{
		ID:              AnsiblePlusTailsclaleEqualGreatComboID,
		URL:             "/p/public/articles/" + strconv.FormatUint(AnsiblePlusTailsclaleEqualGreatComboID, 10),
		Name:            "Ansible + Tailscale = 🎉 ",
		Description:     "Manage VMs in a private network with Ansible and Tailscale.",
		ReadTimeMinutes: 4,
	},
	{
		ID:              ZeroTrustHomelabV2ID,
		URL:             "/p/public/articles/" + strconv.FormatUint(ZeroTrustHomelabV2ID, 10),
		Name:            "Zero trust homelab V2",
		Description:     "An updated version of my homelab setup using FluxCD, Doppler and my own Terraform provider.",
		ReadTimeMinutes: 6,
	},
	{
		ID:              ZeroTrustHomelabID,
		URL:             "/p/public/articles/" + strconv.FormatUint(ZeroTrustHomelabID, 10),
		Name:            "Zero trust homelab",
		Description:     "My homelab setup using Terraform, Helm, Cloudflare, Tailscale and more...",
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
