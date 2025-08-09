package fcm_notification

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// FCMClient represents the Firebase Cloud Messaging client
type FCMClient struct {
	projectID   string
	endpoint    string
	client      *http.Client
	credentials *ServiceAccountCredentials
	tokenCache  *tokenCache
}

// ServiceAccountCredentials represents Firebase service account credentials
type ServiceAccountCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// tokenCache manages OAuth2 access tokens
type tokenCache struct {
	token     string
	expiresAt time.Time
}

// NewFCMClientFromFile creates a new FCM client from a service account file
func NewFCMClientFromFile(credentialsPath string) (*FCMClient, error) {
	credentialsData, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	return NewFCMClientFromJSON(credentialsData)
}

// NewFCMClientFromJSON creates a new FCM client from service account JSON
func NewFCMClientFromJSON(credentialsJSON []byte) (*FCMClient, error) {
	var credentials ServiceAccountCredentials
	if err := json.Unmarshal(credentialsJSON, &credentials); err != nil {
		return nil, fmt.Errorf("failed to parse credentials JSON: %w", err)
	}

	return &FCMClient{
		projectID:   credentials.ProjectID,
		endpoint:    fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", credentials.ProjectID),
		credentials: &credentials,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenCache: &tokenCache{},
	}, nil
}

// NewFCMClient creates a new FCM client with service account credentials
func NewFCMClient(credentials *ServiceAccountCredentials) *FCMClient {
	return &FCMClient{
		projectID:   credentials.ProjectID,
		endpoint:    fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", credentials.ProjectID),
		credentials: credentials,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenCache: &tokenCache{},
	}
}

// NotificationPayload represents the notification content
type NotificationPayload struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Image string `json:"image,omitempty"`
}

// AndroidConfig represents Android-specific options
type AndroidConfig struct {
	CollapseKey           string               `json:"collapse_key,omitempty"`
	Priority              string               `json:"priority,omitempty"`
	TTL                   string               `json:"ttl,omitempty"`
	RestrictedPackageName string               `json:"restricted_package_name,omitempty"`
	Data                  map[string]string    `json:"data,omitempty"`
	Notification          *AndroidNotification `json:"notification,omitempty"`
	FCMOptions            *AndroidFCMOptions   `json:"fcm_options,omitempty"`
}

// AndroidNotification represents Android notification options
type AndroidNotification struct {
	Title                 string   `json:"title,omitempty"`
	Body                  string   `json:"body,omitempty"`
	Icon                  string   `json:"icon,omitempty"`
	Color                 string   `json:"color,omitempty"`
	Sound                 string   `json:"sound,omitempty"`
	Tag                   string   `json:"tag,omitempty"`
	ClickAction           string   `json:"click_action,omitempty"`
	BodyLocKey            string   `json:"body_loc_key,omitempty"`
	BodyLocArgs           []string `json:"body_loc_args,omitempty"`
	TitleLocKey           string   `json:"title_loc_key,omitempty"`
	TitleLocArgs          []string `json:"title_loc_args,omitempty"`
	ChannelID             string   `json:"channel_id,omitempty"`
	Ticker                string   `json:"ticker,omitempty"`
	Sticky                bool     `json:"sticky,omitempty"`
	EventTime             string   `json:"event_time,omitempty"`
	LocalOnly             bool     `json:"local_only,omitempty"`
	NotificationPriority  string   `json:"notification_priority,omitempty"`
	DefaultSound          bool     `json:"default_sound,omitempty"`
	DefaultVibrateTimings bool     `json:"default_vibrate_timings,omitempty"`
	DefaultLightSettings  bool     `json:"default_light_settings,omitempty"`
	VibrateTimings        []string `json:"vibrate_timings,omitempty"`
	Visibility            string   `json:"visibility,omitempty"`
	NotificationCount     int      `json:"notification_count,omitempty"`
}

// AndroidFCMOptions represents Android FCM options
type AndroidFCMOptions struct {
	AnalyticsLabel string `json:"analytics_label,omitempty"`
}

// APNSConfig represents Apple Push Notification service options
type APNSConfig struct {
	Headers    map[string]string `json:"headers,omitempty"`
	Payload    *APNSPayload      `json:"payload,omitempty"`
	FCMOptions *APNSFCMOptions   `json:"fcm_options,omitempty"`
}

