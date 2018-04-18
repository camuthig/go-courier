package courier

// Address is the combination of a email and name to deliver emails to.
type Address struct {
	Email, Name string
}

// ToRfc2822 converts the Address object into a RFC 2822 compatible address string.
func (a Address) ToRfc2822() string {
	if a.Name != "" {
		return "\"" + a.Name + "\" <" + a.Email + ">"
	}

	return a.Email
}

// Email is just what it sounds like.
type Email struct {
	From        Address
	ReplyTo     *Address
	Subject     string
	To, Cc, Bcc []Address
	Headers     map[string]string
	Attachments []Attachment
	Content     interface{}
}

// SimpleContent is a normal email with a text and a HTML body.
type SimpleContent struct {
	Text, HTML string
}

// TemplatedContent is meant to work with SaaS providers and requies an ID and a map of substitution data.
type TemplatedContent struct {
	TemplateID       string
	SubstitutionData map[string]string
}

// Courier is a class that is capable of sending emails using third-party providers
type Courier interface {
	// Send an email. The return will either be a receipt ID as defined by the courier or an error.
	Send(e Email) (string, error)
}
