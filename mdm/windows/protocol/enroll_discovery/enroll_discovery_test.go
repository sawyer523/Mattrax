package enrolldiscovery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/matryer/is"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// TODO: Tests should be routed through Mux from NewMockServer response (this will tests mounted url due to it being hardcoded by device)

func TestDiscoveryGET(t *testing.T) {
	is := is.New(t)

	req, err := http.NewRequest("GET", "/EnrollmentServer/Discovery.svc", nil)
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	GETHandler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusOK) // Request should response status OK
	is.Equal(res.Body.Len(), 0)       // Body should be empty
}

func TestDiscoveryPOST(t *testing.T) {
	is := is.New(t)

	body := []byte(`<s:Envelope xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:s="http://www.w3.org/2003/05/soap-envelope"><s:Header><a:Action s:mustUnderstand="1">http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover</a:Action><a:MessageID>urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><a:To s:mustUnderstand="1">https://EnterpriseEnrollment.otbeaumont.me:443/EnrollmentServer/Discovery.svc</a:To></s:Header><s:Body><Discover xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment"><request xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><EmailAddress>oscar@otbeaumont.me</EmailAddress><RequestVersion>4.0</RequestVersion><DeviceType>CIMClient_Windows</DeviceType><ApplicationVersion>10.0.18362.0</ApplicationVersion><OSEdition>48</OSEdition><AuthPolicies><AuthPolicy>OnPremise</AuthPolicy><AuthPolicy>Federated</AuthPolicy></AuthPolicies></request></Discover></s:Body></s:Envelope>`)
	req, err := http.NewRequest("POST", "/EnrollmentServer/Discovery.svc", bytes.NewBuffer(body))
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	Handler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusOK) // Request should response status OK
	is.True(res.Body.Len() != 0)      // Body should not be empty

	var cmd Response
	err = xml.NewDecoder(res.Body).Decode(&cmd)
	is.NoErr(err) // Error decoding response body

	is.Equal(cmd.Header.Action.Value, "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse")                                   // Invalid schema in response
	is.True(cmd.Header.ActivityID != "")                                                                                                                                         // ActivityID must be set in response
	is.Equal(cmd.Header.RelatesTo, "urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478")                                                                                              // RelatesTo must match the MessageID of the request
	is.True(cmd.Body.DiscoverResponse.AuthPolicy == "OnPremise" || cmd.Body.DiscoverResponse.AuthPolicy == "Federated" || cmd.Body.DiscoverResponse.AuthPolicy == "Certificate") // AuthPolicy must contain a valid value

	url1, err := url.ParseRequestURI(cmd.Body.DiscoverResponse.EnrollmentPolicyServiceURL)
	is.NoErr(err)
	is.Equal(url1.Scheme, "https")

	url2, err := url.ParseRequestURI(cmd.Body.DiscoverResponse.EnrollmentServiceURL)
	is.NoErr(err)
	is.Equal(url2.Scheme, "https")

	if cmd.Body.DiscoverResponse.AuthPolicy == "Federated" {
		url3, err := url.ParseRequestURI(cmd.Body.DiscoverResponse.AuthenticationServiceURL)
		is.NoErr(err)
		is.Equal(url3.Scheme, "https")
	}
}

func TestDiscoveryPOST_InvalidHeaderToAddr(t *testing.T) {
	is := is.New(t)

	body := []byte(`<s:Envelope xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:s="http://www.w3.org/2003/05/soap-envelope"><s:Header><a:Action s:mustUnderstand="1">http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover</a:Action><a:MessageID>urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><a:To s:mustUnderstand="1">this-is-not-a-url</a:To></s:Header><s:Body><Discover xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment"><request xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><EmailAddress>oscar@otbeaumont.me</EmailAddress><RequestVersion>4.0</RequestVersion><DeviceType>CIMClient_Windows</DeviceType><ApplicationVersion>10.0.18362.0</ApplicationVersion><OSEdition>48</OSEdition><AuthPolicies><AuthPolicy>OnPremise</AuthPolicy><AuthPolicy>Federated</AuthPolicy></AuthPolicies></request></Discover></s:Body></s:Envelope>`)
	req, err := http.NewRequest("POST", "/EnrollmentServer/Discovery.svc", bytes.NewBuffer(body))
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	Handler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusBadRequest) // Request should response status BadRequest
	is.True(res.Body.Len() != 0)              // Body should not be empty

	var cmd soap.FaultEnvelop
	err = xml.NewDecoder(res.Body).Decode(&cmd)
	is.NoErr(err) // Error decoding response body

	is.Equal(cmd.Value, "s:Sender")
	is.Equal(cmd.Subcode, "a:EndpointUnavailable")
	is.True(cmd.Reason.Text != "")
}

func TestDiscoveryPOST_InvalidBodyEmailAddr(t *testing.T) {
	is := is.New(t)

	body := []byte(`<s:Envelope xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:s="http://www.w3.org/2003/05/soap-envelope"><s:Header><a:Action s:mustUnderstand="1">http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover</a:Action><a:MessageID>urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><a:To s:mustUnderstand="1">https://EnterpriseEnrollment.otbeaumont.me:443/EnrollmentServer/Discovery.svc</a:To></s:Header><s:Body><Discover xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment"><request xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><EmailAddress>not-an-email-address</EmailAddress><RequestVersion>4.0</RequestVersion><DeviceType>CIMClient_Windows</DeviceType><ApplicationVersion>10.0.18362.0</ApplicationVersion><OSEdition>48</OSEdition><AuthPolicies><AuthPolicy>OnPremise</AuthPolicy><AuthPolicy>Federated</AuthPolicy></AuthPolicies></request></Discover></s:Body></s:Envelope>`)
	req, err := http.NewRequest("POST", "/EnrollmentServer/Discovery.svc", bytes.NewBuffer(body))
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	Handler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusBadRequest) // Request should response status BadRequest
	is.True(res.Body.Len() != 0)              // Body should not be empty

	var cmd soap.FaultEnvelop
	err = xml.NewDecoder(res.Body).Decode(&cmd)
	is.NoErr(err) // Error decoding response body

	is.Equal(cmd.Value, "s:Sender")
	is.Equal(cmd.Subcode, "s:MessageFormat")
	is.True(cmd.Reason.Text != "")
}

