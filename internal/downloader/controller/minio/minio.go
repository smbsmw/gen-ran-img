package minio

import (
	"bufio"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"time"
	"youtube-downloader-rest/internal/model"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient(url, user, password string, ssl bool) *MinioClient {
	client, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: ssl,
	})

	if err != nil {
		log.Fatalln(err)
	}

	if ok, _ := client.BucketExists(context.Background(), "images"); !ok {
		err = client.MakeBucket(context.Background(), "images", minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal("cant make bucket:", err)
		}
	}

	return &MinioClient{client: client}
}

func (m *MinioClient) UploadFile(ctx context.Context, task *model.Task) (string, error) {
	fileName := fmt.Sprintf("%s-%d.jpg", task.ImageName, time.Now().UnixNano())

	resp, err := http.Get("https://picsum.photos/200")
	r := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	_, err = m.client.PutObject(
		ctx,
		"images",
		fileName,
		r,
		-1,
		minio.PutObjectOptions{ContentType: "image/jpg"},
	)

	if err != nil {
		log.Fatalln(err)
	}

	return fileName, nil
}
