package all

import (
	// activate generic outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/generic"
	// activate keeper outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/keeper"
	// activate pocket outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/pocket"
	// activate readflow outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/readflow"
	// activate s3 outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/s3"
	// activate shaarli outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/shaarli"
	// activate wallabag outgoing webhook
	_ "github.com/ncarlier/readflow/internal/integration/webhook/wallabag"
)
