package sendwithus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/superhuman/backend/lib/errors"
)

const (
	swuEndpoint     = "https://api.sendwithus.com/api/v1"
	apiHeaderClient = "golang-0.0.1"
)

// SWUClient implements a SendWithUs client.
type SWUClient struct {
	Client *http.Client
	apiKey string
	URL    string
}

// SWUTemplate describes a SendWithUs template.
type SWUTemplate struct {
	ID       string        `json:"id,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Created  int64         `json:"created,omitempty"`
	Versions []*SWUVersion `json:"versions,omitempty"`
	Name     string        `json:"name,omitempty"`
}

// SWUVersion describes a SendWithUs version.
type SWUVersion struct {
	Name      string `json:"name,omitempty"`
	ID        string `json:"id,omitempty"`
	Created   int64  `json:"created,omitempty"`
	HTML      string `json:"html,omitempty"`
	Text      string `json:"text,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Published bool   `json:"published,omitempty"`
}

// SWUEmail describes a SendWithUs email.
type SWUEmail struct {
	ID          string            `json:"email_id,omitempty"`
	Recipient   *SWURecipient     `json:"recipient,omitempty"`
	CC          []*SWURecipient   `json:"cc,omitempty"`
	BCC         []*SWURecipient   `json:"bcc,omitempty"`
	Sender      *SWUSender        `json:"sender,omitempty"`
	EmailData   map[string]string `json:"email_data,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Inline      *SWUAttachment    `json:"inline,omitempty"`
	Files       []*SWUAttachment  `json:"files,omitempty"`
	ESPAccount  string            `json:"esp_account,omitempty"`
	VersionName string            `json:"version_name,omitempty"`
}

// SWUDripCampaign describes a SendWithUs drip campaign.
type SWUDripCampaign struct {
	Recipient  *SWURecipient     `json:"recipient,omitempty"`
	CC         []*SWURecipient   `json:"cc,omitempty"`
	BCC        []*SWURecipient   `json:"bcc,omitempty"`
	Sender     *SWUSender        `json:"sender,omitempty"`
	EmailData  map[string]string `json:"email_data,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	ESPAccount string            `json:"esp_account,omitempty"`
	Locale     string            `json:"locale,omitempty"`
}

// SWURecipient describes a SendWithUs recipient.
type SWURecipient struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

// SWUSender describes a SendWithUs sender.
type SWUSender struct {
	SWURecipient
	ReplyTo string `json:"reply_to,omitempty"`
}

// SWUAttachment describes a SendWithUs attachment.
type SWUAttachment struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

// SWULogEvent describes a SendWithUs log event.
type SWULogEvent struct {
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

// SWULogQuery describes a SendWithUs log query.
type SWULogQuery struct {
	Count      int   `json:"count,omitempty" url:"count,omitempty"`
	Offset     int   `json:"offset,omitempty" url:"offset,omitempty"`
	CreatedGT  int64 `json:"created_gt,omitempty" url:"created_gt,omitempty"`
	CreatedGTE int64 `json:"created_gte,omitempty" url:"created_gte,omitempty"`
	CreatedLT  int64 `json:"created_lt,omitempty" url:"created_lt,omitempty"`
	CreatedLTE int64 `json:"created_lte,omitempty" url:"created_lte,omitempty"`
}

// SWULog describes a SendWithUs log.
type SWULog struct {
	SWULogEvent
	ID               string `json:"id,omitempty"`
	RecipientName    string `json:"recipient_name,omitempty"`
	RecipientAddress string `json:"recipient_address,omitempty"`
	Status           string `json:"status,omitempty"`
	EmailID          string `json:"email_id,omitempty"`
	EmailName        string `json:"email_name,omitempty"`
	EmailVersion     string `json:"email_version,omitempty"`
	EventsURL        string `json:"events_url,omitempty"`
}

// SWULogResend describes a SendWithUs log resend.
type SWULogResend struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	ID      string `json:"log_id,omitempty"`
	Email   struct {
		Name        string `json:"name"`
		VersionName string `json:"version_name"`
	} `json:"email"`
}

// SWUError describes a SendWithUs error.
type SWUError struct {
	Code    int
	Message string
}

func newSWUError(res *http.Response, message string) *SWUError {
	s := &SWUError{
		Message: message,
	}
	if res != nil {
		s.Code = res.StatusCode
	}
	return s
}

// Error implements the error interface.
func (e *SWUError) Error() string {
	return fmt.Sprintf("swu.go: Status code: %d, Error: %s", e.Code, e.Message)
}

// New initializes a new SWUClient.
func New(apiKey string) *SWUClient {
	return &SWUClient{
		Client: http.DefaultClient,
		apiKey: apiKey,
		URL:    swuEndpoint,
	}
}

