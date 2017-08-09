# JSONP Go http middleware

JSONP is a common technique used to communicate with a JSON-serving Web Service with a
Web browser over cross-domains, in place of a XHR request. There is a lot written about
JSONP out there, but the tl;dr on it is a Javascript http client requesting JSONP
will write a `<script>` tag to the head of a page, with the `src` to an API endpoint,
with the addition of a `callback` (or `jsonp`) query parameter that represents a
randomly-named listener function that will parse the request when it comes back from
the server.

This middleware will work with anything that supports standard `http.Handler`. The code
is small, so go read it, but it just buffers the response from the rest of the chain,
and if its a JSON request with a callback, then it will wrap the response in the callback
function before writing it to the actual response writer.

Any feedback is welcome and appreciated!

## Example

```go
// JSONP example using Chi http router.. but anything that accepts
// a http.Handler will work
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jsonp"
	"github.com/go-chi/render"
)

func main() {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(jsonp.Handler)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := &SomeObj{"superman"}
		render.JSON(w, r, data)
	})

	err := http.ListenAndServe(":4444", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type SomeObj struct {
	Name string `json:"name"`
}
```

*Output:*

```
$ curl -v "http://localhost:4444/"
*   Trying ::1...
* Connected to localhost (::1) port 4444 (#0)
> GET / HTTP/1.1
> Host: localhost:4444
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Fri, 14 Aug 2015 19:11:44 GMT
< Content-Length: 19
<
* Connection #0 to host localhost left intact
{"name":"superman"}

$ curl -v "http://localhost:4444/?callback=X"
*   Trying ::1...
* Connected to localhost (::1) port 4444 (#0)
> GET /?callback=X HTTP/1.1
> Host: localhost:4444
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Length: 122
< Content-Type: application/javascript
< Date: Fri, 14 Aug 2015 19:11:49 GMT
<
* Connection #0 to host localhost left intact
X({"meta":{"content-length":19,"content-type":"application/json; charset=UTF-8","status":200},"data":{"name":"superman"}})
```

## NOTES

Since JSONP must always respond with a 200, as thats what the browser `<script>`
tag expects, a nice pattern that is also used in the GitHub API is to put the HTTP
response headers in a `"meta"` hash, and the HTTP response body in `"data"`. Like so..

```json
JsonpCallbackFn_abc123etc({
  "meta": {
    "Status": 200,
    "Content-Type": "application/json",
    "Content-Length": "19",
    "etc": "etc"
  },
  "data": { "name": "yummy" }
})
```

## LICENSE

BSD
