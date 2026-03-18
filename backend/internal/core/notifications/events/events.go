package events

const (
	UserSignedUp = "user.signed_up"
)

type UserSignedUpPayload struct {
	UserID      string
	Email       string
	FirstName   string
	WorkspaceID string
}
