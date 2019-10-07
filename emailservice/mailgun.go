package emailservice

type mailgunEmailService struct {
}

func NewMailgunEmailService() EmailService {
	return mailgunEmailService{}
}

func (meg *MailgunEmailService) SendEmail(EmailSendRequest) error {

}
