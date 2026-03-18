package port

import "context"

type EmailMessage struct {
	To      string
	Subject string
	Content string
}

type EmailSender interface {
	Send(ctx context.Context, msg EmailMessage) error
}
