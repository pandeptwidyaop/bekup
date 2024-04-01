package s3

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

var ses *session.Session

func Run(ctx context.Context, in <-chan models.BackupFileInfo, destination config.ConfigDestination, worker int) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)

	wg := sync.WaitGroup{}

	var lists []<-chan models.BackupFileInfo

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		lists = append(lists, run(ctx, in, destination))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ls := range lists {
		go func(c <-chan models.BackupFileInfo) {
			for cc := range ls {
				out <- cc
			}
			wg.Done()
		}(ls)
	}

	return out
}

func run(ctx context.Context, in <-chan models.BackupFileInfo, destination config.ConfigDestination) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)

	go func() {
		defer close(out)

		for f := range in {
			select {
			case out <- doUpload(ctx, f, destination):
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func doUpload(ctx context.Context, f models.BackupFileInfo, d config.ConfigDestination) models.BackupFileInfo {
	if ses == nil {
		c, err := newSession(d)
		if err != nil {
			f.Err = err
			return f
		}

		ses = c
	}

	up := s3manager.NewUploader(ses)

	file, err := os.Open(f.ZipPath)
	if err != nil {
		f.Err = err
		return f
	}
	defer file.Close()

	input := s3manager.UploadInput{
		Bucket: aws.String(d.AWSBucket),
		Body:   file,
	}

	_, err = up.UploadWithContext(ctx, &input)
	if err != nil {
		f.Err = err
		return f
	}

	return f
}

func newSession(src config.ConfigDestination) (*session.Session, error) {
	if src.Region == "" {
		return nil, exception.ErrAwsRegionNotExist
	}

	if src.AWSAccessKey == "" {
		return nil, exception.ErrAwsAccessKeyNotExist
	}

	if src.AWSSecretKey == "" {
		return nil, exception.ErrAwsAccessKeySecretNotExist
	}

	if src.AWSBucket == "" {
		return nil, exception.ErrAwsBucketNotExist
	}

	return session.NewSession(&aws.Config{
		Region: aws.String(src.Region),
		Credentials: credentials.NewStaticCredentials(
			src.AWSAccessKey,
			src.AWSSecretKey,
			"",
		),
	})
}
