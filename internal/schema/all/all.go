package all

import (
	// activate article schema
	_ "github.com/ncarlier/readflow/internal/schema/article"
	// activate category schema
	_ "github.com/ncarlier/readflow/internal/schema/category"
	// activate device schema
	_ "github.com/ncarlier/readflow/internal/schema/device"
	// activate incoming webhook schema
	_ "github.com/ncarlier/readflow/internal/schema/incoming-webhook"
	// activate outgoing webhook schema
	_ "github.com/ncarlier/readflow/internal/schema/outgoing-webhook"
	// activate plan schema
	_ "github.com/ncarlier/readflow/internal/schema/plan"
	// activate user schema
	_ "github.com/ncarlier/readflow/internal/schema/user"
)
