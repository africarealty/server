package domain

import (
	"context"
	"github.com/africarealty/server/src/kit/queue/listener"
	"github.com/africarealty/server/src/service"
)

const (
	EmailRqStatusNotSend   = "not-send"
	EmailRqStatusValErr    = "validation-error"
	EmailRqStatusSent      = "sent"
	EmailRqStatusSmtpError = "smtp-error"

	EmailRequestTopic = "email.request"

	EmailTemplateUserActivation = "auth.registration-activation"
)

type Template struct {
	Id    string // Id template ID
	Title string // Title template title
	Body  string // Body template body
}

// TemplateService is responsible for manage Template entities
type TemplateService interface {
	// CreateTemplate creates template
	CreateTemplate(ctx context.Context, rq *Template) (*Template, error)
	// UpdateTemplate updates template
	UpdateTemplate(ctx context.Context, rq *Template) (*Template, error)
	// DeleteTemplate delete template
	DeleteTemplate(ctx context.Context, id string) error
	// GetTemplate get template by id
	GetTemplate(ctx context.Context, id string) (*Template, error)
	// SearchTemplates searches templates by query
	SearchTemplates(ctx context.Context, query string) ([]*Template, error)
}

type TemplateRequest struct {
	Id   string                 // Id template
	Data map[string]interface{} // Data - object to be used for text generation based on a template
}

type TemplateResponse struct {
	Title string // Title template title
	Body  string // Body template body
}

// TemplateGenerator used for generating body based on templates
type TemplateGenerator interface {
	// Generate - generates a text message based on template ID and data object
	Generate(ctx context.Context, rq *TemplateRequest) (*TemplateResponse, error)
}

// TemplateStorage provides storage for templates
type TemplateStorage interface {
	// Get retrieves a template by Id
	Get(ctx context.Context, id string) (*Template, error)
	// Create creates template
	Create(ctx context.Context, template *Template) error
	// Update updates
	Update(ctx context.Context, template *Template) error
	// Delete delete a template by Id
	Delete(ctx context.Context, id string) error
	// Search searches templates by query
	Search(ctx context.Context, query string) ([]*Template, error)
}

type EmailRequest struct {
	// Id should be empty for a new request
	Id string
	// UserId(Optional) - if sending to specific user
	UserId string
	// Email of user
	Email string
	// Template (Optional) - empty if sending confirmation code
	Template *TemplateRequest
	// Subject - email subject
	Subject string
	//Text - final text of email message
	Text string
	// LinkFileIds (Optional) - email attachments fileID list
	LinkFileIds []string
	// From (Optional) - email From field ('Display From' as well as 'Mail From')
	From string
}

// Email object
type Email struct {
	// Id should be empty for a new request
	Id string
	// UserId(Optional) - if sending to specific user
	UserId string
	// Subject - email subject
	Subject string
	//Text - final text of email message
	Text string
	// Email of user
	Email string
	// Template (Required) - empty if sending text
	Template *TemplateRequest
	// SendStatus - current status of sending
	SendStatus string
	// ErrorDescription gives some more details about error
	ErrorDescription string
	// LinkFileIds - email attachments fileID list
	LinkFileIds []string
	// From - email From field ('Display From' as well as 'Mail From')
	From string
	// Attachments - attached files
	Attachments []*FileContent
}

// EmailService is responsible for sending email messages
type EmailService interface {
	// Init inits with a config
	Init(cfg *service.Config) error
	// RequestHandler consumes all email requests from the queue
	RequestHandler() listener.QueueMessageHandler
	// Send puts a request for sending
	// real sending is going asynchronously
	Send(ctx context.Context, rq *EmailRequest) (*Email, error)
}

// EmailRepository is responsible for real sending email messages
type EmailRepository interface {
	// Send real sends an email
	Send(ctx context.Context, email *Email) error
}

// EmailStorage persists data to storages
type EmailStorage interface {
	// CreateEmail creates SMS request to db
	CreateEmail(ctx context.Context, requests *Email) error
	// UpdateEmail updates SMS request to db
	UpdateEmail(ctx context.Context, requests *Email) error
}
