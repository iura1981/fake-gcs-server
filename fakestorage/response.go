// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fakestorage

import (
	cloudstorage "cloud.google.com/go/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
	"net/http"
	"sort"
)

type listResponse struct {
	Kind     string        `json:"kind"`
	Items    []interface{} `json:"items"`
	Prefixes []string      `json:"prefixes"`
}

func newListBucketsResponse(bucketNames []string) listResponse {
	resp := listResponse{
		Kind:  "storage#buckets",
		Items: make([]interface{}, len(bucketNames)),
	}
	sort.Strings(bucketNames)
	for i, name := range bucketNames {
		resp.Items[i] = newBucketResponse(name)
	}
	return resp
}

type bucketResponse struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newBucketResponse(bucketName string) bucketResponse {
	return bucketResponse{
		Kind: "storage#bucket",
		ID:   bucketName,
		Name: bucketName,
	}
}

func newListObjectsResponse(objs []Object, prefixes []string) listResponse {
	resp := listResponse{
		Kind:     "storage#objects",
		Items:    make([]interface{}, len(objs)),
		Prefixes: prefixes,
	}
	for i, obj := range objs {
		resp.Items[i] = newObjectResponse(obj)
	}
	return resp
}

type objectResponse struct {
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Bucket      string `json:"bucket"`
	Size        int64  `json:"size,string"`
	ContentType string `json:"contentType,omitempty"`
	Crc32c      string `json:"crc32c,omitempty"`
	ACL         string `json:"acl,omitempty"`
	Md5Hash     string `json:"md5hash,omitempty"`
}

func newObjectResponse(obj Object) objectResponse {
	return objectResponse{
		Kind:        "storage#object",
		ID:          obj.id(),
		Bucket:      obj.BucketName,
		Name:        obj.Name,
		Size:        int64(len(obj.Content)),
		ContentType: obj.ContentType,
		Crc32c:      obj.Crc32c,
		ACL:         string(obj.ACL),
		Md5Hash:     obj.Md5Hash,
	}
}

type aclListResponse struct {
	*storage.ObjectAccessControls
}

func newACLListResponse(obj Object) aclListResponse {
	return aclListResponse{
		&storage.ObjectAccessControls{
			ServerResponse: googleapi.ServerResponse{
				HTTPStatusCode: http.StatusOK,
			},
			Items: []*storage.ObjectAccessControl{
				{
					Bucket: obj.BucketName,
					Entity: string(obj.ACL),
					Object: obj.Name,
					Role:   string(cloudstorage.RoleReader),
				},
			},
		},
	}
}

type rewriteResponse struct {
	Kind                string         `json:"kind"`
	TotalBytesRewritten int64          `json:"totalBytesRewritten,string"`
	ObjectSize          int64          `json:"objectSize,string"`
	Done                bool           `json:"done"`
	RewriteToken        string         `json:"rewriteToken"`
	Resource            objectResponse `json:"resource"`
}

func newObjectRewriteResponse(obj Object) rewriteResponse {
	return rewriteResponse{
		Kind:                "storage#rewriteResponse",
		TotalBytesRewritten: int64(len(obj.Content)),
		ObjectSize:          int64(len(obj.Content)),
		Done:                true,
		RewriteToken:        "",
		Resource:            newObjectResponse(obj),
	}
}

type errorResponse struct {
	Error httpError `json:"error"`
}

type httpError struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Errors  []apiError `json:"errors"`
}

type apiError struct {
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func newErrorResponse(code int, message string, errs []apiError) errorResponse {
	return errorResponse{
		Error: httpError{
			Code:    code,
			Message: message,
			Errors:  errs,
		},
	}
}