func TestDiscoveryPOST_UnsupportedBodyRequestVersion(t *testing.T) {
	is := is.New(t)

	body := []byte(`<s:Envelope xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:s="http://www.w3.org/2003/05/soap-envelope"><s:Header><a:Action s:mustUnderstand="1">http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover</a:Action><a:MessageID>urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><a:To s:mustUnderstand="1">https://EnterpriseEnrollment.otbeaumont.me:443/EnrollmentServer/Discovery.svc</a:To></s:Header><s:Body><Discover xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment"><request xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><EmailAddress>oscar@otbeaumont.me</EmailAddress><RequestVersion>3.0</RequestVersion><DeviceType>CIMClient_Windows</DeviceType><ApplicationVersion>10.0.18362.0</ApplicationVersion><OSEdition>48</OSEdition><AuthPolicies><AuthPolicy>OnPremise</AuthPolicy><AuthPolicy>Federated</AuthPolicy></AuthPolicies></request></Discover></s:Body></s:Envelope>`)
	req, err := http.NewRequest("POST", "/EnrollmentServer/Discovery.svc", bytes.NewBuffer(body))
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	Handler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusBadRequest) // Request should response status BadRequest
	is.True(res.Body.Len() != 0)              // Body should not be empty

	var cmd soap.FaultEnvelop
	err = xml.NewDecoder(res.Body).Decode(&cmd)
	is.NoErr(err) // Error decoding response body

	is.Equal(cmd.Value, "s:Sender")
	is.Equal(cmd.Subcode, "a:InternalServiceFault")
	is.True(cmd.Reason.Text != "")
}

func TestDiscoveryPOST_UnsupportedBodyDeviceType(t *testing.T) {
	is := is.New(t)

	body := []byte(`<s:Envelope xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:s="http://www.w3.org/2003/05/soap-envelope"><s:Header><a:Action s:mustUnderstand="1">http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover</a:Action><a:MessageID>urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><a:To s:mustUnderstand="1">https://EnterpriseEnrollment.otbeaumont.me:443/EnrollmentServer/Discovery.svc</a:To></s:Header><s:Body><Discover xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment"><request xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><EmailAddress>oscar@otbeaumont.me</EmailAddress><RequestVersion>4.0</RequestVersion><DeviceType>NONEXISTANT_Device</DeviceType><ApplicationVersion>10.0.18362.0</ApplicationVersion><OSEdition>48</OSEdition><AuthPolicies><AuthPolicy>OnPremise</AuthPolicy><AuthPolicy>Federated</AuthPolicy></AuthPolicies></request></Discover></s:Body></s:Envelope>`)
	req, err := http.NewRequest("POST", "/EnrollmentServer/Discovery.svc", bytes.NewBuffer(body))
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	Handler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusBadRequest) // Request should response status BadRequest
	is.True(res.Body.Len() != 0)              // Body should not be empty

	var cmd soap.FaultEnvelop
	err = xml.NewDecoder(res.Body).Decode(&cmd)
	is.NoErr(err) // Error decoding response body

	is.Equal(cmd.Value, "s:Sender")
	is.Equal(cmd.Subcode, "a:InternalServiceFault")
	is.True(cmd.Reason.Text != "")
}

func TestDiscoveryPOST_MismatchedAuthPolicies(t *testing.T) {
	is := is.New(t)

	body := []byte(`<s:Envelope xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:s="http://www.w3.org/2003/05/soap-envelope"><s:Header><a:Action s:mustUnderstand="1">http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover</a:Action><a:MessageID>urn:uuid:748132ec-a575-4329-b01b-6171a9cf8478</a:MessageID><a:ReplyTo><a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address></a:ReplyTo><a:To s:mustUnderstand="1">https://EnterpriseEnrollment.otbeaumont.me:443/EnrollmentServer/Discovery.svc</a:To></s:Header><s:Body><Discover xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment"><request xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><EmailAddress>oscar@otbeaumont.me</EmailAddress><RequestVersion>4.0</RequestVersion><DeviceType>CIMClient_Windows</DeviceType><ApplicationVersion>10.0.18362.0</ApplicationVersion><OSEdition>48</OSEdition><AuthPolicies><AuthPolicy>OnPremise</AuthPolicy></AuthPolicies></request></Discover></s:Body></s:Envelope>`)
	req, err := http.NewRequest("POST", "/EnrollmentServer/Discovery.svc", bytes.NewBuffer(body))
	is.NoErr(err) // Error creating mock request

	res := httptest.NewRecorder()
	Handler(mattrax.NewMockServer(t))(res, req)

	is.Equal(res.Code, http.StatusBadRequest) // Request should response status BadRequest
	is.True(res.Body.Len() != 0)              // Body should not be empty

	var cmd soap.FaultEnvelop
	err = xml.NewDecoder(res.Body).Decode(&cmd)
	is.NoErr(err) // Error decoding response body

	is.Equal(cmd.Value, "s:Sender")
	is.Equal(cmd.Subcode, "a:InternalServiceFault")
	is.True(cmd.Reason.Text != "")
}
