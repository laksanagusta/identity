package dtos

import (
	"strings"
	"unicode"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type CreateNewOrganizationReq struct {
	Name      nullable.NullString `json:"name"`
	Address   nullable.NullString `json:"address"`
	Latitude  nullable.NullString `json:"latitude"`
	Longitude nullable.NullString `json:"longitude"`
	Type      nullable.NullString `json:"type"`
	ParentId  nullable.NullString `json:"parent_id"`
}

func (r CreateNewOrganizationReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Address, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Latitude, validation.Required, is.Latitude),
		validation.Field(&r.Longitude, validation.Required, is.Longitude),
		validation.Field(&r.Latitude, validation.Required),
		validation.Field(&r.Longitude, validation.Required),
		validation.Field(&r.Type, validation.Required, validation.Length(1, 255)),
	)
}

// generateOrganizationCode generates a code from organization name by converting to lowercase and replacing spaces with hyphens
func generateOrganizationCode(name string) string {
	if name == "" {
		return ""
	}

	// Convert to lowercase and replace spaces/multiple spaces with single hyphen
	words := strings.Fields(strings.ToLower(name))
	var codeWords []string

	for _, word := range words {
		// Remove special characters, keep only alphanumeric and spaces
		var cleanWord strings.Builder
		for _, r := range word {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				cleanWord.WriteRune(r)
			}
		}
		if cleanWord.Len() > 0 {
			codeWords = append(codeWords, cleanWord.String())
		}
	}

	return strings.Join(codeWords, "-")
}

func (r CreateNewOrganizationReq) NewOrganization(cred entities.AuthenticatedUser) entities.Organization {
	organization := entities.Organization{
		Name:       r.Name,
		Address:    r.Address,
		Type:       r.Type,
		ParentUUID: r.ParentId,
	}

	baseModel := entities.NewBaseModel(cred.Username)
	organization.BaseModel = baseModel

	// Generate organization code from name or use fallback
	var orgCode string
	if r.Name.IsExists && r.Name.Val != nil {
		orgCode = generateOrganizationCode(*r.Name.Val)
	} else {
		// Fallback: generate code from UUID if name is not available
		orgCode = "org-" + organization.UUID[:8] // Use first 8 characters of UUID
	}
	organization.Code = nullable.NewString(orgCode)

	return organization
}
