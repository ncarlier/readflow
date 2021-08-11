package all

import (
	// activate generic outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/generic"
	// activate keeper outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/keeper"
	// activate pocket outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/pocket"
	// activate s3 outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/s3"
	// activate wallabag outgoing webhook
	_ "github.com/ncarlier/readflow/pkg/integration/webhook/wallabag"
)
