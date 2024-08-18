package omnicron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// HTTPClient interface defines the Do method for HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const APIVersion = "api/v1"

// Client represents the client for the Omnicron API.
type Client struct {
	baseurl    string
	apikey     string
	debug      bool
	httpClient HTTPClient
}

// ClientOption is a type for setting options in the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseurl = baseURL
	}
}

// WithDebug enables or disables debug mode.
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// NewClient creates a new Client with the given API key and options.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apikey:  apiKey,
		debug:   false,
		baseurl: "https://omnicron-latest.onrender.com/", // default base URL for Omnicron
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}

// newJSONPostRequest sends a POST request with a JSON payload.
func (c *Client) newJSONPostRequest(ctx context.Context, path, model string, payload interface{}) ([]byte, error) {
	fullURLPath := c.baseurl + APIVersion + path
	if model != "" {
		fullURLPath = c.withModelQueryParameters(fullURLPath, model)
	}
	if c.debug {
		log.Printf("full URL path: %s", fullURLPath)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if c.debug {
		log.Println(string(body))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURLPath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apikey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apikey))
	}

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if c.debug {
		log.Println(string(resBody))
	}
	if res.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.Unmarshal(resBody, &errResp); err != nil {
			return nil, fmt.Errorf("error unmarshalling: %s", resBody)
		}
		return nil, fmt.Errorf("API request failed: %s", errResp.Error)
	}
	return resBody, nil
}

// newFormWithFilePostRequest sends a POST request with a multipart form-data payload.
func (c *Client) newFormWithFilePostRequest(ctx context.Context, path, model string, payload interface{}) ([]byte, error) {
	fullURLPath := c.baseurl + APIVersion + path
	if model != "" {
		fullURLPath = c.withModelQueryParameters(fullURLPath, model)
	}
	if c.debug {
		log.Printf("full URL path: %s", fullURLPath)
	}

	// Create a buffer to hold the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	v := reflect.ValueOf(payload)
	typeOfParams := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := typeOfParams.Field(i).Tag.Get("form")

		if fieldName == "" {
			continue
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		switch field.Interface().(type) {
		case *os.File:
			if c.debug {
				log.Printf("file is present, fieldname is %s", fieldName)
			}
			if err := addFileField(writer, fieldName, field.Interface().(*os.File)); err != nil {
				log.Printf("Error adding file field: %v", err)
				return nil, err
			}
		case *int:
			if c.debug {
				log.Printf("int is present, fieldname is %s", fieldName)
			}
			if err := addField(writer, fieldName, strconv.Itoa(*field.Interface().(*int))); err != nil {
				log.Printf("Error adding int field: %v", err)
				return nil, err
			}
		case *float64:
			if c.debug {
				log.Printf("float64 is present, fieldname is %s", fieldName)
			}
			if err := addField(writer, fieldName, fmt.Sprintf("%f", *field.Interface().(*float64))); err != nil {
				log.Printf("Error adding float64 field: %v", err)
				return nil, err
			}
		case *string:
			if c.debug {
				log.Printf("string is present, fieldname is %s", fieldName)
			}
			if err := addField(writer, fieldName, *field.Interface().(*string)); err != nil {
				log.Printf("Error adding string field: %v", err)
				return nil, err
			}
		case *bool:
			if c.debug {
				log.Printf("bool is present, fieldname is %s", fieldName)
			}
			if err := addField(writer, fieldName, strconv.FormatBool(*field.Interface().(*bool))); err != nil {
				log.Printf("Error adding bool field: %v", err)
				return nil, err
			}
		default:
			if c.debug {
				log.Printf("Default type, fieldname is required:%s", fieldName)
			}
			if err := addField(writer, fieldName, fmt.Sprintf("%v", field.Interface())); err != nil {
				log.Printf("Error adding default field: %v", err)
				return nil, err
			}
		}
	}

	writer.Close()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURLPath, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	if c.apikey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apikey))
	}

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		log.Println(string(resBody))
	}

	if res.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.Unmarshal(resBody, &errResp); err != nil {
			return nil, fmt.Errorf("error unmarshalling: %s", resBody)
		}
		return nil, fmt.Errorf("API request failed: %s", errResp.Error)
	}

	return resBody, nil
}

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// addField adds a form field with the given key and value.
func addField(writer *multipart.Writer, key string, value string) error {
	key = formatFieldName(key)
	err := writer.WriteField(key, value)
	if err != nil {
		return err
	}
	return nil
}

// addFileField adds a file field to the multipart writer.
func addFileField(writer *multipart.Writer, fieldname string, file *os.File) error {
	defer file.Close()
	fieldname = formatFieldName(fieldname)
	if file == nil {
		return ErrNoFileProvided
	}

	fw, err := writer.CreateFormFile(fieldname, file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err
	}
	return nil
}

// formats the field name, removes the ",omitempty"
func formatFieldName(key string) string {
	keySlice := strings.Split(key, ",")
	return keySlice[0]
}

// withModelQueryParameters appends the model query parameter to the URL.
func (c *Client) withModelQueryParameters(fullURLPath, model string) string {
	params := url.Values{}
	params.Add("model", model)
	return fullURLPath + "?" + params.Encode()
}
