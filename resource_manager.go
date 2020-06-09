package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/commercionetwork/didcomauth"

	"github.com/gorilla/mux"

	redisClient "github.com/go-redis/redis/v7"
)

const (
	rmKeyFmt      = "resourcemanager-%s"
	formFieldHash = "document_hash"
	formFieldCert = "document_cert"
)

type GetDocumentResponse struct {
	Hash string `json:"hash"`
	Cert string `json:"cert"`
}

type ResourceManager struct {
	redis    *redisClient.Client
	savePath string
}

func newResourceManager(rh, savePath string) *ResourceManager {
	return &ResourceManager{
		savePath: savePath,
		redis:    redisClient.NewClient(&redisClient.Options{Addr: rh}),
	}
}

func rmKey(did string) string {
	return fmt.Sprintf(rmKeyFmt, did)
}

func (r *ResourceManager) Add(did, resource string) {
	r.redis.Set(rmKey(did), resource, 0)
}

func (r *ResourceManager) Get(did string) string {
	return r.redis.Get(rmKey(did)).Val()
}

func (r *ResourceManager) Del(did string) {
	r.redis.Del(rmKey(did))
}

func uploadSavePath(base, resource string) string {
	return filepath.Join(base, resource)
}

func (r *ResourceManager) HandleUpload(rw http.ResponseWriter, req *http.Request) {
	// resource
	vars := mux.Vars(req)
	did := req.Header.Get(didcomauth.DIDHeader)
	res := vars["id"]

	if path.Base(r.Get(did)) != res {
		log.Printf("attempted access to resource %s by unauthorized did %s\n", res, did)
		writeError(rw, http.StatusForbidden, errors.New("access to resource denied"))
		return
	}

	defer r.Del(did)

	// get form data
	err := req.ParseForm()
	if err != nil {
		writeError(rw, http.StatusBadRequest, errors.New("malformed form data"))
		return
	}

	data := req.PostForm.Encode()
	savePath := uploadSavePath(r.savePath, res)

	_, err = os.Stat(savePath)
	if !os.IsNotExist(err) {
		log.Printf("savePath resource %s has already been used before\n", savePath)
		writeError(rw, http.StatusBadRequest, errors.New("invalid resource"))
		return
	}

	err = ioutil.WriteFile(savePath, []byte(data), 0755)
	if err != nil {
		log.Println(fmt.Errorf("could not save cert/signature in storage, %w", err))
		writeError(rw, http.StatusInternalServerError, errors.New("could not save file"))
		return
	}
}

func (r *ResourceManager) HandleGetDocument(rw http.ResponseWriter, req *http.Request) {
	// resource
	vars := mux.Vars(req)
	res := vars["id"]

	savePath := uploadSavePath(r.savePath, res)
	rawData, err := ioutil.ReadFile(savePath)
	if err != nil {
		log.Printf("access to unknown id %s has been tried, refusing", res)
		writeError(rw, http.StatusForbidden, errors.New("forbidden"))
		return
	}

	defer os.Remove(savePath)

	values, err := url.ParseQuery(string(rawData))
	if err != nil {
		log.Printf("could not parse query string from file, %s", err.Error())
		writeError(rw, http.StatusInternalServerError, errors.New("error while reading data from disk"))
		return
	}

	ret := GetDocumentResponse{
		Hash: values.Get(formFieldHash),
		Cert: values.Get(formFieldCert),
	}

	jenc := json.NewEncoder(rw)
	err = jenc.Encode(ret)

	if err != nil {
		log.Printf("could not marshal response to getdocument, %s/n", err.Error())
		writeError(rw, http.StatusInternalServerError, errors.New("could not marshal response"))
		return
	}
}

func (r *ResourceManager) HandleAdd(rw http.ResponseWriter, req *http.Request) {
	// resource
	did := req.Header.Get(didcomauth.DIDHeader)
	res := req.Header.Get(didcomauth.ResourceHeader)

	if did == "" {
		log.Println("empty did in Add request")
		writeError(rw, http.StatusBadRequest, errors.New("invalid did"))
		return
	}

	if res == "" {
		log.Println("empty res in Add request")
		writeError(rw, http.StatusBadRequest, errors.New("invalid res"))
		return
	}

	savePath := uploadSavePath(r.savePath, res)

	_, err := os.Stat(savePath)
	if !os.IsNotExist(err) {
		log.Printf("savePath resource %s has already been used before\n", savePath)
		writeError(rw, http.StatusBadRequest, errors.New("invalid resource"))
		return
	}

	r.Add(did, res)
}
