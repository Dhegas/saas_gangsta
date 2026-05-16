package storage

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageRule struct {
	MaxSize     int64 // bytes
	MaxWidth    int
	Quality     int
	AllowedExt  []string
	AllowedMime []string
}

var ImageRules = map[string]ImageRule{
	"tenant_menus": {
		MaxSize:     1 * 1024 * 1024,
		MaxWidth:    1200,
		Quality:     75,
		AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp"},
		AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
	},
	"tenant_logos": {
		MaxSize:     2 * 1024 * 1024,
		MaxWidth:    800,
		Quality:     70,
		AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp", ".svg"},
		AllowedMime: []string{"image/jpeg", "image/png", "image/webp", "image/svg+xml"},
	},
	"user_profiles": {
		MaxSize:     2 * 1024 * 1024,
		MaxWidth:    500,
		Quality:     70,
		AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp"},
		AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
	},
	"tenant_banners": {
		MaxSize:     5 * 1024 * 1024,
		MaxWidth:    1920,
		Quality:     80,
		AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp"},
		AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
	},
}

type ImageProcessor interface {
	Process(fileHeader *multipart.FileHeader, folder string) (io.Reader, string, error)
}

type imageProcessor struct{}

func NewImageProcessor() ImageProcessor {
	return &imageProcessor{}
}

func (p *imageProcessor) Process(fileHeader *multipart.FileHeader, folder string) (io.Reader, string, error) {
	// 1. Get rules for folder
	rule, ok := ImageRules[strings.Trim(folder, "/")]
	if !ok {
		// Default rule if folder not found
		rule = ImageRule{
			MaxSize:     2 * 1024 * 1024,
			MaxWidth:    1024,
			Quality:     75,
			AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp"},
			AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
		}
	}

	// 2. Validate Size
	if fileHeader.Size > rule.MaxSize {
		return nil, "", fmt.Errorf("file size exceeds limit of %d bytes", rule.MaxSize)
	}

	// 3. Validate Extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	isAllowedExt := false
	for _, allowed := range rule.AllowedExt {
		if ext == allowed {
			isAllowedExt = true
			break
		}
	}
	if !isAllowedExt {
		return nil, "", fmt.Errorf("extension %s not allowed", ext)
	}

	// 4. Open File
	src, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 5. Validate Mime Type
	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		return nil, "", fmt.Errorf("failed to read file header: %w", err)
	}
	mimeType := http.DetectContentType(buff)
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	isAllowedMime := false
	for _, allowed := range rule.AllowedMime {
		if mimeType == allowed {
			isAllowedMime = true
			break
		}
	}
	if !isAllowedMime {
		return nil, "", fmt.Errorf("mime type %s not allowed", mimeType)
	}

	// 6. Special case for SVG - skip processing as it's vector
	if mimeType == "image/svg+xml" {
		return src, mimeType, nil
	}

	// 7. Decode Image
	img, err := imaging.Decode(src)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// 8. Resize if needed
	bounds := img.Bounds()
	if bounds.Dx() > rule.MaxWidth {
		img = imaging.Resize(img, rule.MaxWidth, 0, imaging.Lanczos)
	}

	// 9. Compress and Encode
	var out bytes.Buffer
	outputMime := mimeType

	switch mimeType {
	case "image/png":
		// PNG optimization
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		if err := encoder.Encode(&out, img); err != nil {
			return nil, "", fmt.Errorf("failed to encode png: %w", err)
		}
	default:
		// Default to JPEG for compression
		if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: rule.Quality}); err != nil {
			return nil, "", fmt.Errorf("failed to encode jpeg: %w", err)
		}
		outputMime = "image/jpeg"
	}

	return &out, outputMime, nil
}