// Templates executes a SendWithUs api call.
func (c *SWUClient) Templates() ([]*SWUTemplate, error) {
	return c.Emails()
}

// Emails executes a SendWithUs api call.
func (c *SWUClient) Emails() ([]*SWUTemplate, error) {
	var parse []*SWUTemplate
	if err := c.makeRequest("GET", "/templates", nil, &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return parse, nil
}

// GetTemplate executes a SendWithUs api call.
func (c *SWUClient) GetTemplate(id string) (*SWUTemplate, error) {
	var parse SWUTemplate
	if err := c.makeRequest("GET", "/templates/"+id, nil, &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// GetTemplateVersion executes a SendWithUs api call.
func (c *SWUClient) GetTemplateVersion(id, version string) (*SWUVersion, error) {
	var parse SWUVersion
	if err := c.makeRequest("GET", "/templates/"+id+"/versions/"+version, nil, &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// UpdateTemplateVersion executes a SendWithUs api call.
func (c *SWUClient) UpdateTemplateVersion(id, version string, template *SWUVersion) (*SWUVersion, error) {
	var parse SWUVersion
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := c.makeRequest("PUT", "/templates/"+id+"/versions/"+version, bytes.NewReader(payload), &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// CreateTemplate executes a SendWithUs api call.
func (c *SWUClient) CreateTemplate(template *SWUVersion) (*SWUTemplate, error) {
	var parse SWUTemplate
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := c.makeRequest("POST", "/templates", bytes.NewReader(payload), &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// CreateTemplateVersion executes a SendWithUs api call.
func (c *SWUClient) CreateTemplateVersion(id string, template *SWUVersion) (*SWUTemplate, error) {
	var parse SWUTemplate
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := c.makeRequest("POST", "/templates/"+id+"/versions", bytes.NewReader(payload), &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// Send executes a SendWithUs api call.
func (c *SWUClient) Send(email *SWUEmail) error {
	payload, err := json.Marshal(email)
	if err != nil {
		return errors.Wrap(err)
	}
	if err := c.makeRequest("POST", "/send", bytes.NewReader(payload), nil); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// ActivateDripCampaign executes a SendWithUs api call.
func (c *SWUClient) ActivateDripCampaign(id string, dripCampaign *SWUDripCampaign) error {
	payload, err := json.Marshal(dripCampaign)
	if err != nil {
		return errors.Wrap(err)
	}
	if err := c.makeRequest("POST", "/drip_campaigns/"+id+"/activate", bytes.NewReader(payload), nil); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// GetLogs executes a SendWithUs api call.
func (c *SWUClient) GetLogs(q *SWULogQuery) ([]*SWULog, error) {
	var parse []*SWULog
	payload, err := query.Values(q)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := c.makeRequest("GET", "/logs?"+payload.Encode(), nil, &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return parse, nil
}

// GetLog executes a SendWithUs api call.
func (c *SWUClient) GetLog(id string) (*SWULog, error) {
	var parse SWULog
	if err := c.makeRequest("GET", "/logs/"+id, nil, &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// GetLogEvents executes a SendWithUs api call.
func (c *SWUClient) GetLogEvents(id string) (*SWULogEvent, error) {
	var parse SWULogEvent
	if err := c.makeRequest("GET", "/logs/"+id+"/events", nil, &parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return &parse, nil
}

// ResendLog executes a SendWithUs api call.
func (c *SWUClient) ResendLog(id string) (*SWULogResend, error) {
	parse := &SWULogResend{
		ID: id,
	}
	payload, err := json.Marshal(parse)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := c.makeRequest("POST", "/resend", bytes.NewReader(payload), parse); err != nil {
		return nil, errors.Wrap(err)
	}
	return parse, nil
}

func (c *SWUClient) makeRequest(method, endpoint string, body io.Reader, result interface{}) error {
	r, err := http.NewRequest(method, c.URL+endpoint, body)
	if err != nil {
		return errors.Wrap(err)
	}
	r.SetBasicAuth(c.apiKey, "")
	r.Header.Set("X-SWU-API-CLIENT", apiHeaderClient)
	res, err := c.Client.Do(r)
	if err != nil {
		return errors.Wrap(newSWUError(res, err.Error()))
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(newSWUError(res, err.Error()))
	}
	if res.StatusCode >= 300 {
		return errors.Wrap(newSWUError(res, string(b)))
	}
	if result != nil {
		return buildRespJSON(b, result)
	}
	return nil
}

func buildRespJSON(b []byte, parse interface{}) error {
	if err := json.Unmarshal(b, parse); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
