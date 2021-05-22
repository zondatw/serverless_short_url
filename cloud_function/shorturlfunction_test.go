package shorturlfunction

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestConvertToShort(t *testing.T) {
	shortHash := convertToShort("http://test.org/original_url")
	if "ad6e1f62aa3fc5e4" != shortHash {
		t.Errorf("Get wrong shortHash %v", shortHash)
	}
}

func TestRegister(t *testing.T) {
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	defer redisServer.Close()

	os.Setenv("REDISHOST", redisServer.Host())
	os.Setenv("REDISPORT", redisServer.Port())
	os.Setenv("SHORTURLBASE", "http://short_url/")

	originalUrl := "http://test.org/original_url"
	shortUrl := "http://short_url/ad6e1f62aa3fc5e4"

	rg := registerRequestStruct{Url: originalUrl}
	jsonByte, _ := json.Marshal(rg)
	req := httptest.NewRequest("POST", "/", strings.NewReader(string(jsonByte)))
	rr := httptest.NewRecorder()

	Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Register got status %v, want %v", rr.Code, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("Read recorder result body error %v", err)
	}

	checkBody := `{"url":"http://short_url/ad6e1f62aa3fc5e4"}`
	if strings.Trim(string(body), "\n") != checkBody {
		t.Errorf("Body got %v, want %v", body, checkBody)
	}

	// Check redis value
	if got, err := redisServer.Get(originalUrl); err != nil {
		t.Errorf("Got redis value error %v", err)
	} else if got != shortUrl {
		t.Errorf("Got redis value %v, want %v", got, shortUrl)
	}
}

func TestRegisterInValidData(t *testing.T) {
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	defer redisServer.Close()

	os.Setenv("REDISHOST", redisServer.Host())
	os.Setenv("REDISPORT", redisServer.Port())
	os.Setenv("SHORTURLBASE", "http://short_url/")

	req := httptest.NewRequest("POST", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()

	Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Register got status %v, want %v", rr.Code, http.StatusBadRequest)
	}
}
