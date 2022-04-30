package all

import (
	// activate HTML exporter
	_ "github.com/ncarlier/readflow/pkg/exporter/html"
	// activate EPUB exporter
	_ "github.com/ncarlier/readflow/pkg/exporter/epub"
	// activate txtpaper exporter
	_ "github.com/ncarlier/readflow/pkg/exporter/txtpaper"
)
