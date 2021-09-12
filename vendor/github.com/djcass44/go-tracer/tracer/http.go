/*
 *    Copyright 2020 Django Cass
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package tracer

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ContextKey int

const (
	// Deprecated: replaced by concrete type ContextKey
	DefaultContextKey            = "id"
	ContextKeyID      ContextKey = iota
)

// SetRequestId returns a shallow copy of the origin request, with a request ID created or extracted from the X-Request-ID header
func SetRequestId(r *http.Request) *http.Request {
	id := r.Header.Get(DefaultRequestHeader)
	// if we didn't get an id, create one
	if id == "" {
		id = uuid.New().String()
		log.WithField("id", id).Debugf("failed to locate existing request ID, generating a new one...")
	}
	// store the legacy value for compat
	ctx := context.WithValue(r.Context(), DefaultContextKey, id)
	// update the request context
	return r.WithContext(context.WithValue(ctx, ContextKeyID, id))
}

// GetRequestId extracts the request ID from the current request context
func GetRequestId(r *http.Request) string {
	return GetContextId(r.Context())
}

// GetContextId extracts the request ID from the given context
func GetContextId(ctx context.Context) string {
	id := ctx.Value(ContextKeyID)
	if id == nil {
		id = ctx.Value(DefaultContextKey)
	}
	if id == nil {
		return ""
	}
	// safely convert the id into a string
	return fmt.Sprintf("%v", id)
}
