package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/mattrax/Mattrax/mdm/agent"
)

func enroll(upn string) error {
	if info, err := os.Stat(AgentConfigPath); info != nil && !info.IsDir() {
		return fmt.Errorf("error uneroll before enrolling in a new server")
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error retrieving existance of current agent configuration: %w", err)
	}

	managementServerRoot, err := parseUPN(upn)
	if err != nil {
		return err
	}
	managementServerRoot = "enterpriseenrollment." + managementServerRoot

	enrollEndpoint := (&url.URL{
		Scheme: "https",
		Host:   managementServerRoot,
		Path:   "/Manage/Enroll.svc",
	}).String()

	res, err := client.Get(enrollEndpoint)
	if err != nil {
		return fmt.Errorf("error discovering mdm server: %w", err)
	} else if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error mdm server not found")
	}

	// TODO: Authentication phase

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error retrieving devices hostname: %w", err)
	}

	udid, err := machineid.ProtectedID(managementServerRoot)
	if err != nil {
		return fmt.Errorf("error retrieving unique device id: %w", err)
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("error generating device private key: %w", err)
	}

	csrRaw, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: udid,
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}, certKey)
	if err != nil {
		return fmt.Errorf("error retrieving unique device id: %w", err)
	}

	var enrollmentReq = agent.EnrollRequest{
		UDID:               udid,
		Hostname:           hostname,
		CertificateRequest: csrRaw,
	}

	enrollmentReqBody, err := json.Marshal(enrollmentReq)
	if err != nil {
		return fmt.Errorf("error marshalling enrollment json body: %w", err)
	}
	req, err := http.NewRequest("POST", enrollEndpoint, bytes.NewBuffer(enrollmentReqBody))
	req.Header.Set("Content-Type", "application/json+dm")

	res, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error enrolling with mdm server: %w", err)
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error mdm server had trouble enrolling. returned status: %v", res.StatusCode)
	} else if res.Header.Get("Content-Type") != "application/json+dm" {
		return fmt.Errorf("error mdm server response invalid. returned content type: %v", res.Header.Get("Content-Type"))
	}

	var enrollRes agent.EnrollResponse
	if err := json.NewDecoder(res.Body).Decode(&enrollRes); err != nil {
		return fmt.Errorf("error decoding mdm server enroll response: %w", err)
	}

	var agentConfig = AgentConfig{enrollRes, base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(certKey))}
	agentConfigFile, err := os.Create(AgentConfigPath)
	if err != nil {
		return fmt.Errorf("error saving agent configuration: %w", err)
	}
	if err := json.NewEncoder(agentConfigFile).Encode(agentConfig); err != nil {
		return fmt.Errorf("error marshalling agent configuration: %w", err)
	}
	agentConfigFile.Close()

	return nil
}

func parseUPN(upn string) (domain string, err error) {
	at := strings.LastIndex(upn, "@")
	if at >= 0 {
		return upn[at+1:], nil
	} else {
		return "", fmt.Errorf("error invalid upn format: %v", upn)
	}
}
