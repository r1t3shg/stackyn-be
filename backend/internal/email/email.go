// Package email provides email sending functionality using AWS SES.
package email

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Service struct {
	sesClient *ses.SES
	fromEmail string
}

func NewService() (*Service, error) {
	// Get AWS credentials from environment variables
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1" // Default region
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	fromEmail := os.Getenv("AWS_SES_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "noreply@stackyn.com" // Default from email
	}

	// Create AWS session
	config := &aws.Config{
		Region: aws.String(awsRegion),
	}

	// If credentials are provided, use them
	if awsAccessKeyID != "" && awsSecretAccessKey != "" {
		config.Credentials = credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	} else {
		// Otherwise, use default credential chain (IAM role, environment, etc.)
		log.Println("[EMAIL] Using default AWS credential chain")
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	sesClient := ses.New(sess)

	return &Service{
		sesClient: sesClient,
		fromEmail: fromEmail,
	}, nil
}

// SendOTPEmail sends an OTP verification email to the user
func (s *Service) SendOTPEmail(toEmail, otp string) error {
	subject := "Verify your Stackyn account"
	body := fmt.Sprintf(`
Hello,

Your verification code for Stackyn is: %s

This code will expire in 5 minutes.

If you didn't request this code, please ignore this email.

Best regards,
The Stackyn Team
`, otp)

	// HTML version
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center; border-radius: 8px 8px 0 0;">
		<h1 style="color: white; margin: 0; font-size: 28px;">Stackyn</h1>
	</div>
	<div style="background: #ffffff; padding: 40px; border: 1px solid #e0e0e0; border-top: none; border-radius: 0 0 8px 8px;">
		<h2 style="color: #333; margin-top: 0;">Verify your email address</h2>
		<p style="color: #666; font-size: 16px;">Your verification code is:</p>
		<div style="background: #f5f5f5; border: 2px dashed #667eea; border-radius: 8px; padding: 20px; text-align: center; margin: 30px 0;">
			<code style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 4px;">%s</code>
		</div>
		<p style="color: #666; font-size: 14px;">This code will expire in 5 minutes.</p>
		<p style="color: #999; font-size: 12px; margin-top: 30px; border-top: 1px solid #e0e0e0; padding-top: 20px;">If you didn't request this code, please ignore this email.</p>
	</div>
</body>
</html>
`, otp)

	// Create email input
	input := &ses.SendEmailInput{
		Source: aws.String(s.fromEmail),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toEmail)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
				Html: &ses.Content{
					Data:    aws.String(htmlBody),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	// Send email
	result, err := s.sesClient.SendEmail(input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	log.Printf("[EMAIL] OTP email sent successfully to %s (MessageId: %s)", toEmail, *result.MessageId)
	return nil
}

