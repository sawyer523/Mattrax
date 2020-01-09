package soap

import (
	"net/http"
	"strconv"

	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/rs/zerolog/log"
)

// FaultEnvelop is the response body used to report an error to the device.
type FaultEnvelop struct {
	XMLName                      xml.Name                     `xml:"s:Envelope"`
	NamespaceS                   string                       `xml:"xmlns:s,attr"`
	NamespaceA                   string                       `xml:"xmlns:a,attr"`
	Value                        string                       `xml:"s:Body>s:Fault>s:Code>s:Value"`
	Subcode                      string                       `xml:"s:Body>s:Fault>s:Code>s:Subcode>s:Value"`
	Reason                       FaultReason                  `xml:"s:Body>s:Fault>s:Reason>s:Text"`
	DeviceEnrollmentServiceError DeviceEnrollmentServiceError `xml:"s:Body>s:Fault>s:Detail>DeviceEnrollmentServiceError,omitempty"` // TODO: Remove s:Details when this not set + ,omitempty not working on TraceID
}

// Response marshals and sends the SOAP fault to the HTTP response
// It verifies the values in the fault to prevent incorrect values.
// If it is unable to respond to the client it logs the error (error isn't handled because nothing can be done)
func (fault FaultEnvelop) Response(w http.ResponseWriter) {
	faultLogger := log.With().Str("value", fault.Value).Str("subcode", fault.Subcode).Str("reason", fault.Reason.Text).Str("error-type", fault.DeviceEnrollmentServiceError.Type).Str("message", fault.DeviceEnrollmentServiceError.Message).Str("traceID", fault.DeviceEnrollmentServiceError.TraceID).Logger()

	if fault.Value != "s:Sender" && fault.Value != "s:Receiver" {
		log.Error().Msg("|||||||||||||||||||||||| DEVELOPER BUG ||||||||||||||||||||||||")
		faultLogger.Error().Bool("developer-bug-please-report", true).Msg("error: fault.Respond(): invalid value")
		fault.Value = "s:Receiver"
	}

	if fault.Subcode != "s:MessageFormat" && fault.Subcode != "s:Authentication" && fault.Subcode != "s:Authorization" && fault.Subcode != "s:CertificateRequest" && fault.Subcode != "s:EnrollmentServer" && fault.Subcode != "a:InternalServiceFault" && fault.Subcode != "a:InvalidSecurity" && fault.Subcode != "a:ActionMismatch" && fault.Subcode != "a:EndpointUnavailable" {
		log.Error().Msg("|||||||||||||||||||||||| DEVELOPER BUG ||||||||||||||||||||||||")
		faultLogger.Error().Bool("developer-bug-please-report", true).Msg("error: fault.Respond(): invalid subcode")
		fault.Subcode = "s:EnrollmentServer"
	}

	if fault.DeviceEnrollmentServiceError.Type != "" && fault.DeviceEnrollmentServiceError.Type != "DeviceCapReached" && fault.DeviceEnrollmentServiceError.Type != "DeviceNotSupported" && fault.DeviceEnrollmentServiceError.Type != "NotSupported" && fault.DeviceEnrollmentServiceError.Type != "NotEligibleToRenew" && fault.DeviceEnrollmentServiceError.Type != "InMaintenance" && fault.DeviceEnrollmentServiceError.Type != "UserLicense" && fault.DeviceEnrollmentServiceError.Type != "InvalidEnrollmentData" && fault.DeviceEnrollmentServiceError.Type != "EnrollmentServer" {
		log.Error().Msg("|||||||||||||||||||||||| DEVELOPER BUG ||||||||||||||||||||||||")
		faultLogger.Error().Bool("developer-bug-please-report", true).Msg("error: fault.Respond(): invalid device enrollment service error type")
		fault.DeviceEnrollmentServiceError.Type = ""
		fault.DeviceEnrollmentServiceError.Message = ""
	}

	response, err := xml.Marshal(fault)
	if err != nil {
		_, err := w.Write([]byte("HTTP 500: Internal Service Fault"))
		if err != nil {
			faultLogger.Error().Err(err).Msg("error: fault.Respond(): failed to send error response after xml marshaling failed")
		}
		return
	}

	w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))

	if fault.Value == "s:Sender" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(response)
	if err != nil {
		faultLogger.Error().Err(err).Msg("error: fault.Respond(): failed to send fault response body")
	}
}

// FaultReason contains the language and reason that caused the fault
type FaultReason struct {
	Language string `xml:"xml:lang,attr"`
	Text     string `xml:",innerxml"`
}

// DeviceEnrollmentServiceError contains specific details about the fault
type DeviceEnrollmentServiceError struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment DeviceEnrollmentServiceError"`
	Type    string   `xml:"ErrorType"`
	Message string   `xml:"Message"`
	TraceID string   `xml:"TraceID"`
}

// NewBasicFault is an simple helper to create a new FaultEnvelop
// It fills all static values exposing only the values you need access to.
func NewBasicFault(value string, subcode string, reason string) FaultEnvelop {
	return FaultEnvelop{
		NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
		NamespaceA: "http://www.w3.org/2005/08/addressing",
		Value:      value,
		Subcode:    subcode,
		Reason: FaultReason{
			Language: "en-US",
			Text:     reason,
		},
	}
}

// NewEnrollmentFault is an simple helper to create a new FaultEnvelop
// It fills all static values exposing only the values you need access to.
func NewEnrollmentFault(value string, subcode string, reason string, errorType string, message string, traceID string) FaultEnvelop {
	return FaultEnvelop{
		NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
		NamespaceA: "http://www.w3.org/2005/08/addressing",
		Value:      value,
		Subcode:    subcode,
		Reason: FaultReason{
			Language: "en-US",
			Text:     reason,
		},
		DeviceEnrollmentServiceError: DeviceEnrollmentServiceError{
			Type:    errorType,
			Message: message,
			TraceID: traceID,
		},
	}
}
