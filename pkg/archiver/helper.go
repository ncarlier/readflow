package archiver

import (
	"encoding/base64"
	"fmt"
)

// createDataURL returns base64 encoded data URL
func createDataURL(content []byte, contentType string) string {
	b64encoded := base64.StdEncoding.EncodeToString(content)
	return fmt.Sprintf("data:%s;base64,%s", contentType, b64encoded)
}
