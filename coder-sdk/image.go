package coder

import (
	"context"
	"net/http"
)

// Image describes a Coder Image
type Image struct {
	ID              string  `json:"id"`
	OrganizationID  string  `json:"organization_id"`
	Repository      string  `json:"repository"`
	Description     string  `json:"description"`
	URL             string  `json:"url"` // user-supplied URL for image
	DefaultCPUCores float32 `json:"default_cpu_cores"`
	DefaultMemoryGB int     `json:"default_memory_gb"`
	DefaultDiskGB   int     `json:"default_disk_gb"`
	Deprecated      bool    `json:"deprecated"`
}

// NewRegistryRequest describes a docker registry used in importing an image
type NewRegistryRequest struct {
	FriendlyName string `json:"friendly_name"`
	Registry     string `json:"registry"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// ImportImageRequest is used to import new images and registries into Coder
type ImportImageRequest struct {
	// RegistryID is used to import images to existing registries.
	RegistryID *string `json:"registry_id"`
	// NewRegistry is used when adding a new registry.
	NewRegistry *NewRegistryRequest `json:"new_registry"`
	// Repository refers to the image. For example: "codercom/ubuntu".
	Repository      string  `json:"repository"`
	Tag             string  `json:"tag"`
	DefaultCPUCores float32 `json:"default_cpu_cores"`
	DefaultMemoryGB int     `json:"default_memory_gb"`
	DefaultDiskGB   int     `json:"default_disk_gb"`
	Description     string  `json:"description"`
	URL             string  `json:"url"`
}

// ImportImage creates a new image and optionally a new registry
func (c Client) ImportImage(ctx context.Context, orgID string, req ImportImageRequest) (*Image, error) {
	var img *Image
	err := c.requestBody(
		ctx,
		http.MethodPost, "/api/orgs/"+orgID+"/images",
		req,
		img,
	)
	return img, err
}