// APNSPayload represents APNS payload
type APNSPayload struct {
	Aps        *APSPayload            `json:"aps,omitempty"`
	CustomData map[string]interface{} `json:"-"`
}

// APSPayload represents APS payload
type APSPayload struct {
	Alert            interface{} `json:"alert,omitempty"`
	Badge            interface{} `json:"badge,omitempty"`
	Sound            interface{} `json:"sound,omitempty"`
	ContentAvailable int         `json:"content-available,omitempty"`
	MutableContent   int         `json:"mutable-content,omitempty"`
	Category         string      `json:"category,omitempty"`
	ThreadID         string      `json:"thread-id,omitempty"`
}

// APNSFCMOptions represents APNS FCM options
type APNSFCMOptions struct {
	AnalyticsLabel string `json:"analytics_label,omitempty"`
	Image          string `json:"image,omitempty"`
}

// WebpushConfig represents Webpush options
type WebpushConfig struct {
	Headers      map[string]string    `json:"headers,omitempty"`
	Data         map[string]string    `json:"data,omitempty"`
	Notification *WebpushNotification `json:"notification,omitempty"`
	FCMOptions   *WebpushFCMOptions   `json:"fcm_options,omitempty"`
}

// WebpushNotification represents Webpush notification
type WebpushNotification struct {
	Title              string                 `json:"title,omitempty"`
	Body               string                 `json:"body,omitempty"`
	Icon               string                 `json:"icon,omitempty"`
	Image              string                 `json:"image,omitempty"`
	Badge              string                 `json:"badge,omitempty"`
	Tag                string                 `json:"tag,omitempty"`
	Data               map[string]interface{} `json:"data,omitempty"`
	Dir                string                 `json:"dir,omitempty"`
	Lang               string                 `json:"lang,omitempty"`
	Renotify           bool                   `json:"renotify,omitempty"`
	RequireInteraction bool                   `json:"requireInteraction,omitempty"`
	Silent             bool                   `json:"silent,omitempty"`
	Timestamp          int64                  `json:"timestamp,omitempty"`
	Vibrate            []int                  `json:"vibrate,omitempty"`
}

// WebpushFCMOptions represents Webpush FCM options
type WebpushFCMOptions struct {
	Link           string `json:"link,omitempty"`
	AnalyticsLabel string `json:"analytics_label,omitempty"`
}

// FCMMessage represents the complete FCM message structure (v1 API)
type FCMMessage struct {
	Name         string               `json:"name,omitempty"`
	Data         map[string]string    `json:"data,omitempty"`
	Notification *NotificationPayload `json:"notification,omitempty"`
	Android      *AndroidConfig       `json:"android,omitempty"`
	Webpush      *WebpushConfig       `json:"webpush,omitempty"`
	APNS         *APNSConfig          `json:"apns,omitempty"`
	FCMOptions   *FCMOptions          `json:"fcm_options,omitempty"`
	Token        string               `json:"token,omitempty"`
	Topic        string               `json:"topic,omitempty"`
	Condition    string               `json:"condition,omitempty"`
}

// BatchResponse represents response for batch operations
type BatchResponse struct {
	SuccessCount int             `json:"success_count"`
	FailureCount int             `json:"failure_count"`
	Responses    []TokenResponse `json:"responses"`
}

// TokenResponse represents response for individual token
type TokenResponse struct {
	Token    string       `json:"token"`
	Success  bool         `json:"success"`
	Response *FCMResponse `json:"response,omitempty"`
	Error    error        `json:"error,omitempty"`
}

// FCMOptions represents FCM options
type FCMOptions struct {
	AnalyticsLabel string `json:"analytics_label,omitempty"`
}

// FCMRequest wraps the message for the v1 API
type FCMRequest struct {
	ValidateOnly bool        `json:"validate_only,omitempty"`
	Message      *FCMMessage `json:"message"`
}

// FCMResponse represents the response from FCM server (v1 API)
type FCMResponse struct {
	Name  string    `json:"name,omitempty"`
	Error *FCMError `json:"error,omitempty"`
	Token string    `json:"token,omitempty"`
}

