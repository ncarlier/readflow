package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/converter"
	"github.com/ncarlier/readflow/pkg/integration/webhook"
	"github.com/ncarlier/readflow/pkg/model"
)

// ProviderConfig is the structure definition of a S3 configuration
type ProviderConfig struct {
	Endpoint       string `json:"endpoint"`
	AccesKeyID     string `json:"acess_key_id"`
	AccesKeySecret string `json:"acess_key_secret"`
	Region         string `json:"region"`
	Bucket         string `json:"bucket"`
	Format         string `json:"format"`
}

// Provider is the structure definition of a S3 outbound service
type Provider struct {
	config    ProviderConfig
	client    *minio.Client
	converter converter.ArticleConverter
}

func newS3Provider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	config := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &config); err != nil {
		return nil, err
	}

	// Get article converter
	conv, err := converter.GetArticleConverter(config.Format)
	if err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Create S3 client
	client, err := minio.New(endpoint.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccesKeyID, config.AccesKeySecret, ""),
		Secure: endpoint.Scheme == "https",
	})
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		config:    config,
		client:    client,
		converter: conv,
	}

	return provider, nil
}

// Send article to Webhook endpoint.
func (s3p *Provider) Send(ctx context.Context, article model.Article) error {
	asset, err := s3p.converter.Convert(ctx, &article)
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
