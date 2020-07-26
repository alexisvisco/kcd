<p align="center">
	<img width="460" height="300" src="./.github/golang-ss.gif">
</p>
<p align="center">
	<img width="93" height="20" src="https://github.com/expectedsh/kcd/workflows/Go/badge.svg">
	<img width="78" height="20" src="https://goreportcard.com/badge/github.com/expectedsh/kcd">
	<img src="https://codecov.io/gh/expectedsh/kcd/branch/master/graph/badge.svg" />
</p>

------

## :stars2: KCD 

KCD is a grandiose REST helper that wrap your shiny handler into a classic http handler. It manage all you want for building REST services.

This library is **opinionated** by default but **fully customizable** which mean it uses some other libraries like Chi for instance. KCD is modular so each pieces of the code that rely on a specific library can be changed. 

## What KCD does exactly 

Okay so KCD will wrap your cool handler into a http handler. The magic happen with this function:

`kcd.Handler(YourShinyHandler,  http.StatusOK)` (which returns a http.HandlerFunc)

Your handler is the `YourShinyHandler` parameter, it accepts: 
```go
func([response http.ResponseWriter], [request *http.Request], [input object ptr]) ([output object], error)
```
 
The only parameter in your shiny handler that is required is the returned error. 

**If there are any errors at some point KCD will call the [error hook](hooks.go#L69) to provide a REST generic error**.

1. If there is a custom input parameter (a pointer to a structure) it will:
    1. Run all [extractors](extractors.go) to extract values from the request into the input (query parameters, path, header, default value ...)
    2. Run the JSON body [bind hook](hooks.go#L144)
    3. Validate the input through the [vaidate hook](hooks.go#L18)
3. If all is good it will then call your shiny handler with all required arguments
4. Then if there is an output parameter it will call the [render hook](hooks.go#L39)

That's all. Well that's it, that's all you should have done if you didn't have KCD. 

## :coffee: Benefits

- More readable code
- Focus on what it matters: business code
- No more code duplication with unmarshalling, verifying, validating, marshalling ...
- You could have one interface for the client and server implementation

## :muscle: Example

- [*examples/simple/main.go*](./examples/simple/main.go)

