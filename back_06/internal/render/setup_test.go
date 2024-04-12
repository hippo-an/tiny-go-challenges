package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())

}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (mw *myWriter) WriteHeader(statusCode int) {

}
