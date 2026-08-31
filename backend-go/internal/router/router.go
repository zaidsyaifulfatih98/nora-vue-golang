package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/auth"
	"nora-photobooth-backend/internal/backdrops"
	"nora-photobooth-backend/internal/config"
	"nora-photobooth-backend/internal/features"
	"nora-photobooth-backend/internal/finance"
	"nora-photobooth-backend/internal/frametemplates"
	"nora-photobooth-backend/internal/gallery"
	"nora-photobooth-backend/internal/logging"
	appmw "nora-photobooth-backend/internal/middleware"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/packages"
	"nora-photobooth-backend/internal/photoboothframes"
	"nora-photobooth-backend/internal/photoboothresults"
	"nora-photobooth-backend/internal/reviews"
	"nora-photobooth-backend/internal/upload"
	"nora-photobooth-backend/internal/voicemessages"
)

func New(cfg *config.Config, db *gorm.DB, sheetsSvc *sheets.Service) (*gin.Engine, error) {
	r := gin.New()
	r.Use(appmw.Recovery())
	r.Use(logging.RequestLogger())
	r.Use(cors.New(config.BuildCORS(cfg)))
	r.Use(appmw.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": gin.H{}})
	})

	jwtVerify := appmw.JWTVerify(cfg.JWTSecret)
	adminOnly := appmw.RoleVerify(string(models.RoleSuperAdmin), string(models.RoleAdmin))
	superAdminOnly := appmw.RoleVerify(string(models.RoleSuperAdmin))

	api := r.Group("/api")

	// Auth
	authService := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/me", jwtVerify, authHandler.Me)
		authGroup.POST("/register", jwtVerify, superAdminOnly, authHandler.Register)
	}

	// Packages
	pkgHandler := packages.NewHandler(db)
	pkgGroup := api.Group("/packages")
	{
		pkgGroup.GET("", pkgHandler.List())
		pkgGroup.POST("", jwtVerify, adminOnly, pkgHandler.Create)
		pkgGroup.PATCH("/:id", jwtVerify, adminOnly, pkgHandler.Update)
		pkgGroup.DELETE("/:id", jwtVerify, adminOnly, pkgHandler.Delete())
	}

	// Reviews
	reviewHandler := reviews.NewHandler(db)
	reviewGroup := api.Group("/reviews")
	{
		reviewGroup.GET("", reviewHandler.List())
		reviewGroup.POST("", jwtVerify, adminOnly, reviewHandler.Create)
		reviewGroup.PATCH("/:id", jwtVerify, adminOnly, reviewHandler.Update)
		reviewGroup.DELETE("/:id", jwtVerify, adminOnly, reviewHandler.Delete())
	}

	uploader, err := upload.NewUploader(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		return nil, err
	}

	// Gallery
	galleryHandler := gallery.NewHandler(db, uploader)
	galleryGroup := api.Group("/gallery")
	{
		galleryGroup.GET("", galleryHandler.List())
		galleryGroup.POST("", jwtVerify, adminOnly, galleryHandler.Create)
		galleryGroup.PATCH("/:id", jwtVerify, adminOnly, galleryHandler.Update)
		galleryGroup.DELETE("/:id", jwtVerify, adminOnly, galleryHandler.Delete())
	}

	// Features
	featureHandler := features.NewHandler(db, uploader)
	featureGroup := api.Group("/features")
	{
		featureGroup.GET("", featureHandler.List())
		featureGroup.POST("", jwtVerify, adminOnly, featureHandler.Create)
		featureGroup.PATCH("/:id", jwtVerify, adminOnly, featureHandler.Update)
		featureGroup.DELETE("/:id", jwtVerify, adminOnly, featureHandler.Delete())
	}

	// Frame templates
	frameHandler := frametemplates.NewHandler(db, uploader)
	frameGroup := api.Group("/frame-templates")
	{
		frameGroup.GET("", frameHandler.List())
		frameGroup.POST("", jwtVerify, adminOnly, frameHandler.Create)
		frameGroup.PATCH("/:id", jwtVerify, adminOnly, frameHandler.Update)
		frameGroup.DELETE("/:id", jwtVerify, adminOnly, frameHandler.Delete())
	}

	// Photobooth frames (transparent-window PNGs used by the digital
	// photobooth "try it" flow — distinct from the landing page's frame templates)
	photoboothFrameHandler := photoboothframes.NewHandler(db, uploader)
	photoboothFrameGroup := api.Group("/photobooth-frames")
	{
		photoboothFrameGroup.GET("", photoboothFrameHandler.List())
		photoboothFrameGroup.POST("", jwtVerify, adminOnly, photoboothFrameHandler.Create)
		photoboothFrameGroup.PATCH("/:id", jwtVerify, adminOnly, photoboothFrameHandler.Update)
		photoboothFrameGroup.DELETE("/:id", jwtVerify, adminOnly, photoboothFrameHandler.Delete())
	}

	// Photobooth results — saves a guest's finished digital photobooth image
	// and returns a public download link (turned into a QR code on the
	// frontend). Public: guests use this, not admins.
	photoboothResultHandler := photoboothresults.NewHandler(db, uploader)
	photoboothResultGroup := api.Group("/photobooth-results")
	{
		photoboothResultGroup.POST("", photoboothResultHandler.Create)
		photoboothResultGroup.GET("", jwtVerify, adminOnly, photoboothResultHandler.List)
		photoboothResultGroup.DELETE("/:id", jwtVerify, adminOnly, photoboothResultHandler.Delete)
	}

	// Voice messages — guests record a greeting from the digital photobooth
	// result screen (public); admins listen to the collection in the
	// dashboard (auth required).
	voiceMessageHandler := voicemessages.NewHandler(db, uploader)
	voiceMessageGroup := api.Group("/voice-messages")
	{
		voiceMessageGroup.POST("", voiceMessageHandler.Create)
		voiceMessageGroup.GET("", jwtVerify, adminOnly, voiceMessageHandler.List)
		voiceMessageGroup.DELETE("/:id", jwtVerify, adminOnly, voiceMessageHandler.Delete)
	}

	// Backdrops
	backdropHandler := backdrops.NewHandler(db, uploader)
	backdropGroup := api.Group("/backdrops")
	{
		backdropGroup.GET("", backdropHandler.List())
		backdropGroup.POST("", jwtVerify, adminOnly, backdropHandler.Create)
		backdropGroup.PATCH("/:id", jwtVerify, adminOnly, backdropHandler.Update)
		backdropGroup.DELETE("/:id", jwtVerify, adminOnly, backdropHandler.Delete())
	}

	// Finance (Google Sheets backed, whole router requires auth)
	if sheetsSvc != nil {
		sheetsRepo := finance.NewSheetsRepository(sheetsSvc, cfg.GoogleSheetsSpreadsheetID)
		financeService := finance.NewService(db, sheetsRepo)
		financeHandler := finance.NewHandler(financeService)
		financeGroup := api.Group("/finance", jwtVerify, adminOnly)
		{
			financeGroup.GET("", financeHandler.List)
			financeGroup.GET("/summary", financeHandler.Summary)
			financeGroup.POST("", financeHandler.Create)
			financeGroup.PATCH("/:id", financeHandler.Update)
			financeGroup.DELETE("/:id", financeHandler.Delete)
		}
	}

	return r, nil
}
