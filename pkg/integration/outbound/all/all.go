package all

import (
	// activate keeper outbound service
	_ "github.com/ncarlier/readflow/pkg/integration/outbound/keeper"
	// activate wallabag outbound service
	_ "github.com/ncarlier/readflow/pkg/integration/outbound/wallabag"
	// activate webhook outbound service
	_ "github.com/ncarlier/readflow/pkg/integration/outbound/webhook"
)
