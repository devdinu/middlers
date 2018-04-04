# Middlers
[![Build Status](https://travis-ci.org/devdinu/middlers.svg?branch=master)](https://travis-ci.org/devdinu/middlers)

Go HTTP Middlewares

## Installation
Add middlers package as your dependency, with go get or either of your dependency management tool.
- Dep `dep ensure add -v github.com/devdinu/middlers`
- Glide `glide install -v github.com/devdinu/middlers`

### Request Logger Middleware

You could wrap it with a handler `http.HandlerFunc` or `http.Handler`
``` 
    h := gomw.RequestLogger(handler)
```

You could use a custom logger or log.New(io.Writer...) any interface which have Println implementation
```
    h := gomw.RequestLogger(handler, customLogger)
```

### Timeout Middleware

This ensures the request process completes without timeout, else writes `GatewayTimeout` header
```
    duration := 100*time.Millisecond
    withTimeout := gomw.Timeout(handler, duration)
```
This changes the `http.Request` context to `context.WithTimeout(r.Context(), duration)`

### Recovery Handler
You could use this handler to recover from any panic from your handlers, and return `500` Internal Server Error.
```
    logger // adheres to Println(...interface{}), also could be nil
    withRecovery := gomw.Recovery(logger)(next)
```

### Filter Middleware
Filter middleware could be used to block requests based on some `predicate`, you could use this to validate request based on header, url or body
 ```
    predicate := func(r *http.Request) bool {
        var result bool
        // bool to decide whether to block / pass the request
        return result
    }
    withFilter := gomw.Filter(predicate, handler)
```

### Stats Middleware
Stats middleware reports the url calls to stats, it uses `Increment(string)` interface to increment the url along with status code. `host:/some/url, with status 200`
would increment stats as `some_url_ok` `http.StatusText` used to convert the statuscode, and `/` is replaced with `_`
```
    withReporter := gomw.StatsReporter(reporter)(next)
```

### BeforeAfter Middleware
You could run a custom function before executes before the handler, and after executes after handler completion. After is executed even if handler panics

```
    before := func() { ... }
    after := func() { ... }
    mw := gomw.ExecutionHooks(before, after)(next)
```


## Contribution
- create issues or share your opinions/future enhancements
- clone repo and make the changes you want, new features, run
```
    dep ensure -v # manage dependencies
    go test -v
    golint . | grep -iEv 'exported.*should have comment'
    go vet .
```
fix lint and vet errors if any, and create a PR.

reach out to me dineshkumar in gophers slack.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## TODO:
- Support for negroni middleware
- Make all middlewares adhere to `Middleware`
- stats middleware report statuscode as tags
