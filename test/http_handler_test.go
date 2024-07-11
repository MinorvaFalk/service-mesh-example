package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"worker-mesh/internal/handler"
	"worker-mesh/pkg/router"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProducer struct{ mock.Mock }

func newMockProducer() *mockProducer {
	return &mockProducer{}
}

func (m *mockProducer) Publish(topic string, body []byte) error {
	return nil
}

func (m *mockProducer) DeferredPublish(topic string, delay time.Duration, body []byte) error {
	return nil
}

func (m *mockProducer) Ping() error {
	return nil
}

func TestHTTPHandler(t *testing.T) {
	body := &handler.Notification{
		Title: "Hello World",
		Body:  "Hello World",
	}

	tests := []struct {
		description  string
		route        string
		body         []byte
		expectedCode int
	}{
		{
			description:  "request /publish without body",
			route:        "/publish/",
			body:         nil,
			expectedCode: 400,
		},
		{
			description: "request /publish with body",
			route:       "/publish",
			body:        body.Byte(),
			// body: (*handler.Notification).Byte(body),
			expectedCode: 200,
		},
		{
			description:  "request /publish/defer without body",
			route:        "/publish/defer/1000",
			body:         nil,
			expectedCode: 400,
		},
		{
			description:  "request /publish/defer without delay param",
			route:        "/publish/defer/",
			body:         body.Byte(),
			expectedCode: 404,
		},
		{
			description:  "request /publish/defer with body and delay param",
			route:        "/publish/defer/1000",
			body:         body.Byte(),
			expectedCode: 200,
		},
	}

	producer := newMockProducer()
	producer.On("Publish", "notification", nil).Return(nil)
	producer.On("DeferredPublish", "notification", "1000", nil).Return(nil)

	httpHandler := handler.NewHttpHandler(producer)
	app := router.NewFiber(httpHandler)

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")
			defer req.Body.Close()

			res, _ := app.Test(req, 1000)

			assert.Equal(t, test.expectedCode, res.StatusCode)
		})
	}
}
