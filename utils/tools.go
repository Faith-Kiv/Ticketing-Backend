package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Faith-Kiv/Ticketing-Backend/models"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

const (
	GORM_SEARCH_INPUT_COUNT  = 2
	RANGE_SEARCH_PARAM_COUNT = 2
)

func ReadFile(file string) (string, error) {
	openedFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	byteValue, _ := io.ReadAll(openedFile)

	return string(byteValue[:]), nil

}

func ISODateConversion(timestamp string) (date time.Time, err error) {
	date, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(timestamp))
	if err != nil {
		logrus.Error(err)
		err = errors.New("invalid date")
		return
	}
	return
}

func ModelValidationResponse(err error) ([]map[string]string, error) {

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		invalid := make([]map[string]string, len(ve))
		for i, fe := range ve {
			invalid[i] = map[string]string{"field": fe.Field(), "tag": fe.Tag()}
		}
		return invalid, nil
	}

	return nil, err
}

func ExtractDateRange(date_range string) (start time.Time, end time.Time, err error) {
	dates := strings.Split(date_range, "to")
	if len(dates) != 2 {
		err = errors.New("invalid date range")
		return
	}

	start, err1 := ISODateConversion(dates[0])
	end, err = ISODateConversion(dates[1])
	if err1 != nil || err != nil {
		err = errors.New("invalid date format")
		return
	}
	return
}

func Request(request string, headers map[string][]string, urlPath string, method string) (string, error) {

	reqURL, _ := url.Parse(urlPath)

	reqBody := io.NopCloser(strings.NewReader(request))

	req := &http.Request{
		Method: method,
		URL:    reqURL,
		Header: headers,
		Body:   reqBody,
	}

	res, err := ExternalRequestTimer(req)
	if err != nil {
		logrus.Errorf("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | ERROR : %v", urlPath, method, request, err)
		return "", err
	}

	data, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	resbody := string(data)

	logrus.Infof("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d | RESPONSE : %s", urlPath, method, request, res.Status, res.StatusCode, resbody)

	if res.StatusCode > 299 || res.StatusCode <= 199 {
		logrus.Errorf("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, request, res.Status, res.StatusCode)
		return resbody, fmt.Errorf("%d", res.StatusCode)
	}

	// logrus.Infof("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, resbody, res.Status, res.StatusCode)

	return resbody, nil
}

func ExternalRequestTimer(req *http.Request) (*http.Response, error) {

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			dns = time.Now()
			logrus.Debug(dsi)
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			logrus.Debug(ddi)
			logrus.Infof("DNS Done: %v", time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			// logrus.Debug(cs, err)
			logrus.Infof("TLS Handshake: %v", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) {
			connect = time.Now()
			logrus.Debug(network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			logrus.Debug(network, addr, err)
			logrus.Infof("Connect time: %v", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			logrus.Warnf("TAT : %v", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	// NOTE: Below line is to ignore ssl certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func SendEmail(recipients []string, subject, contentType, content string, attachments []models.InMemFile) error {
	smtpUser := os.Getenv("EMAIL_SENDER")
	smtpPass := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMT_HOST")
	smtpPort := os.Getenv("SMT_PORT")
	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		return errors.New("error converting string to int")
	}

	// Create the MIME email headers
	m := mail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, content)

	// Attach the file if provided
	for _, attachment := range attachments {
		m.Attach(attachment.FileName, mail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Buffer.Bytes())
			if err != nil {
				logrus.Error(err)
			}
			return err
		}))
	}

	// Set up the SMTP server details
	d := mail.NewDialer(smtpHost, smtpPortInt, smtpUser, smtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email
	logrus.Infoln("Starting to send email...")
	if err := d.DialAndSend(m); err != nil {
		logrus.Errorf("An error occurred while sending the email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logrus.Infoln("Email sent successfully!")
	return nil
}

// custom validator: fails if the field is all whitespace.
func NotBlank(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return strings.TrimSpace(s) != ""
}
