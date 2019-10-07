package emailservice

// EmailSendRequest represents the data necessary for any EmailService to
// send an email
type EmailSendRequest struct {
	To         string
	ToName     string
	From       string
	FromName   string
	RawSubject string
	RawBody    string
}

// EmailService is an interface that knows how to send emails
// using a third party service
type EmailService interface {
	SendEmail(EmailSendRequest) error
}
