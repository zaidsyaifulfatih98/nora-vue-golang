package upload

import (
	"context"
	"errors"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"nora-photobooth-backend/internal/apperror"
)

// cloudinaryPublicIDPattern pulls the public_id (including any folder
// prefix, e.g. "uploads/abcd1234") out of a Cloudinary secure_url. Every
// asset here is uploaded without a custom public_id, so Cloudinary always
// shapes the URL as .../upload/[v<version>/]<public_id>.<ext> — matching
// this lets existing records that only ever stored the URL (not the
// public_id separately) still be deleted from Cloudinary itself.
var cloudinaryPublicIDPattern = regexp.MustCompile(`/upload/(?:v\d+/)?(.+)\.[a-zA-Z0-9]+$`)

func ExtractPublicID(secureURL string) (string, error) {
	matches := cloudinaryPublicIDPattern.FindStringSubmatch(secureURL)
	if len(matches) != 2 {
		return "", errors.New("not a recognizable Cloudinary upload URL: " + secureURL)
	}
	return matches[1], nil
}

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

// DeleteImage removes a previously-uploaded image from Cloudinary, deriving
// its public_id from the stored secure_url. Best-effort by design: callers
// should log/ignore a failure here rather than block the user-facing delete
// on Cloudinary being reachable — an orphaned Cloudinary asset is a much
// smaller problem than a delete button that stops working.
func (u *Uploader) DeleteImage(secureURL string) error {
	publicID, err := ExtractPublicID(secureURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = u.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}

// DeleteAudio removes a previously-uploaded voice message recording.
// ResourceType "video" must match what UploadAudio used, or Cloudinary
// won't find the asset to destroy.
func (u *Uploader) DeleteAudio(secureURL string) error {
	publicID, err := ExtractPublicID(secureURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = u.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID, ResourceType: "video"})
	return err
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
