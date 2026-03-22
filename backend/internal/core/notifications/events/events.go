package events

const (
	UserSignedUp         = "user.signed_up"
	UserVerificationEmail = "user.verification_email"
	UserPasswordReset    = "user.password_reset"
)

type UserSignedUpPayload struct {
	UserID         string
	Email          string
	Name           string
	OrganizationID string
}

type UserVerificationEmailPayload struct {
	UserID          string
	Email           string
	Name            string
	VerificationURL string
}

type UserPasswordResetPayload struct {
	UserID   string
	Email    string
	Name     string
	ResetURL string
}
