package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/commercionetwork/didcomauth"
)

func main() {
	rawKey, err := ioutil.ReadFile("./private_signing_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(rawKey)

	pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	rpk := pk.(*rsa.PrivateKey)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:9999/auth/challenge", nil)
	req.Header.Set("X-DID", "did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	req.Header.Set("X-Resource", "/protected/upload/testing")

	c := &http.Client{}
	data, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	jdec := json.NewDecoder(data.Body)

	var cc didcomauth.Challenge
	err = jdec.Decode(&cc)
	if err != nil {
		log.Fatal(err)
	}

	sp := cc.SignaturePayload()
	spH := sha256.Sum256(sp)

	rbytes, err := rsa.SignPKCS1v15(rand.Reader, rpk, crypto.SHA256, spH[:])
	if err != nil {
		log.Fatal(err)
	}

	response := base64.StdEncoding.EncodeToString(rbytes)

	rr := didcomauth.AuthResponse{
		Challenge: cc,
		Response:  response,
	}

	rrb, _ := json.Marshal(rr)

	req, _ = http.NewRequest(http.MethodPost, "http://localhost:9999/auth/challenge", bytes.NewReader(rrb))
	req.Header.Set("X-DID", "did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	req.Header.Set("X-Resource", "/protected/upload/testing")

	d2, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if d2.StatusCode != http.StatusOK {
		d, _ := ioutil.ReadAll(d2.Body)
		log.Fatal(string(d))
	}

	jdec = json.NewDecoder(d2.Body)
	var j didcomauth.ReleaseJWTResponse
	err = jdec.Decode(&j)
	if err != nil {
		log.Fatal(err)
	}

	uploadPayload := url.Values{}
	uploadPayload.Set("document_hash", "hash")
	uploadPayload.Set("document_cert", "cert")

	req, _ = http.NewRequest(http.MethodPost, "http://localhost:9999/protected/upload/testing", strings.NewReader(uploadPayload.Encode()))
	req.Header.Set("X-DID", "did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	req.Header.Set("X-Resource", "/protected/upload/testing")
	req.Header.Set("Authorization", "Bearer "+j.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	d2, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	aaa, _ := ioutil.ReadAll(d2.Body)
	fmt.Println(string(aaa))

}
