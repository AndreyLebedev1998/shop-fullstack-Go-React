package productstats

import (
	"context"
	"database/sql"
	"fmt"
	"products-microservice/models"
	"strings"
)

func ProductStats(db *sql.DB, stats []models.ProductStats) error {
	if len(stats) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(stats))
	valueArgs := make([]interface{}, 0, len(stats)*2)

	for i, p := range stats {
		n := i * 2
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", n+1, n+2))
		valueArgs = append(valueArgs, p.ProductId, p.PurchaseCount)
	}

	query := fmt.Sprintf(`INSERT INTO product_stats (product_id, purchase_count)
		VALUES %s
		ON CONFLICT (product_id) DO UPDATE
		SET purchase_count = product_stats.purchase_count + EXCLUDED.purchase_count`, strings.Join(valueStrings, ","))

	ctx := context.Background()

	_, err := db.ExecContext(ctx, query, valueArgs...)
	return err
}
