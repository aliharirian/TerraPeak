package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aliharirian/TerraPeak/logger"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetVersionList(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "name")

	// Generate cache key for this request
	cacheKey := fmt.Sprintf("registry/v1/versions/%s/%s", namespace, name)

	// Check if response exists in cache
	if cachedResponse := s.getCachedResponse(cacheKey); cachedResponse != nil {
		logger.Infof("Cache HIT: Serving cached version list for %s/%s", namespace, name)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache-Status", "HIT")
		w.WriteHeader(http.StatusOK)
		w.Write(cachedResponse)
		return
	}

	// Cache miss - fetch from upstream
	logger.Infof("Cache MISS: Fetching version list for %s/%s from upstream", namespace, name)
	upstreamURL := s.cfg.Terraform.RegistryUrl + "/v1/providers/" + namespace + "/" + name + "/versions"

	resp, err := http.Get(upstreamURL)
	if err != nil {
		http.Error(w, "upstream unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	// Cache the response
	s.cacheResponse(cacheKey, respBody)

	// Send response to client
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("X-Cache-Status", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (s *Service) GetProviderDownloadDetails(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "name")
	version := chi.URLParam(r, "version")
	os := chi.URLParam(r, "os")
	arch := chi.URLParam(r, "arch")

	// Generate cache key for this request
	cacheKey := fmt.Sprintf("registry/v1/download/%s/%s/%s/%s/%s", namespace, name, version, os, arch)

	// Check if response exists in cache
	if cachedResponse := s.getCachedResponse(cacheKey); cachedResponse != nil {
		logger.Infof("Cache HIT: Serving cached download details for %s/%s/%s/%s/%s", namespace, name, version, os, arch)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache-Status", "HIT")
		w.WriteHeader(http.StatusOK)
		w.Write(cachedResponse)
		return
	}

	// Cache miss - fetch from upstream
	logger.Infof("Cache MISS: Fetching download details for %s/%s/%s/%s/%s from upstream", namespace, name, version, os, arch)
	upstreamURL := s.cfg.Terraform.RegistryUrl + "/v1/providers/" + namespace + "/" + name + "/" + version + "/download/" + os + "/" + arch

	resp, err := http.Get(upstreamURL)
	if err != nil {
		http.Error(w, "upstream unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	var body map[string]any
	if err := json.Unmarshal(respBody, &body); err != nil {
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	// Modify URLs to point to our cacher
	body["download_url"] = AppFirstURL(body["download_url"], s.cfg.Server.Domain)
	body["shasums_signature_url"] = AppFirstURL(body["shasums_signature_url"], s.cfg.Server.Domain)
	body["shasums_url"] = AppFirstURL(body["shasums_url"], s.cfg.Server.Domain)

	// Encode modified response
	modifiedResponse, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	// Cache the modified response
	s.cacheResponse(cacheKey, modifiedResponse)

	// Send response to client
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("X-Cache-Status", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(modifiedResponse)
}

func AppFirstURL(base any, cacherURL string) any {
	// Parse the base URL
	baseURL, err := url.Parse(base.(string))
	if err != nil {
		return base // Return original if parsing fails
	}

	// Parse the cacher URL
	cacherURLParsed, err := url.Parse(cacherURL)
	if err != nil {
		return base // Return original if parsing fails
	}

	// Create new URL: cacher host + original host + original path
	newPath := "/" + baseURL.Host + baseURL.Path
	newURL := &url.URL{
		Scheme: cacherURLParsed.Scheme,
		Host:   cacherURLParsed.Host,
		Path:   newPath,
	}

	return newURL.String()
}

// getCachedResponse retrieves cached API response from storage
func (s *Service) getCachedResponse(cacheKey string) []byte {
	if s.store == nil {
		return nil
	}

	// Check if file exists in storage
	if !s.store.FileExists(cacheKey) {
		return nil
	}

	// Read from storage
	data, err := s.store.ReadFromStorage(cacheKey)
	if err != nil {
		logger.Debugf("Failed to read cached response for %s: %v", cacheKey, err)
		return nil
	}

	return data
}

// cacheResponse stores API response in storage
func (s *Service) cacheResponse(cacheKey string, data []byte) {
	if s.store == nil {
		return
	}

	err := s.store.Save(cacheKey, data)
	if err != nil {
		logger.Warnf("Failed to cache response for %s: %v", cacheKey, err)
	} else {
		logger.Debugf("Successfully cached response for %s (%d bytes)", cacheKey, len(data))
	}
}
