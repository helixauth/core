package email

import (
	"context"
	"fmt"
	"log"

	"github.com/helixauth/core/src/entity"
	"github.com/helixauth/core/src/lib/mapper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type sesGateway struct {
	awsSession *session.Session
	sesSession *ses.SES
}

func NewSESGateway(ctx context.Context, tenant *entity.Tenant) (Gateway, error) {
	awsSess, err := session.NewSession(&aws.Config{
		Region:      tenant.AWSRegionID,
		Credentials: credentials.NewStaticCredentials(mapper.String(tenant.AWSAccessKeyID), mapper.String(tenant.AWSSecretAccessKey), ""),
	})
	sesSess := ses.New(awsSess)
	return &sesGateway{
		awsSession: awsSess,
		sesSession: sesSess,
	}, err
}

func (g *sesGateway) SendEmail(ctx context.Context, sender string, recipient string, subject string, htmlBody string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	result, err := g.sesSession.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	if result.MessageId != nil {
		log.Printf("Sent message %s", *result.MessageId)
	} else {
		log.Printf("Missing message ID")
	}

	return err
}
