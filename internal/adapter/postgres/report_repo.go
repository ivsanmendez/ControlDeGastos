package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ivsanmendez/ControlDeContabilidad/internal/domain/report"
)

// ReportRepo implements report.Repository.
type ReportRepo struct {
	db *sql.DB
}

func NewReportRepo(db *sql.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) AggregateIncomeByMonth(ctx context.Context, year int) ([]report.MonthAggregate, error) {
	const q = `
		SELECT EXTRACT(MONTH FROM payment_date)::int, COALESCE(SUM(amount), 0)
		FROM contributions
		WHERE EXTRACT(YEAR FROM payment_date)::int = $1
		GROUP BY EXTRACT(MONTH FROM payment_date)
		ORDER BY EXTRACT(MONTH FROM payment_date)`

	return r.scanAggregates(ctx, q, year)
}

func (r *ReportRepo) AggregateExpensesByMonth(ctx context.Context, year int) ([]report.MonthAggregate, error) {
	const q = `
		SELECT EXTRACT(MONTH FROM date)::int, COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE EXTRACT(YEAR FROM date)::int = $1
		GROUP BY EXTRACT(MONTH FROM date)
		ORDER BY EXTRACT(MONTH FROM date)`

	return r.scanAggregates(ctx, q, year)
}

func (r *ReportRepo) scanAggregates(ctx context.Context, query string, args ...any) ([]report.MonthAggregate, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("report aggregate: %w", err)
	}
	defer rows.Close()

	var result []report.MonthAggregate
	for rows.Next() {
		var a report.MonthAggregate
		if err := rows.Scan(&a.Month, &a.Amount); err != nil {
			return nil, fmt.Errorf("scan aggregate: %w", err)
		}
		result = append(result, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("report aggregate: %w", err)
	}
	return result, nil
}
