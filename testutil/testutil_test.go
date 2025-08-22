package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
)

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNewTestServer(t *testing.T) {
	g := NewWithT(t)
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv, cleanup := NewTestServer(t, handler)
	defer cleanup()

	resp, err := http.Get(srv.URL)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}

func TestWriteJSON(t *testing.T) {
	g := NewWithT(t)
	
	testData := testStruct{Name: "Test", Age: 30}

	rr := httptest.NewRecorder()
	err := WriteJSON(rr, http.StatusOK, testData)

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rr.Code).To(Equal(http.StatusOK))
	g.Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))

	var result testStruct
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).To(Equal(testData))
}

func TestMustJSON(t *testing.T) {
	g := NewWithT(t)
	
	testData := testStruct{Name: "Test", Age: 30}

	// Test successful marshaling
	data := MustJSON(t, testData)

	var result testStruct
	err := json.Unmarshal(data, &result)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).To(Equal(testData))
}
