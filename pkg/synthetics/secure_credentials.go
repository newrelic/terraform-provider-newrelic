package synthetics

import (
	"fmt"
)

// GetSecureCredentials is used to retrieve all secure credentials from your New Relic account.
func (s *Synthetics) GetSecureCredentials() ([]*SecureCredential, error) {
	url := "/v1/secure-credentials"
	resp := getSecureCredentialsResponse{}

	_, err := s.client.Get(url, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.SecureCredentials, nil
}

// GetSecureCredential is used to retrieve a specific secure credential from your New Relic account.
func (s *Synthetics) GetSecureCredential(key string) (*SecureCredential, error) {
	url := fmt.Sprintf("/v1/secure-credentials/%s", key)
	var sc SecureCredential

	_, err := s.client.Get(url, nil, &sc)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

// AddSecureCredential is used to add a secure credential to your New Relic account.
func (s *Synthetics) AddSecureCredential(key, value, description string) (*SecureCredential, error) {
	url := "/v1/secure-credentials"
	sc := &SecureCredential{
		Key:         key,
		Value:       value,
		Description: description,
	}

	_, err := s.client.Post(url, nil, sc, nil)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

// UpdateSecureCredential is used to update a secure credential in your New Relic account.
func (s *Synthetics) UpdateSecureCredential(key, value, description string) (*SecureCredential, error) {
	url := fmt.Sprintf("/v1/secure-credentials/%s", key)
	sc := &SecureCredential{
		Key:         key,
		Value:       value,
		Description: description,
	}

	_, err := s.client.Put(url, nil, sc, nil)

	if err != nil {
		return nil, err
	}

	return sc, nil
}

// DeleteSecureCredential deletes a secure credential from your New Relic account.
func (s *Synthetics) DeleteSecureCredential(key string) error {
	url := fmt.Sprintf("/v1/secure-credentials/%s", key)

	_, err := s.client.Delete(url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

type getSecureCredentialsResponse struct {
	SecureCredentials []*SecureCredential `json:"secureCredentials"`
	Count             int                 `json:"count"`
}
