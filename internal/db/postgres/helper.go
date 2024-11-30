package postgres

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ncarlier/readflow/pkg/utils"
)

// hashtagsExpr create ARRAY SQL expression from an array of hastags
// Replacing the #prefix by ~prefix in order to avoid lexeme convertion
func hashtagsExpr(values []string) interface{} {
	if len(values) == 0 {
		return "{}"
	}
	args := make([]interface{}, len(values))
	for i, value := range values {
		args[i] = utils.ReplaceHashtagsPrefix(value, "~")
	}
	expr := strings.Join(strings.Split(strings.Repeat("?", len(values)), ""), ",")
	expr = fmt.Sprintf("ARRAY[%s]", expr)
	return sq.Expr(expr, args...)
}
