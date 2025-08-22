package testing

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewServer(t *testing.T) {
	g := NewWithT(t)
	
	srv := NewServer()
	defer srv.Teardown()

	g.Expect(srv).NotTo(BeNil(), "Server should not be nil")
	g.Expect(srv.Mux).NotTo(BeNil(), "Mux should not be nil")
	g.Expect(srv.Server).NotTo(BeNil(), "Server.Server should not be nil")
}

func TestServer_Teardown(t *testing.T) {
	g := NewWithT(t)
	
	srv := NewServer()
	
	// Should not panic when tearing down
	g.Expect(func() { srv.Teardown() }).NotTo(Panic())
}

func TestServer_URL(t *testing.T) {
	g := NewWithT(t)
	
	srv := NewServer()
	defer srv.Teardown()

	g.Expect(srv.URL).NotTo(BeEmpty())
	g.Expect(srv.URL).To(HavePrefix("http://"))
}

func TestServer_Integration(t *testing.T) {
	g := NewWithT(t)
	
	srv := NewServer()
	defer srv.Teardown()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv.Mux.HandleFunc("/test", testHandler)

	// Create a test client and make a request to our test server
	client := &http.Client{}
	req, err := http.NewRequest("GET", srv.URL+"/test", nil)
	g.Expect(err).NotTo(HaveOccurred())

	resp, err := client.Do(req)
	g.Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
