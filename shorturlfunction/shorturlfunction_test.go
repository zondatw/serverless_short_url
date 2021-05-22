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

	shortHash := "ad6e1f62aa3fc5e4"
	originalUrl := "http://test.org/original_url"
	shortUrl := "http://short_url/" + shortHash

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

	checkBody := `{"url":"` + shortUrl + `"}`
	if strings.Trim(string(body), "\n") != checkBody {
		t.Errorf("Body got %v, want %v", body, checkBody)
	}

	// Check redis value
	if got, err := redisServer.Get(shortHash); err != nil {
		t.Errorf("Got redis value error %v", err)
	} else if got != originalUrl {
		t.Errorf("Got redis value %v, want %v", got, originalUrl)
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

func TestRedirect(t *testing.T) {
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	defer redisServer.Close()

	os.Setenv("REDISHOST", redisServer.Host())
	os.Setenv("REDISPORT", redisServer.Port())

	shortHash := "ad6e1f62aa3fc5e4"
	originalUrl := "http://test.org/original_url"
	shortUrl := "http://short_url/" + shortHash

	redisServer.Set(shortHash, originalUrl)

	req := httptest.NewRequest("GET", shortUrl, strings.NewReader(""))
	rr := httptest.NewRecorder()

	Redirect(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Register got status %v, want %v", rr.Code, http.StatusSeeOther)
	}

	returnUrl := rr.HeaderMap.Get("Location")
	if originalUrl != returnUrl {
		t.Errorf("Redirect got url %v, want %v", returnUrl, originalUrl)
	}

}

func TestRedirectNotExist(t *testing.T) {
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis.Run: %v", err)
	}
	defer redisServer.Close()

	os.Setenv("REDISHOST", redisServer.Host())
	os.Setenv("REDISPORT", redisServer.Port())

	shortHash := "ad6e1f62aa3fc5e4"
	shortUrl := "http://short_url/" + shortHash

	req := httptest.NewRequest("GET", shortUrl, strings.NewReader(""))
	rr := httptest.NewRecorder()

	Redirect(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Register got status %v, want %v", rr.Code, http.StatusNotFound)
	}

}