// FCMError represents FCM error response
type FCMError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NotificationBuilder helps build notifications with a fluent interface
type NotificationBuilder struct {
	message *FCMMessage
}

// NewNotificationBuilder creates a new notification builder
func NewNotificationBuilder() *NotificationBuilder {
	return &NotificationBuilder{
		message: &FCMMessage{},
	}
}

// ToToken sets the target device token
func (nb *NotificationBuilder) ToToken(token string) *NotificationBuilder {
	nb.message.Token = token
	return nb
}

// ToTopic sets the target topic
func (nb *NotificationBuilder) ToTopic(topic string) *NotificationBuilder {
	nb.message.Topic = topic
	return nb
}

// ToCondition sets a condition for targeting
func (nb *NotificationBuilder) ToCondition(condition string) *NotificationBuilder {
	nb.message.Condition = condition
	return nb
}

// WithNotification sets the notification payload
func (nb *NotificationBuilder) WithNotification(title, body string) *NotificationBuilder {
	nb.message.Notification = &NotificationPayload{
		Title: title,
		Body:  body,
	}
	return nb
}

// WithNotificationAndImage sets the notification payload with image
func (nb *NotificationBuilder) WithNotificationAndImage(title, body, image string) *NotificationBuilder {
	nb.message.Notification = &NotificationPayload{
		Title: title,
		Body:  body,
		Image: image,
	}
	return nb
}

// WithData sets custom data payload
func (nb *NotificationBuilder) WithData(data map[string]string) *NotificationBuilder {
	nb.message.Data = data
	return nb
}

// WithAndroidConfig sets Android-specific configuration
func (nb *NotificationBuilder) WithAndroidConfig(config *AndroidConfig) *NotificationBuilder {
	nb.message.Android = config
	return nb
}

// WithAPNSConfig sets APNS-specific configuration
func (nb *NotificationBuilder) WithAPNSConfig(config *APNSConfig) *NotificationBuilder {
	nb.message.APNS = config
	return nb
}

// WithWebpushConfig sets Webpush-specific configuration
func (nb *NotificationBuilder) WithWebpushConfig(config *WebpushConfig) *NotificationBuilder {
	nb.message.Webpush = config
	return nb
}

// WithAnalyticsLabel sets analytics label for tracking
func (nb *NotificationBuilder) WithAnalyticsLabel(label string) *NotificationBuilder {
	if nb.message.FCMOptions == nil {
		nb.message.FCMOptions = &FCMOptions{}
	}
	nb.message.FCMOptions.AnalyticsLabel = label
	return nb
}

// Build returns the built FCM message
func (nb *NotificationBuilder) Build() *FCMMessage {
	return nb.message
}

// getAccessToken generates or retrieves cached OAuth2 access token
func (c *FCMClient) getAccessToken(ctx context.Context) (string, error) {
	// Check if we have a valid cached token
	if c.tokenCache.token != "" && time.Now().Before(c.tokenCache.expiresAt) {
		return c.tokenCache.token, nil
	}

	// Generate new JWT token
	token, err := c.generateJWT()
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Exchange JWT for access token
	accessToken, expiresIn, err := c.exchangeJWTForAccessToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed to exchange JWT for access token: %w", err)
	}

	// Cache the token
	c.tokenCache.token = accessToken
	c.tokenCache.expiresAt = time.Now().Add(time.Duration(expiresIn-300) * time.Second) // 5 minutes buffer

	return accessToken, nil
}

