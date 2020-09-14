# avoxi-challenge
[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=github.com/jamalyusuf/avoxi-challenge.git)

Code Challenge for avoxi.com interview. 

### Scenario (note, this is fictional, but gives a sense of the types of requests we might encounter):

Our product team has heard from several customers that we restrict users to logging in to their UI accounts from selected countries to prevent them from outsourcing their work to others.  For an initial phase one we’re not going to worry about VPN connectivity, only the presented IP address.

The team has designed a solution where the customer database will hold the white listed countries and the API gateway will capture the requesting IP address, check the target customer for restrictions, and send the data elements to a new service you are going to create.  

The new service will be an HTTP-based API that receives an IP address and a white list of countries.  The API should return an indicator if the IP address is within the blacklisted listed countries.  You can get a data set of IP address to country mappings from https://dev.maxmind.com/geoip/geoip2/geolite2/.

We do our backend development in Go (Golang) and prefer solutions in that language.

We'll be explicitly looking at coding style, code organization, API design, and operational/maintenance aspects such as logging and error handling.  We'll also be giving bonus points for things like
* Documenting a plan for keeping the mapping data up to date.  Extra bonus points for implementing the solution.
* Including a Docker file for the running service
* Including a Kubernetes YAML file for running the service in an existing cluster
* Exposing the service as gRPC in addition to HTTP
* Other extensions to the service you think would be worthwhile.  If you do so, please include a brief description of the feature and justification for its inclusion.  Think of this as what you would have said during the design meeting to convince your team the effort was necessary.


## Running
Running `main.go` starts a web server on https://0.0.0.0:11000/. You can configure
the port used with the `$PORT` environment variable, and to serve on HTTP set
`$SERVE_HTTP=true`.

```
$ go run ./cmd/api/
```

An OpenAPI UI is served on https://0.0.0.0:11000/.

## grpcurl payloads
grpcurl is a command-line tool that lets you interact with gRPC servers. It's basically curl for gRPC servers.
## install via brew
```
$ brew install grpcurl
```

## from source 
You can use the `go` tool to install `grpcurl`:
```
go get github.com/fullstorydev/grpcurl/...
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

## **Method GeoIPCheck**
### payload
```
$ grpcurl -insecure -d '{"IP": "95.31.18.119", "AllowedCountries": [ "RU", "US"] }' localhost:10000 v1.IPFilterService.GeoIPCheck
```

### request 
```
$ grpcurl -insecure localhost:10000 describe v1.GeoIPCheckRequest
```

## **Method IPLocation**

### payload
```
$ grpcurl -insecure -d '{"IP": "95.31.18.119"}' localhost:10000 v1.IPFilterService.IPLocation
```

### request 
```
$ grpcurl -insecure localhost:10000 describe v1.IPLocationRequest
```

### **Method Health**

### payload
```
$ grpcurl -insecure -d '{}' localhost:10000 v1.IPFilterService.Health
```

### request 
```
$ grpcurl -insecure localhost:10000 describe v1.HealthRequest
```
## Requirements

Generating the files requires the `protoc` protobuf compiler.
Please install it according to the
[installation instructions](https://github.com/google/protobuf#protocol-compiler-installation)
for your specific platform.

