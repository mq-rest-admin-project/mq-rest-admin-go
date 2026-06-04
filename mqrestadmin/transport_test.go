package mqrestadmin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPTransport_PostJSON_Success(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		if r.Header.Get("X-Custom") != "header-value" {
			t.Errorf("X-Custom header = %q, want header-value", r.Header.Get("X-Custom"))
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		w.Header().Set("X-Response", "resp-value")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer server.Close()

	certPool := x509.NewCertPool()
	certPool.AddCert(server.Certificate())
	transport := &HTTPTransport{
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: certPool},
	}

	response, err := transport.PostJSON(
		context.Background(),
		server.URL+"/test",
		map[string]any{"key": "value"},
		map[string]string{"X-Custom": "header-value"},
		30*time.Second,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", response.StatusCode)
	}
	if response.Headers["X-Response"] != "resp-value" {
		t.Errorf("X-Response header = %q, want resp-value", response.Headers["X-Response"])
	}
	if response.Body == "" {
		t.Error("expected non-empty response body")
	}
}

func TestHTTPTransport_PostJSON_NetworkError(t *testing.T) {
	transport := &HTTPTransport{}

	_, err := transport.PostJSON(
		context.Background(),
		"https://127.0.0.1:1/nonexistent",
		map[string]any{"key": "value"},
		nil,
		1*time.Second,
	)
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T: %v", err, err)
	}
}

func TestBuildClient_NilTLSConfig(t *testing.T) {
	transport := &HTTPTransport{TLSConfig: nil}
	client := transport.buildClient(10 * time.Second)
	if client.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", client.Timeout)
	}
}

func TestBuildClient_WithTLSConfig(t *testing.T) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "test-server",
	}
	transport := &HTTPTransport{TLSConfig: tlsConfig}
	client := transport.buildClient(10 * time.Second)

	httpTransport := client.Transport.(*http.Transport)
	if httpTransport.TLSClientConfig.ServerName != "test-server" {
		t.Error("expected TLSConfig to be cloned with ServerName")
	}
	if httpTransport.TLSClientConfig.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify = false (TLS is always verified)")
	}
}

func TestHTTPTransport_PostJSON_InvalidURL(t *testing.T) {
	transport := &HTTPTransport{}

	_, err := transport.PostJSON(
		context.Background(),
		"://invalid-url",
		map[string]any{"key": "value"},
		nil,
		1*time.Second,
	)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}

	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T: %v", err, err)
	}
}

func TestNewSession_WithTLSCAFile(t *testing.T) {
	certPEM, _ := generateSelfSignedCert(t)
	caFile := writeTempFile(t, "ca.pem", certPEM)

	session, err := NewSession("https://localhost:9443/x", "QM1",
		BasicAuth{Username: "u", Password: "p"}, WithTLSCAFile(caFile))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	httpTransport, ok := session.transport.(*HTTPTransport)
	if !ok {
		t.Fatalf("expected *HTTPTransport, got %T", session.transport)
	}
	if httpTransport.TLSConfig == nil || httpTransport.TLSConfig.RootCAs == nil {
		t.Fatal("expected TLSConfig with RootCAs set")
	}
}

func TestNewSession_WithTLSCAFile_MissingFile(t *testing.T) {
	_, err := NewSession("https://localhost:9443/x", "QM1",
		BasicAuth{Username: "u", Password: "p"}, WithTLSCAFile("/nonexistent/ca.pem"))
	if err == nil {
		t.Fatal("expected error for missing CA file")
	}
}

func TestNewSession_WithTLSCAFile_InvalidPEM(t *testing.T) {
	caFile := writeTempFile(t, "bad.pem", []byte("not a certificate"))
	_, err := NewSession("https://localhost:9443/x", "QM1",
		BasicAuth{Username: "u", Password: "p"}, WithTLSCAFile(caFile))
	if err == nil {
		t.Fatal("expected error for invalid PEM CA file")
	}
}

func TestNewSession_CertificateAuthWithTLSCAFile(t *testing.T) {
	certPEM, keyPEM := generateSelfSignedCert(t)
	certFile := writeTempFile(t, "cert.pem", certPEM)
	keyFile := writeTempFile(t, "key.pem", keyPEM)
	caFile := writeTempFile(t, "ca.pem", certPEM)

	session, err := NewSession("https://localhost:9443/x", "QM1",
		CertificateAuth{CertPath: certFile, KeyPath: keyFile}, WithTLSCAFile(caFile))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	httpTransport, ok := session.transport.(*HTTPTransport)
	if !ok {
		t.Fatalf("expected *HTTPTransport, got %T", session.transport)
	}
	if len(httpTransport.TLSConfig.Certificates) == 0 {
		t.Error("expected client certificate to be set")
	}
	if httpTransport.TLSConfig.RootCAs == nil {
		t.Error("expected RootCAs to be set")
	}
}

