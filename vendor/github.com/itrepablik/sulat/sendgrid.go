package sulat

import (
	"errors"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SGC collects required configurations for SendGrid SMTP server
type SGC struct {
	SendGridAPIKey, SendGridEndPoint, SendGridHost string
}

func (s *SendMail) optionsSG(sm *SendMail) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	s = sm

	// Set the required options prior to sending an email
	from := mail.NewEmail(s.From.Name, s.From.Address)
	to := mail.NewEmail(s.To.Name, s.To.Address)
	content := mail.NewContent("text/html", s.HTMLBody) // Use only the HTML format
	m := mail.NewV3MailInit(from, s.Subject, to, content)

	// Optional configs
	if len(strings.TrimSpace(s.CC.Address)) > 0 {
		cc := mail.NewEmail(s.CC.Name, s.CC.Address)
		m.Personalizations[0].AddCCs(cc)
	}
	if len(strings.TrimSpace(s.BCC.Address)) > 0 {
		bcc := mail.NewEmail(s.BCC.Name, s.BCC.Address)
		m.Personalizations[0].AddBCCs(bcc)
	}
	return mail.GetRequestBody(m)
}

func (s *SendMail) sendSG(byteMailOpt []byte, sgc *SGC) (bool, error) {
	if len(strings.TrimSpace(sgc.SendGridAPIKey)) == 0 {
		return false, errors.New("sendgrid api key is required")
	}
	if len(strings.TrimSpace(sgc.SendGridEndPoint)) == 0 {
		return false, errors.New("sendgrid endpoint is required")
	}
	if len(strings.TrimSpace(sgc.SendGridHost)) == 0 {
		return false, errors.New("sendgrid host is required")
	}

	request := sendgrid.GetRequest(sgc.SendGridAPIKey, sgc.SendGridEndPoint, sgc.SendGridHost)
	request.Method = "POST"
	var Body = byteMailOpt
	request.Body = Body

	_, err := sendgrid.API(request)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SendEmailSG dispatch an automatic email notification using SendGrid SMTP
func SendEmailSG(s *SendMail, emf *EmailHTMLFormat, sgCon *SGC) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	emailContent := ""
	if emf.IsFullHTML {
		emailContent, err = NewFullHTMLContent(emf.FullHTMLTemplate)
	} else {
		emailContent, err = NewEmailContent(emf.HTMLHeader, emf.HTMLBody, emf.HTMLFooter)
	}
	if err != nil {
		return false, err
	}

	mailOpt := SM.optionsSG(&SendMail{
		Subject:  s.Subject,
		From:     s.From,
		To:       s.To,
		CC:       s.CC,
		BCC:      s.BCC,
		HTMLBody: emailContent,
	})
	isSend, err := SM.sendSG(mailOpt, sgCon)

	if err != nil {
		return false, err
	}
	if !isSend {
		return false, errors.New("failed to send an automatic email")
	}
	return true, nil
}
