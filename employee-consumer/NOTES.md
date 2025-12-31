## TODO

- customize logging [logging deps](https://blog.logrocket.com/five-structured-logging-packages-for-go/)
- add Swagger [Go Swagger](https://github.com/go-swagger/go-swagger)[Swagger doc](https://goswagger.io/use/spec/model.html)
- add k8 conf
  - add k8 conf for database
- Add logic in a new controller to add websocket option consume
  - [Option 1](https://blog.friendsofgo.tech/posts/introduccion-a-los-websockets-en-go/) (Gorilla)
  - [Option 2](https://yalantis.com/blog/how-to-build-websockets-in-go/) (STDLIB)
- investigate more [Middleware](https://blog.friendsofgo.tech/posts/middlewares-en-go/)
- investigate more [Request Context](https://pkg.go.dev/context)
- investigate gzip compression
- investigate performace [Profiling](https://blog.golang.org/pprof) and [Creating Web Applications with Go](https://app.pluralsight.com/course-player?clipId=f06c5b9c-22d0-4fae-a7ad-9975455e74ec)
- server push for http/2
- [Project Layout](https://github.com/golang-standards/project-layout)
- Lifecycle of a package (always execute init function)
- [JSON vs protocol buffers](https://www.bizety.com/2018/11/12/protocol-buffers-vs-json/)
  - [Performance json vs protobuf](https://auth0.com/blog/beating-json-performance-with-protobuf/)
- Go frameworks for distributed systems
  - [Go Micro](https://github.com/asim/go-micro)
  - [Go Kit](https://gokit.io/)
- Go frameworks for Web
  - [Gin](https://github.com/gin-gonic/gin)
  - [Echo](https://echo.labstack.com/guide/)
- Command line tools with Go
  - [Link](https://www.rapid7.com/blog/post/2016/08/04/build-a-simple-cli-tool-with-golang/)
  - [Link2](https://gobyexample.com/command-line-flags)
- Commands [Documentation](https://golang.org/doc/cmd)
- [Http2 Fundamentals](https://developers.google.com/web/fundamentals/performance/http2?hl=es)
- service-to-service communication
  - [twirp](https://github.com/twitchtv/twirp)
  - Grpc
  - Rest
  - Proto buffers
- [Generics in 1.8](https://go.dev/doc/tutorial/generics)

### Date and Time:

```golang
s1, _:= time.Parse(time.RFC3339, “2018-12-12T11:45:26.371Z”)
// https://flaviocopes.com/go-date-time-format/
```

### Static Server

```golang
http.ListenAndServe(":8080", http.FileServer(http.Dir("public"))) // server public dir
```

### Generate ssl certificate

[More Details](https://gist.github.com/denji/12b3a568f092ab951456)

```bash
# get help 
go run /usr/local/opt/go/libexec/src/crypto/tls/generate_cert.go -h 
go run $GOROOT/src/crypto/tls/generate_cert.go -h
# generate certificate for localhost
go run $GOROOT/src/crypto/tls/generate_cert.go -host localhost
```

### See all versions for a dep

`go list -m -versions github.com/google/uuid`

### Userful commands

```bash
go mod why
go mod vendor
go mod graph
go mod edit
go mod download
```

### Add OAUTH2

https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/
https://github.com/go-oauth2/oauth2
https://auth0.com/blog/authentication-in-golang/

https://www.sohamkamani.com/golang/jwt-authentication/

### Unit testing notes

https://medium.com/@elliotchance/a-new-simpler-way-to-do-dependency-injection-in-go-9e191bef50d5
https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
https://segment.com/blog/5-advanced-testing-techniques-in-go/
