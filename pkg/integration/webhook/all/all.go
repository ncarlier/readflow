package all

import (
	// activate generic outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/generic"
	// activate keeper outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/keeper"
	// activate wallabag outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/wallabag"
)
