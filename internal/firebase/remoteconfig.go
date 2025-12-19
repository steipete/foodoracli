package firebase

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type APKFirebaseConfig struct {
	APIKey      string
	ProjectID   string
	ProjectNum  string
	AppID       string
	PackageName string
	CertSHA1    string
}

// HU values, extracted from the provided XAPK.
var NetPincerHU = APKFirebaseConfig{
	APIKey:      "AIzaSyBgXz7y2YA0gMbh1TLng4_o2gw2V40ug78",
	ProjectID:   "netpincer-1239",
	ProjectNum:  "954610850075",
	AppID:       "1:954610850075:android:289a63ab386cb722",
	PackageName: "hu.viala.newiapp",
	// SHA1 of signing cert (no colons): from META-INF/BNDLTOOL.RSA
	CertSHA1: "D106C8D37D4BF7D6F9DA124821503EE89CCB073A",
}

type RemoteConfigClient struct {
	cfg  APKFirebaseConfig
	http *http.Client
}

func NewRemoteConfigClient(cfg APKFirebaseConfig) *RemoteConfigClient {
	return &RemoteConfigClient{
		cfg: cfg,
		http: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

type Installation struct {
	FID       string
	AuthToken string
}

func (c *RemoteConfigClient) Fetch(ctx context.Context) (RemoteConfigFetchResponse, error) {
	inst, err := c.createInstallation(ctx)
	if err != nil {
		return RemoteConfigFetchResponse{}, err
	}
	return c.fetchRemoteConfig(ctx, inst)
}

type RemoteConfigFetchResponse struct {
	State           string            `json:"state"`
	Entries         map[string]string `json:"entries"`
	TemplateVersion string            `json:"templateVersion"`
}

func (c *RemoteConfigClient) createInstallation(ctx context.Context) (Installation, error) {
	fid, err := newFID()
	if err != nil {
		return Installation{}, err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   "firebaseinstallations.googleapis.com",
		Path:   fmt.Sprintf("/v1/projects/%s/installations", c.cfg.ProjectNum),
	}
	q := u.Query()
	q.Set("key", c.cfg.APIKey)
	u.RawQuery = q.Encode()

	body := map[string]any{
		"fid":         fid,
		"appId":       c.cfg.AppID,
		"authVersion": "FIS_v2",
		"sdkVersion":  "a:19.0.0",
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return Installation{}, err
	}
	c.addAndroidKeyRestrictionHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return Installation{}, err
	}
	defer res.Body.Close()
	rb, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Installation{}, fmt.Errorf("firebase installations: http %d: %s", res.StatusCode, string(rb))
	}

	var out struct {
		FID       string `json:"fid"`
		AuthToken struct {
			Token string `json:"token"`
		} `json:"authToken"`
	}
	if err := json.Unmarshal(rb, &out); err != nil {
		return Installation{}, fmt.Errorf("firebase installations: decode: %w", err)
	}
	if out.FID == "" || out.AuthToken.Token == "" {
		return Installation{}, fmt.Errorf("firebase installations: missing fid/authToken")
	}
	return Installation{FID: out.FID, AuthToken: out.AuthToken.Token}, nil
}

func (c *RemoteConfigClient) fetchRemoteConfig(ctx context.Context, inst Installation) (RemoteConfigFetchResponse, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "firebaseremoteconfig.googleapis.com",
		Path:   fmt.Sprintf("/v1/projects/%s/namespaces/firebase:fetch", c.cfg.ProjectNum),
	}
	q := u.Query()
	q.Set("key", c.cfg.APIKey)
	u.RawQuery = q.Encode()

	// Best-effort minimal client payload; server accepts additional fields.
	reqBody := map[string]any{
		"appId":       c.cfg.AppID,
		"packageName": c.cfg.PackageName,
		// FIS is used via headers; keep instance id fields for compatibility.
		"appInstanceId": inst.FID,
		"languageCode":  "en",
		"timeZone":      time.Now().Location().String(),
	}
	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return RemoteConfigFetchResponse{}, err
	}
	c.addAndroidKeyRestrictionHeaders(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.cfg.APIKey)
	req.Header.Set("X-Goog-Firebase-Installations-Id", inst.FID)
	req.Header.Set("X-Goog-Firebase-Installations-Auth", inst.AuthToken)

	res, err := c.http.Do(req)
	if err != nil {
		return RemoteConfigFetchResponse{}, err
	}
	defer res.Body.Close()
	rb, _ := io.ReadAll(io.LimitReader(res.Body, 4<<20))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return RemoteConfigFetchResponse{}, fmt.Errorf("firebase remote config: http %d: %s", res.StatusCode, string(rb))
	}

	var out RemoteConfigFetchResponse
	if err := json.Unmarshal(rb, &out); err != nil {
		return RemoteConfigFetchResponse{}, fmt.Errorf("firebase remote config: decode: %w", err)
	}
	if out.Entries == nil {
		out.Entries = map[string]string{}
	}
	return out, nil
}

func (c *RemoteConfigClient) addAndroidKeyRestrictionHeaders(req *http.Request) {
	if c.cfg.PackageName != "" {
		req.Header.Set("X-Android-Package", c.cfg.PackageName)
	}
	if c.cfg.CertSHA1 != "" {
		req.Header.Set("X-Android-Cert", c.cfg.CertSHA1)
	}
}

func newFID() (string, error) {
	// 16 bytes => 22 chars base64url (no padding), matches typical FID length.
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}
