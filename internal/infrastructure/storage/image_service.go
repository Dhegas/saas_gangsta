package storage

import (
	"context"
	"fmt"
	"mime/multipart"
)

type ImageService interface {
	UploadOptimizedImage(ctx context.Context, bucket, folder string, fileHeader *multipart.FileHeader) (string, error)
}

type imageService struct {
	storage   StorageService
	processor ImageProcessor
}

func NewImageService(storage StorageService, processor ImageProcessor) ImageService {
	return &imageService{
		storage:   storage,
		processor: processor,
	}
}

func (s *imageService) UploadOptimizedImage(ctx context.Context, bucket, folder string, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("file header is nil")
	}

	// 1. Process image (validate, resize, compress)
	processedFile, mimeType, err := s.processor.Process(fileHeader, folder)
	if err != nil {
		return "", fmt.Errorf("failed to process image for folder %s: %w", folder, err)
	}

	// 2. Upload to storage
	// We use the same filename but the content is optimized
	publicURL, err := s.storage.UploadFile(ctx, bucket, folder, fileHeader.Filename, processedFile, mimeType)
	if err != nil {
		return "", fmt.Errorf("failed to upload optimized image: %w", err)
	}

	return publicURL, nil
}
