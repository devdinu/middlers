# Middlers
Go HTTP Middlewares

## Request Logger Middleware

You could wrap it with a handler `http.HandlerFunc` or `http.Handler`
``` 
    gomw.RequestLogger(handler)
```

You could use a custom logger or log.New(io.Writer...) any interface which have Println implementation
```
    gomw.RequestLogger(handler, customLogger)
```
