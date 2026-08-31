package upload

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"nora-photobooth-backend/internal/apperror"
)

type Uploader struct {
	cld *cloudinary.Cloudinary
}

func NewUploader(cloudName, apiKey, apiSecret string) (*Uploader, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &Uploader{cld: cld}, nil
}

// UploadImage mirrors cloudinary.utils.ts's cloudinaryUpload: uploads to the
// "uploads" folder with a 120s timeout, mapping timeout -> 504 and any other
// failure -> 502.
func (u *Uploader) UploadImage(file multipart.File) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	result, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "uploads"})
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", apperror.New("Upload gambar timeout", 504)
		}
		return "", apperror.New("Gagal mengunggah gambar", 502)
	}

	// The SDK's Upload() silently discards API-level errors (e.g. invalid/
	// missing Cloudinary credentials) instead of returning them as `err`, so
	// we must check the response body ourselves: a failed upload comes back
	// with `result.Error.Message` set and no `secure_url`.
	if result.Error.Message != "" || result.SecureURL == "" {
		msg := result.Error.Message
		if msg == "" {
			msg = "empty response from Cloudinary"
		}
		return "", apperror.New("Gagal mengunggah gambar: "+msg, 502)
	}

	return result.SecureURL, nil
}

// UploadAudio uploads a voice message recording. ResourceType "video" is
// Cloudinary's container for audio-only files (there is no separate "audio"
// resource type).
func (u *Uploader) UploadAudio(file multipart.File) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	result, err := u.cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "uploads/voice-messages", ResourceType: "video"})
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", apperror.New("Upload audio timeout", 504)
		}
		return "", apperror.New("Gagal mengunggah audio", 502)
	}

	if result.Error.Message != "" || result.SecureURL == "" {
		msg := result.Error.Message
		if msg == "" {
			msg = "empty response from Cloudinary"
		}
		return "", apperror.New("Gagal mengunggah audio: "+msg, 502)
	}

	return result.SecureURL, nil
}
