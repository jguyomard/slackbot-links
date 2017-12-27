package api

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jguyomard/slackbot-links/src/config"
	"github.com/jguyomard/slackbot-links/src/links"
)

var (
	lastReply string
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	lastReply = string(p)
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	lastReply = s
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func init() {
	configFilePtr := flag.String("config-file", "/etc/slackbot-links/config.yaml", "conf file path")
	flag.Parse()
	config.SetFilePath(*configFilePtr)

	// Disable Mercury
	config.Get().MercuryAPIKey = ""

	// Add some links to test
	links.Init()
	links.Restore("testdata/links.json")
}

func TestListenPort(t *testing.T) {
	if getListenPort() != config.Get().APIListenPort {
		t.Fatal("TestListenPort: ListenPort doesn't match")
	}
}

func TestListen_listLinks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/links", nil)
	testResponseEqualsFixture(t, req, "listLinks")
}

func TestListen_searchLinks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/links?search=hugo", nil)
	testResponseEqualsFixture(t, req, "searchLinks")
}

func TestListen_searchLinksWithoutResult(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/links?search=toto", nil)
	testResponseEqualsFixture(t, req, "searchLinksWithoutResult")
}

func TestListen_getLink(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/links/test", nil)
	testResponseEqualsFixture(t, req, "getLink")
}

func TestListen_getLinkNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/links/not-found", nil)
	testResponseEqualsFixture(t, req, "getLinkNotFound")
}

func TestListen_deleteLink(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/v1/links/test", nil)
	testResponseEqualsFixture(t, req, "deleteLink")

	req, _ = http.NewRequest("DELETE", "/v1/links/test2", nil)
	testResponseEqualsFixture(t, req, "deleteLink")
}

func TestListen_deleteLinkNotFound(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/v1/links/not-found", nil)
	testResponseEqualsFixture(t, req, "deleteLinkNotFound")
}

func TestListen_notFound(t *testing.T) {
	req, _ := http.NewRequest("Get", "/not/found", nil)
	testResponseEqualsFixture(t, req, "notFound")
}

func testResponseEqualsFixture(t *testing.T, req *http.Request, fixtureName string) {
	lastReply = ""

	w := new(mockResponseWriter)
	router := NewAPIMiddleware(createRoutes())
	router.ServeHTTP(w, req)

	fixture := getFixture(fixtureName)
	if strings.TrimSpace(lastReply) != strings.TrimSpace(fixture) {
		t.Fatalf("TestListen: reply doesn't match:\n- have: %s\n- want: %s\n", lastReply, fixture)
	}
}

func getFixture(name string) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.txt", name))
	if err != nil {
		panic(err)
	}
	return string(data)
}
