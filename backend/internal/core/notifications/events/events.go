package events

const (
	UserSignedUp = "user.signed_up"
)

type UserSignedUpPayload struct {
	UserID         string
	Email          string
	Name           string
	OrganizationID string
}
