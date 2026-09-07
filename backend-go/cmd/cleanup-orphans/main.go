// One-off report of Cloudinary assets orphaned before the delete/update
// handlers started cleaning up after themselves: cross-checks every
// Cloudinary asset under the "uploads" folder against every URL still
// referenced in the database, and prints whatever isn't referenced
// anywhere. Read-only — it never deletes anything; use the printed
// public_id list to remove them manually from the Cloudinary dashboard
// (Media Library → search the public_id → delete).
//
//	go run ./cmd/cleanup-orphans
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/config"
	"nora-photobooth-backend/internal/logging"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/upload"
)

func main() {
	logging.Init()
	cfg := config.Load()

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to init cloudinary client")
	}

	knownImageIDs, knownAudioIDs := collectKnownPublicIDs(db)
	fmt.Printf("Found %d referenced image(s) and %d referenced audio file(s) in the database.\n\n", len(knownImageIDs), len(knownAudioIDs))

	ctx := context.Background()
	imageOrphans := findOrphans(ctx, cld, api.Image, "uploads", knownImageIDs)
	audioOrphans := findOrphans(ctx, cld, api.Video, "uploads/voice-messages", knownAudioIDs)

	orphans := append(imageOrphans, audioOrphans...)
	if len(orphans) == 0 {
		fmt.Println("No orphaned Cloudinary assets found. Nothing to do.")
		return
	}

	var totalBytes int
	for _, o := range orphans {
		totalBytes += o.Bytes
		fmt.Printf("  orphan: %-45s  %8.1f KB  uploaded %s\n", o.PublicID, float64(o.Bytes)/1024, o.CreatedAt.Format("2006-01-02"))
	}
	fmt.Printf("\n%d orphaned asset(s), %.1f MB total.\n", len(orphans), float64(totalBytes)/1024/1024)
	fmt.Println("\nThis tool is read-only — nothing was deleted. Remove these manually from the Cloudinary dashboard (Media Library → search the public_id) if you want to reclaim the space.")
}

// collectKnownPublicIDs reads every Cloudinary URL still referenced by a
// live (non-soft-deleted) row across every resource that uploads to
// Cloudinary, and converts each one to the public_id Cloudinary actually
// tracks it under.
func collectKnownPublicIDs(db *gorm.DB) (images map[string]bool, audio map[string]bool) {
	images = map[string]bool{}
	audio = map[string]bool{}

	addImage := func(url string) {
		if url == "" {
			return
		}
		if id, err := upload.ExtractPublicID(url); err == nil {
			images[id] = true
		}
	}
	addAudio := func(url string) {
		if url == "" {
			return
		}
		if id, err := upload.ExtractPublicID(url); err == nil {
			audio[id] = true
		}
	}

	var features []models.Feature
	db.Find(&features)
	for _, f := range features {
		addImage(f.ImageURL)
	}

	var frameTemplates []models.FrameTemplate
	db.Find(&frameTemplates)
	for _, f := range frameTemplates {
		addImage(f.ImageURL)
	}

	var backdrops []models.Backdrop
	db.Find(&backdrops)
	for _, b := range backdrops {
		addImage(b.ImageURL)
	}

	var galleryPhotos []models.GalleryPhoto
	db.Find(&galleryPhotos)
	for _, g := range galleryPhotos {
		addImage(g.URL)
	}

	var photoboothFrames []models.PhotoboothFrame
	db.Find(&photoboothFrames)
	for _, p := range photoboothFrames {
		addImage(p.ImageURL)
	}

	var photoboothResults []models.PhotoboothResult
	db.Find(&photoboothResults)
	for _, p := range photoboothResults {
		addImage(p.ImageURL)
	}

	var voiceMessages []models.VoiceMessage
	db.Find(&voiceMessages)
	for _, v := range voiceMessages {
		addAudio(v.AudioURL)
		// A voice message's PhotoURL points at a PhotoboothResult's own
		// asset, already covered above — not a separate upload.
	}

	return images, audio
}

type orphanAsset struct {
	PublicID  string
	AssetType string
	Bytes     int
	CreatedAt time.Time
}

// findOrphans pages through every Cloudinary asset under prefix (for the
// given resource type) and returns the ones whose public_id isn't in known.
func findOrphans(ctx context.Context, cld *cloudinary.Cloudinary, assetType api.AssetType, prefix string, known map[string]bool) []orphanAsset {
	var orphans []orphanAsset
	cursor := ""

	for {
		res, err := cld.Admin.Assets(ctx, admin.AssetsParams{
			AssetType:    assetType,
			DeliveryType: "upload",
			Prefix:       prefix,
			MaxResults:   500,
			NextCursor:   cursor,
		})
		if err != nil {
			fmt.Printf("failed to list %s assets under %q: %v\n", assetType, prefix, err)
			return orphans
		}
		if res.Error.Message != "" {
			fmt.Printf("Cloudinary API error listing %s assets under %q: %s\n", assetType, prefix, res.Error.Message)
			return orphans
		}

		for _, asset := range res.Assets {
			if !known[asset.PublicID] {
				orphans = append(orphans, orphanAsset{
					PublicID:  asset.PublicID,
					AssetType: asset.AssetType,
					Bytes:     asset.Bytes,
					CreatedAt: asset.CreatedAt,
				})
			}
		}

		if res.NextCursor == "" {
			break
		}
		cursor = res.NextCursor
	}

	return orphans
}
