package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	client *cloudinary.Cloudinary
}

var cloudinaryInstance *CloudinaryUploader

// InitCloudinary initializes the Cloudinary client
func InitCloudinary(cloudName, apiKey, apiSecret string) (*CloudinaryUploader, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	cloudinaryInstance = &CloudinaryUploader{
		client: cld,
	}
	return cloudinaryInstance, nil
}

// GetCloudinaryClient returns the singleton Cloudinary instance
func GetCloudinaryClient() *CloudinaryUploader {
	return cloudinaryInstance
}

// UploadEmployeeDocument uploads a base64 encoded document to Cloudinary
func (cu *CloudinaryUploader) UploadEmployeeDocument(ctx context.Context, documentBase64, documentName string, employeeID uint) (string, error) {
	if cu.client == nil {
		return "", fmt.Errorf("cloudinary client not initialized")
	}

	// Decode base64
	decodedData, err := base64.StdEncoding.DecodeString(documentBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		Folder:         "employee-documents",
		PublicID:       fmt.Sprintf("employee_%d_%s", employeeID, documentName),
		ResourceType:   "raw",
		UniqueFilename: &[]bool{false}[0],
	}

	uploadResult, err := cu.client.Upload.Upload(ctx, bytes.NewReader(decodedData), uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload document to cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

// UploadEmployeeDocumentFromReader uploads a document file to Cloudinary from an io.Reader
func (cu *CloudinaryUploader) UploadEmployeeDocumentFromReader(ctx context.Context, fileReader io.Reader, fileName string, employeeID uint) (string, error) {
	if cu.client == nil {
		return "", fmt.Errorf("cloudinary client not initialized")
	}

	// Remove file extension from fileName to avoid double extensions
	fileNameWithoutExt := fileName
	if lastDot := strings.LastIndex(fileName, "."); lastDot != -1 {
		fileNameWithoutExt = fileName[:lastDot]
	}

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		Folder:         "employee-documents",
		PublicID:       fmt.Sprintf("employee_%d_%s", employeeID, fileNameWithoutExt),
		ResourceType:   "raw",
		UniqueFilename: &[]bool{false}[0],
	}

	uploadResult, err := cu.client.Upload.Upload(ctx, fileReader, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload document to cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

// DeleteEmployeeDocument deletes an employee document from Cloudinary
func (cu *CloudinaryUploader) DeleteEmployeeDocument(ctx context.Context, publicID string) error {
	if cu.client == nil {
		return fmt.Errorf("cloudinary client not initialized")
	}

	_, err := cu.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete document from cloudinary: %w", err)
	}

	return nil
}
