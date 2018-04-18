package courier

import (
	"fmt"
	"reflect"
	"strings"

	sp "github.com/SparkPost/gosparkpost"
)

// SparkPostCourier supports sending emails through SparkPost
type SparkPostCourier struct {
	client sp.Client
}

// Send an email
func (spc SparkPostCourier) Send(e Email) (string, error) {
	switch cType := e.Content.(type) {
	case SimpleContent:
		return spc.sendSimpleEmail(e)
	case TemplatedContent:
		return spc.sendTemplatedEmail(e)
	default:
		return "", fmt.Errorf("Unsupport type %s", reflect.TypeOf(cType))
	}
}

func (spc SparkPostCourier) sendSimpleEmail(e Email) (string, error) {
	content := e.Content.(SimpleContent)
	var contentHeaders = e.Headers

	// @TODO Make sure this is the correct check
	if e.Cc != nil {
		contentHeaders["CC"] = createRecipientsHeader(e.Cc)
	}

	tx := &sp.Transmission{
		Recipients: createRecipients(e),
		Content: sp.Content{
			HTML:        content.HTML,
			Text:        content.Text,
			From:        e.From.ToRfc2822(),
			Subject:     e.Subject,
			Headers:     contentHeaders,
			Attachments: createAttachments(e),
		},
	}

	id, _, err := spc.client.Send(tx)

	return id, err
}

func (spc SparkPostCourier) sendTemplatedEmail(e Email) (string, error) {
	content := e.Content.(TemplatedContent)
	substitutionData := content.SubstitutionData

	if e.Cc != nil {
		substitutionData["CC"] = createRecipientsHeader(e.Cc)
	}

	if e.ReplyTo != nil {
		substitutionData["replyTo"] = e.ReplyTo.ToRfc2822()
	}

	if e.From.Name != "" {
		substitutionData["fromName"] = e.From.Name
	}

	substitutionData["subject"] = e.Subject
	substitutionData["fromEmail"] = strings.Split(e.From.Email, "@")[0]
	substitutionData["fromDomain"] = strings.SplitN(e.From.Email, "@", 2)[1]

	tx := &sp.Transmission{
		Recipients: createRecipients(e),
		Content: map[string]string{
			"template_id": content.TemplateID,
		},
		SubstitutionData: substitutionData,
	}

	id, _, err := spc.client.Send(tx)

	return id, err
}

func createRecipients(e Email) []sp.Recipient {
	var as []sp.Recipient

	for _, a := range e.To {
		as = append(as, sp.Recipient{Address: sp.Address{
			Email:    a.Email,
			HeaderTo: createRecipientsHeader(e.To),
		}})
	}

	for _, a := range e.Cc {
		as = append(as, sp.Recipient{Address: sp.Address{
			Email:    a.Email,
			HeaderTo: createRecipientsHeader(e.To),
		}})
	}

	for _, a := range e.Bcc {
		as = append(as, sp.Recipient{Address: sp.Address{
			Email:    a.Email,
			HeaderTo: createRecipientsHeader(e.To),
		}})
	}

	return as
}

func createRecipientsHeader(to []Address) string {
	var as []string

	for _, a := range to {
		as = append(as, a.ToRfc2822())
	}

	return strings.Join(as, ",")
}

func createAttachments(e Email) []sp.Attachment {
	r := []sp.Attachment{}

	for _, a := range e.Attachments {
		// @TODO Handle the error on this
		c, _ := a.Base64Content()

		r = append(r, sp.Attachment{
			MIMEType: a.ContentType(),
			Filename: a.Name(),
			B64Data:  c,
		})
	}

	return r
}
