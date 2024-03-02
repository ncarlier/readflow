package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ncarlier/readflow/internal/config"
	"github.com/ncarlier/readflow/internal/exporter"
	"github.com/ncarlier/readflow/internal/global"
	"github.com/ncarlier/readflow/internal/integration/webhook"
	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/downloader"
)

// ProviderConfig is the structure definition of a S3 configuration
type ProviderConfig struct {
	Endpoint    string `json:"endpoint"`
	AccessKeyID string `json:"access_key_id"`
	Region      string `json:"region"`
	Bucket      string `json:"bucket"`
	Format      string `json:"format"`
}

// Provider is the structure definition of a S3 outbound service
type Provider struct {
	config          ProviderConfig
	client          *minio.Client
	AccessKeySecret string
}

func newS3Provider(srv model.OutgoingWebhook, conf config.Config) (webhook.Provider, error) {
	cfg := ProviderConfig{}
	if err := json.Unmarshal([]byte(srv.Config), &cfg); err != nil {
		return nil, err
	}

	// Validate endpoint URL
	endpoint, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	// Validate secrets
	accessKeySecret, ok := srv.Secrets["access_key_secret"]
	if !ok {
		return nil, fmt.Errorf("missing access key secret")
	}

	// Create S3 client
	client, err := minio.New(endpoint.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, accessKeySecret, ""),
		Secure: endpoint.Scheme == "https",
	})
	if err != nil {
		return nil, err
	}

	provider := &Provider{
		config: cfg,
		client: client,
	}

	return provider, nil
}

// Send article to Webhook endpoint.
func (s3p *Provider) Send(ctx context.Context, article model.Article) (*webhook.Result, error) {
	// Get download from context
	// /!\ this is a ugly hack required to simplify service coupling
	ctxValue := ctx.Value(global.ContextDownloader)
	if ctxValue == nil {
		return nil, errors.New("downloader not found inside the context")
	}
	dl := ctxValue.(downloader.Downloader)

	// Get article exporter
	exp, err := exporter.NewArticleExporter(s3p.config.Format, dl)
	if err != nil {
		return nil, err
	}
	asset, err := exp.Export(ctx, &article)
	if err != nil {
		return nil, err
	}
	data := bytes.NewReader(asset.Data)

	_, err = s3p.client.PutObject(ctx, s3p.config.Bucket, asset.Name, data, int64(len(asset.Data)), minio.PutObjectOptions{ContentType: asset.ContentType})
	if err != nil {
		return nil, err
	}

	return &webhook.Result{}, nil
}

func init() {
	webhook.Register("s3", &webhook.Def{
		Name:   "S3",
		Desc:   "Export article(s) to a S3 bucket.",
		Create: newS3Provider,
	})
}
