# Middlers
Go HTTP Middlewares

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## TODO:
- Support for negroni middleware