// generateJWT creates a JWT token for authentication
func (c *FCMClient) generateJWT() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   c.credentials.ClientEmail,
		"scope": "https://www.googleapis.com/auth/firebase.messaging",
		"aud":   c.credentials.TokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Parse private key
	block, _ := pem.Decode([]byte(c.credentials.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to parse private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not RSA")
	}

	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// exchangeJWTForAccessToken exchanges JWT for OAuth2 access token
func (c *FCMClient) exchangeJWTForAccessToken(ctx context.Context, jwt string) (string, int, error) {
	data := fmt.Sprintf("grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&assertion=%s", jwt)

	req, err := http.NewRequestWithContext(ctx, "POST", c.credentials.TokenURI, bytes.NewBufferString(data))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	return tokenResponse.AccessToken, tokenResponse.ExpiresIn, nil
}

// SendNotification sends a notification using the v1 API
func (c *FCMClient) SendNotification(ctx context.Context, message *FCMMessage) (*FCMResponse, error) {
	return c.send(ctx, message, false)
}

// ValidateNotification validates a notification without sending it
func (c *FCMClient) ValidateNotification(ctx context.Context, message *FCMMessage) (*FCMResponse, error) {
	return c.send(ctx, message, true)
}

// SendToToken sends a simple notification to a single device token
func (c *FCMClient) SendToToken(ctx context.Context, token, title, body string, data map[string]string) (*FCMResponse, error) {
	message := NewNotificationBuilder().
		ToToken(token).
		WithNotification(title, body).
		WithData(data).
		Build()

	return c.send(ctx, message, false)
}

// SendToTopic sends a notification to a topic
func (c *FCMClient) SendToTopic(ctx context.Context, topic, title, body string, data map[string]string, condition string) (*FCMResponse, error) {
	message := NewNotificationBuilder().
		ToTopic(topic).
		WithNotification(title, body).
		WithData(data).
		ToCondition(condition).
		Build()

	return c.send(ctx, message, false)
}

func (c *FCMClient) SendToTokens(ctx context.Context, tokens []string, title, body string, data map[string]string) ([]*FCMResponse, error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	var (
		responses []*FCMResponse
		mu        sync.Mutex
		wg        sync.WaitGroup
	)
	for _, token := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			message := NewNotificationBuilder().
				ToToken(token).
				WithNotification(title, body).
				WithData(data).
				Build()

			response, _ := c.send(ctx, message, false)
			if response != nil && response.Error != nil {
				mu.Lock()
				responses = append(responses, response)
				mu.Unlock()
			}
		}(token)
	}
	wg.Wait()
	return responses, nil
}

// SendDataOnly sends a data-only message
func (c *FCMClient) SendDataOnly(ctx context.Context, token string, data map[string]string) (*FCMResponse, error) {
	message := NewNotificationBuilder().
		ToToken(token).
		WithData(data).
		Build()

	return c.send(ctx, message, false)
}

// send performs the actual HTTP request to FCM v1 API
func (c *FCMClient) send(ctx context.Context, message *FCMMessage, validateOnly bool) (*FCMResponse, error) {
	// Validate message
	if err := c.validateMessage(message); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	// Get access token
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Create request payload
	request := &FCMRequest{
		ValidateOnly: validateOnly,
		Message:      message,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var fcmResponse FCMResponse
	if err := json.NewDecoder(resp.Body).Decode(&fcmResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fcmResponse.Token = message.Token // Include the token in the response for tracking

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		if fcmResponse.Error != nil {
			return &fcmResponse, fmt.Errorf("FCM error %d: %s", fcmResponse.Error.Code, fcmResponse.Error.Message)
		}
		return &fcmResponse, fmt.Errorf("FCM server returned status %d", resp.StatusCode)
	}

	return &fcmResponse, nil
}

// validateMessage validates the FCM message before sending
func (c *FCMClient) validateMessage(message *FCMMessage) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Must have exactly one target
	targets := 0
	if message.Token != "" {
		targets++
	}
	if message.Topic != "" {
		targets++
	}
	if message.Condition != "" {
		targets++
	}

	if targets != 1 {
		return fmt.Errorf("message must have exactly one target (token, topic, or condition)")
	}

	// Must have either notification or data payload
	if message.Notification == nil && message.Data == nil {
		return fmt.Errorf("message must have either notification or data payload")
	}

	return nil
}

// Helper functions for error checking
func IsRetryableError(err *FCMError) bool {
	if err == nil {
		return false
	}

	retryableCodes := []int{500, 502, 503, 504}
	for _, code := range retryableCodes {
		if err.Code == code {
			return true
		}
	}

	return err.Status == "UNAVAILABLE" || err.Status == "INTERNAL"
}

func IsInvalidTokenError(err *FCMError) bool {
	if err == nil {
		return false
	}

	return err.Status == "NOT_FOUND" ||
		err.Status == "INVALID_ARGUMENT" ||
		(err.Code == 400 && (err.Message == "The registration token is not a valid FCM registration token" ||
			err.Message == "Requested entity was not found."))
}
