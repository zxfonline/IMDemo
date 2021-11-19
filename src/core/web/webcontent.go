package web

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/zxfonline/IMDemo/core/golangtrace"
)

var contextType reflect.Type

var EnableTracing = true

func init() {
	contextType = reflect.TypeOf(Context{})
}

type Context struct {
	Request     *http.Request
	RequestBody []byte
	Params      map[string]string
	Server      *Server
	http.ResponseWriter
	tr golangtrace.Trace
}

var (
	acceptsHtmlRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXmlRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJsonRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
)

func (ctx *Context) TraceFinish() {
	if ctx.tr != nil {
		ctx.tr.Finish()
		ctx.tr = nil
	}
}
func (ctx *Context) TracePrintf(format string, a ...interface{}) {
	if ctx.tr != nil {
		ctx.tr.LazyPrintf(format, a...)
	}
}

func (ctx *Context) TraceErrorf(format string, a ...interface{}) {
	if ctx.tr != nil {
		ctx.tr.LazyPrintf(format, a...)
		ctx.tr.SetError()
	}
}
func (ctx *Context) Protocol() string {
	return ctx.Request.Proto
}

func (ctx *Context) Uri() string {
	return ctx.Request.RequestURI
}

func (ctx *Context) Url() string {
	return ctx.Request.URL.Path
}

func (ctx *Context) Site() string {
	return ctx.Scheme() + "://" + ctx.Domain()
}

func (ctx *Context) Scheme() string {
	if ctx.Request.URL.Scheme != "" {
		return ctx.Request.URL.Scheme
	}
	if ctx.Request.TLS == nil {
		return "http"
	}
	return "https"
}

func (ctx *Context) Domain() string {
	return ctx.Host()
}

func (ctx *Context) Host() string {
	if ctx.Request.Host != "" {
		hostParts := strings.Split(ctx.Request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return ctx.Request.Host
	}
	return "localhost"
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

//GET,POST,HEAD,OPTIONS,PUT,DELETE,PATCH...
func (ctx *Context) Is(method string) bool {
	return ctx.Method() == method
}

func (ctx *Context) IsAjax() bool {
	return ctx.Header("X-Requested-With") == "XMLHttpRequest"
}

func (ctx *Context) IsSecure() bool {
	return ctx.Scheme() == "https"
}

func (ctx *Context) IsWebsocket() bool {
	return ctx.Header("Upgrade") == "websocket"
}

func (ctx *Context) AcceptsHtml() bool {
	return acceptsHtmlRegex.MatchString(ctx.Header("Accept"))
}

func (ctx *Context) AcceptsXml() bool {
	return acceptsXmlRegex.MatchString(ctx.Header("Accept"))
}

func (ctx *Context) AcceptsJson() bool {
	return acceptsJsonRegex.MatchString(ctx.Header("Accept"))
}

// IP returns request client ip.
// if in proxy, return first proxy id.
// if error, return 127.0.0.1.
func (ctx *Context) IP() string {
	ips := ctx.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	ip := strings.Split(ctx.Request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (ctx *Context) Proxy() []string {
	if ips := ctx.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (ctx *Context) UserAgent() string {
	return ctx.Header("User-Agent")
}

func (ctx *Context) Header(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *Context) CopyBody() []byte {
	requestBody, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body.Close()
	bf := bytes.NewBuffer(requestBody)
	ctx.Request.Body = ioutil.NopCloser(bf)
	ctx.Request.ContentLength = int64(len(requestBody))
	ctx.RequestBody = requestBody
	return requestBody
}

func (ctx *Context) IsUpload() bool {
	return strings.Contains(ctx.Header("Content-Type"), "multipart/form-data")
}

func (ctx *Context) ParseFormOrMutliForm(maxMemory int64) error {
	if ctx.IsUpload() {
		if err := ctx.Request.ParseMultipartForm(maxMemory); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := ctx.Request.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}

func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

func (ctx *Context) WriteBytes(content []byte) {
	ctx.ResponseWriter.Write(content)
}

// 4xx or 5xx
func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

// 4xx or 5xx
func (ctx *Context) AbortBytes(status int, body []byte) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write(body)
}

// 3xx
func (ctx *Context) Redirect(status int, url_ string) {
	ctx.ResponseWriter.Header().Set("Location", url_)
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte("Redirecting to: " + url_))
}

// 304
func (ctx *Context) NotModified() {
	ctx.ResponseWriter.WriteHeader(304)
}

// 404
func (ctx *Context) NotFound(message string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(message))
}

//401
func (ctx *Context) Unauthorized() {
	ctx.ResponseWriter.WriteHeader(401)
}

//403
func (ctx *Context) Forbidden() {
	ctx.ResponseWriter.WriteHeader(403)
}

func (ctx *Context) ContentType(val string) string {
	var ctype string
	if strings.ContainsRune(val, '/') {
		ctype = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}
		ctype = mime.TypeByExtension(val)
	}
	if ctype != "" {
		ctx.SetHeader("Content-Type", ctype, true)
	}
	return ctype
}

func (ctx *Context) SetHeader(hdr string, val string, overwritten bool) {
	if overwritten {
		ctx.ResponseWriter.Header().Set(hdr, val)
	} else {
		ctx.ResponseWriter.Header().Add(hdr, val)
	}
}

func (ctx *Context) SetExpires(expires time.Duration) {
	ctx.SetHeader("Expires", webTime(time.Now().Add(expires)), true)
}

func (ctx *Context) SetCacheControl(expires time.Duration) {
	if expires > time.Second {
		ctx.SetHeader("Expires", webTime(time.Now().Add(expires)), true)
		ctx.SetHeader("Cache-Control", fmt.Sprintf("max-age=%d", int(expires.Seconds())), true)
	} else {
		ctx.SetHeader("Expires", webTime(time.Now()), true)
		ctx.SetHeader("Cache-Control", "no-cache", true)
	}
}

func (ctx *Context) SetLastModified(modTime time.Time) {
	ctx.SetHeader("Last-Modified", webTime(modTime), true)
}

func webTime(t time.Time) string {
	ftime := t.UTC().Format(time.RFC1123)
	if strings.HasSuffix(ftime, "UTC") {
		ftime = ftime[0:len(ftime)-3] + "GMT"
	}
	return ftime
}
