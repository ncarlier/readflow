package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/exporter"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	"github.com/ncarlier/readflow/pkg/model"
)

// ProviderConfig is the structure definition of a S3 configuration
type ProviderConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	Format          string `json:"format"`
}

// Provider is the structure definition of a S3 outbound service
type Provider struct {
	config ProviderConfig
	client *minio.Client
}

func newS3Provider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Create S3 client
	client, err := minio.New(endpoint.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.AccessKeySecret, ""),
		Secure: endpoint.Scheme == "https",
	})
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		config: config,
		client: client,
	}

	return provider, nil
}

// Send article to Webhook endpoint.
func (s3p *Provider) Send(ctx context.Context, article model.Article) error {
	// Get download from context
	// /!\ this is a ugly hack required to simplify service coupling
	ctxValue := ctx.Value(constant.ContextDownloader)
	if ctxValue == nil {
		return errors.New("downloader not found inside the context")
	}
	downloader := ctxValue.(exporter.Downloader)

	// Get article exporter
	exp, err := exporter.NewArticleExporter(s3p.config.Format, downloader)
	if err != nil {
		return err
	}
	asset, err := exp.Export(ctx, &article)
	if err != nil {
		return err
	}
	data := bytes.NewReader(asset.Data)

	_, err = s3p.client.PutObject(ctx, s3p.config.Bucket, asset.Name, data, int64(len(asset.Data)), minio.PutObjectOptions{ContentType: asset.ContentType})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	webhook.Register("s3", &webhook.Def{
		Name:   "S3",
		Desc:   "Export article(s) to a S3 bucket.",
		Create: newS3Provider,
	})
}
