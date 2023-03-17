package interfaces

import (
	"context"
	"encoding/json"
)

// type constrain the response in these functions a little more
// for more type safety.
type PayloadConstructor func(string, int, string) (Messagable, string)

type DeliveryBasePayload struct {
	EmailTo           string
	EmailFrom         string
	MailgunTemplateID string
	Subject           string
	SchoolID          int
}

// BasePayload implements the Payloadable interface.
func (dbp *DeliveryBasePayload) BasePayload() DeliveryBasePayload {
	// NOTE: potential panic here with nil deref, maybe be more careful
	return *dbp
}

type AbandonedCartMetadata struct {
	NewPlan         string `json:"new_plan"`
	UserName        string `json:"user_name"`
	Email           string `json:"email"`
	SchoolDomainUrl string `json:"school_primary_domain_url"`
}

type AbandonedCartPayload struct {
	DeliveryBasePayload
	MailgunTemplateVars AbandonedCartMetadata
}

// TemplateVars implements the Templatable interface.
func (acp *AbandonedCartPayload) TemplateVars() interface{} {
	return &acp.MailgunTemplateVars
}

// fake requests.Message
type Message struct {
	To                string
	From              string
	Subject           string
	Template          string
	TemplateVariables string
}

type Templatable interface {
	TemplateVars() interface{}
}

type Payloadable interface {
	BasePayload() DeliveryBasePayload
}

type Messagable interface {
	Templatable
	Payloadable
}

// originally this:
// func (acp *AbandonedCartPayload) CreateMessageFromPayload(logger *logging.AppLogger) requests.Message {
// removed logger for simplicity in the example
// any payload type that embeds `DeliveryBasePayload` needs a stub implementation of  `TemplateVars`
// to be usable in this function.
func CreateMessageFromPayload(payload Messagable) Message {
	basePayload := payload.BasePayload()
	msg := Message{
		To:       basePayload.EmailTo,
		From:     basePayload.EmailFrom,
		Subject:  basePayload.Subject,
		Template: basePayload.MailgunTemplateID,
	}

	vars := payload.TemplateVars()
	templateVars, err := json.Marshal(vars)
	if vars == nil || err != nil {
		// do your logging here
		return msg
	}

	msg.TemplateVariables = string(templateVars)
	return msg
}

// originally this:
// func (acp *AbandonedCartPayload) HandleEmailSend(ctx context.Context, logger *logging.AppLogger, gdb *gorm.DB, emSvc *requests.EmailService) error {
// simplifying the signature for the example
// any payload type that embeds `DeliveryBasePayload` needs a stub implementation of  `TemplateVars`
// to be usable in this function.
func HandleEmailSend(ctx context.Context, payload Messagable) error {
	msg := CreateMessageFromPayload(payload)
	basePayload := payload.BasePayload()
	// I commented out usages, this is here to make the example compile
	_ = msg
	_ = basePayload

	// send your email and get your response
	// resp, id, err := emSvc.Send(ctx, &msg)

	// I commented out the declaration, this is here to make the example compile
	var err error
	if err != nil {
		// send your response and handle errors
		// err := InsertDBRecordHelper(
		// 	logger,
		// 	gdb,
		// 	basePayload.EmailTo,
		// 	basePayload.SchoolID,
		// 	basePayload.MailgunTemplateID,
		// 	models.FAILED)
		// if err != nil {
		// 	return err
		// }
		// logger.Info("EmailDelivery Failed")
		// return err
	}

	// do other handling here
	// logger.Info(
	// 	"Email sent",
	// 	logging.StringMetadata("email", basePayload.EmailTo),
	// 	logging.StringMetadata("templateID", basePayload.MailgunTemplateID),
	// 	logging.StringMetadata("messageID", id),
	// 	logging.StringMetadata("response", resp),
	// )
	// err = InsertDBRecordHelper(
	// 	logger,
	// 	gdb,
	// 	basePayload.EmailTo,
	// 	basePayload.SchoolID,
	// 	basePayload.MailgunTemplateID,
	// 	models.DELIVERED,
	// )
	// if err != nil {
	// 	return err
	// }

	return nil
}
