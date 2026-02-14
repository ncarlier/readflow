package all

import (
	// activate all ctl commands
	_ "github.com/ncarlier/readflow/ctl/cmd/articles"
	_ "github.com/ncarlier/readflow/ctl/cmd/categories"
	_ "github.com/ncarlier/readflow/ctl/cmd/incoming-webhooks"
	_ "github.com/ncarlier/readflow/ctl/cmd/init-config"
	_ "github.com/ncarlier/readflow/ctl/cmd/me"
	_ "github.com/ncarlier/readflow/ctl/cmd/outgoing-webhooks"
	_ "github.com/ncarlier/readflow/ctl/cmd/version"
)
