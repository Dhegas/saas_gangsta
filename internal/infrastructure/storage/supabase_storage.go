package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StorageService interface {
	UploadFile(ctx context.Context, bucket, folder, filename string, file io.Reader, contentType string) (string, error)
}

type supabaseStorage struct {
	baseURL string
	apiKey  string
}

func NewSupabaseStorage(url, key string) StorageService {
	// Menghapus trailing slash jika ada
	url = strings.TrimSuffix(url, "/")
	return &supabaseStorage{
		baseURL: fmt.Sprintf("%s/storage/v1/object", url),
		apiKey:  key,
	}
}

func (s *supabaseStorage) UploadFile(ctx context.Context, bucket, folder, filename string, file io.Reader, contentType string) (string, error) {
	// Menghapus leading/trailing slashes untuk konsistensi
	bucket = strings.Trim(bucket, "/")
	folder = strings.Trim(folder, "/")
	
	path := fmt.Sprintf("%s/%s", folder, filename)
	uploadURL := fmt.Sprintf("%s/%s/%s", s.baseURL, bucket, path)

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, file)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Content-Type", contentType)
	// Upsert: true memungkinkan menimpa file jika nama sama
	req.Header.Set("x-upsert", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase storage error: %s - %s", resp.Status, string(body))
	}

	// Untuk bucket publik, format URL-nya adalah:
	// [SUPABASE_URL]/storage/v1/object/public/[bucket]/[folder]/[filename]
	publicURL := fmt.Sprintf("%s/public/%s/%s", s.baseURL, bucket, path)
	return publicURL, nil
}
