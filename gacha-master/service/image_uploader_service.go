package service

import (
	"cloud.google.com/go/storage"
	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"context"
	"errors"
	"fmt"
	"gacha-master/helper"
	"gacha-master/model/web"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"time"
)

type UploaderService interface {
	UploadCharacterImage(ctx context.Context, request *web.ImageCharacterUploadRequest) string
	DeleteCharacterImage(ctx context.Context, characterId int, gachaSystemId int)
	DeleteGachaSystemCharacterImage(ctx context.Context, gachaSystemId int)
	Close()
}

type UploaderServiceImpl struct {
	client        *storage.Client
	clientControl *control.StorageControlClient
	projectID     string
	bucketName    string
}

func NewUploaderServiceImpl() UploaderService {
	ctx := context.Background()
	var client *storage.Client
	var clientControl *control.StorageControlClient
	var err error

	if os.Getenv("ENVIRONMENT") != "prod" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
		if err != nil {
			panic(err)
		}

		clientControl, err = control.NewStorageControlClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
		if err != nil {
			panic(err)
		}
	} else {
		client, err = storage.NewClient(ctx)
		if err != nil {
			panic(err)
		}

		clientControl, err = control.NewStorageControlClient(ctx)
		if err != nil {
			panic(err)
		}
	}

	return &UploaderServiceImpl{
		client:        client,
		clientControl: clientControl,
		projectID:     os.Getenv("GOOGLE_CLOUD_PROJECT"),
		bucketName:    os.Getenv("GOOGLE_CLOUD_BUCKET"),
	}
}

func (uploader *UploaderServiceImpl) Close() {
	if err := uploader.client.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}
	if err := uploader.clientControl.Close(); err != nil {
		log.Printf("Error closing client control: %v", err)
	}
}

func (uploader *UploaderServiceImpl) UploadCharacterImage(ctx context.Context, request *web.ImageCharacterUploadRequest) string {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := uploader.client.Bucket(uploader.bucketName)

	uniqueFileName := fmt.Sprintf("gacha/%d/%d", request.GachaSystemId, request.Id)
	object := bucket.Object(uniqueFileName)

	writerContext := object.NewWriter(ctx)
	_, err := io.Copy(writerContext, request.CharacterImage)
	helper.PanicIfError(err, "failed to upload file to GCS")

	err = writerContext.Close()
	helper.PanicIfError(err, "failed to close GCS writer")

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", uploader.bucketName, uniqueFileName)

}

func (uploader *UploaderServiceImpl) DeleteCharacterImage(ctx context.Context, characterId int, gachaSystemId int) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := uploader.client.Bucket(uploader.bucketName)

	uniqueFileName := fmt.Sprintf("gacha/%d/%d", gachaSystemId, characterId)
	object := bucket.Object(uniqueFileName)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := object.Attrs(ctx)
	if err != nil && (errors.Is(err, storage.ErrBucketNotExist) || errors.Is(err, storage.ErrObjectNotExist)) {
		return
	}
	helper.PanicIfError(err, "failed to get object attributes")

	object = object.If(storage.Conditions{GenerationMatch: attrs.Generation})

	err = object.Delete(ctx)
	if err != nil && (errors.Is(err, storage.ErrBucketNotExist) || errors.Is(err, storage.ErrObjectNotExist)) {
		return
	}
	helper.PanicIfError(err, "failed to delete object")
}

func (uploader *UploaderServiceImpl) DeleteGachaSystemCharacterImage(ctx context.Context, gachaSystemId int) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	gachaSystemFolder := fmt.Sprintf("gacha/%d", gachaSystemId)

	// Construct folder path including the bucket name.
	folderPath := fmt.Sprintf("projects/_/buckets/%v/folders/%v", uploader.bucketName, gachaSystemFolder)

	req := &controlpb.DeleteFolderRequest{
		Name: folderPath,
	}
	err := uploader.clientControl.DeleteFolder(ctx, req)

	if err != nil && (errors.Is(err, storage.ErrBucketNotExist) || errors.Is(err, storage.ErrObjectNotExist)) {
		return
	}
}
