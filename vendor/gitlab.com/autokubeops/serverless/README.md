# Serverless

This project provides Serverless runtimes for use with Knative.

## Runtimes

* Go


### Go

Usage:

```go
package main

import (
	"gitlab.com/autokubeops/serverless"
	"net/http"
)

func main() {
	serverless.Run(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// put your application logic here
		_, _ = w.Write([]byte("OK"))
	}))
}
```