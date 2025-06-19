package handlers

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadPosterToCloudinary(file multipart.File, fileHeader *multipart.FileHeader, cfgCloudName, cfgAPIKey, cfgAPISecret string) (string, error) {
	cld, err := cloudinary.NewFromParams(cfgCloudName, cfgAPIKey, cfgAPISecret)
	if err != nil {
		return "", err
	}
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicID:       fileHeader.Filename,
		Folder:         "movie_posters",
		UniqueFilename: func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}
