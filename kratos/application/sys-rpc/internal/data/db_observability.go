package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
	"time"

	appmetrics "github.com/gcc798/microservice-kit/kratos/pkg/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type dbObservability struct {
	slowThreshold time.Duration
}

func openObservedPostgresDB(dsn string, obs dbObservability) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	connector := stdlib.GetConnector(*connConfig)
	return sql.OpenDB(observedConnector{base: connector, obs: obs}), nil
}

type observedConnector struct {
	base driver.Connector
	obs  dbObservability
}

func (c observedConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return observedConn{Conn: conn, obs: c.obs}, nil
}

func (c observedConnector) Driver() driver.Driver {
	return c.base.Driver()
}

type observedConn struct {
	driver.Conn
	obs dbObservability
}

func (c observedConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return observedStmt{Stmt: stmt, query: query, obs: c.obs}, nil
}

func (c observedConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if preparer, ok := c.Conn.(driver.ConnPrepareContext); ok {
		stmt, err := preparer.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return observedStmt{Stmt: stmt, query: query, obs: c.obs}, nil
	}
	return c.Prepare(query)
}

func (c observedConn) Begin() (driver.Tx, error) {
	tx, err := c.Conn.Begin()
	if err != nil {
		return nil, err
	}
	return observedTx{Tx: tx}, nil
}

func (c observedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if beginner, ok := c.Conn.(driver.ConnBeginTx); ok {
		tx, err := beginner.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		return observedTx{Tx: tx}, nil
	}
	return c.Begin()
}

func (c observedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	start := time.Now()
	result, err := execer.ExecContext(ctx, query, args)
	c.obs.observeDB(query, argsToAny(args), start, err)
	return result, err
}

func (c observedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	start := time.Now()
	rows, err := queryer.QueryContext(ctx, query, args)
	c.obs.observeDB(query, argsToAny(args), start, err)
	return rows, err
}

func (c observedConn) CheckNamedValue(value *driver.NamedValue) error {
	if checker, ok := c.Conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(value)
	}
	return nil
}

func (c observedConn) Ping(ctx context.Context) error {
	if pinger, ok := c.Conn.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}
	return nil
}

func (c observedConn) ResetSession(ctx context.Context) error {
	if resetter, ok := c.Conn.(driver.SessionResetter); ok {
		return resetter.ResetSession(ctx)
	}
	return nil
}

func (c observedConn) IsValid() bool {
	if validator, ok := c.Conn.(driver.Validator); ok {
		return validator.IsValid()
	}
	return true
}

type observedStmt struct {
	driver.Stmt
	query string
	obs   dbObservability
}

func (s observedStmt) Exec(args []driver.Value) (driver.Result, error) {
	start := time.Now()
	result, err := s.Stmt.Exec(args)
	s.obs.observeDB(s.query, valuesToAny(args), start, err)
	return result, err
}

func (s observedStmt) Query(args []driver.Value) (driver.Rows, error) {
	start := time.Now()
	rows, err := s.Stmt.Query(args)
	s.obs.observeDB(s.query, valuesToAny(args), start, err)
	return rows, err
}

func (s observedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if execer, ok := s.Stmt.(driver.StmtExecContext); ok {
		start := time.Now()
		result, err := execer.ExecContext(ctx, args)
		s.obs.observeDB(s.query, argsToAny(args), start, err)
		return result, err
	}
	return nil, driver.ErrSkip
}

func (s observedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if queryer, ok := s.Stmt.(driver.StmtQueryContext); ok {
		start := time.Now()
		rows, err := queryer.QueryContext(ctx, args)
		s.obs.observeDB(s.query, argsToAny(args), start, err)
		return rows, err
	}
	return nil, driver.ErrSkip
}

func (s observedStmt) ColumnConverter(idx int) driver.ValueConverter {
	if converter, ok := s.Stmt.(driver.ColumnConverter); ok {
		return converter.ColumnConverter(idx)
	}
	return driver.DefaultParameterConverter
}

type observedTx struct {
	driver.Tx
}

func (o dbObservability) observeDB(query string, args []any, start time.Time, err error) {
	duration := time.Since(start)
	operation := sqlOperation(query)
	status := observeStatus(err)
	appmetrics.ObserveDB(operation, status, duration)
	if o.slowThreshold > 0 && duration >= o.slowThreshold {
		log.Printf("level=WARN msg=%q operation=%s duration_ms=%d status=%s sql=%q args=%v err=%v",
			"slow db operation",
			operation,
			duration.Milliseconds(),
			status,
			normalizeSQL(query),
			args,
			err,
		)
	}
}

func sqlOperation(query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return "unknown"
	}
	parts := strings.Fields(query)
	if len(parts) == 0 {
		return "unknown"
	}
	return strings.ToLower(parts[0])
}

func normalizeSQL(query string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(query)), " ")
}

func argsToAny(values []driver.NamedValue) []any {
	items := make([]any, 0, len(values))
	for _, item := range values {
		items = append(items, item.Value)
	}
	return items
}

func valuesToAny(values []driver.Value) []any {
	items := make([]any, 0, len(values))
	for _, item := range values {
		items = append(items, item)
	}
	return items
}

func observeStatus(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}

func startDBPoolMetricsSampler(db *sql.DB, interval time.Duration, stop <-chan struct{}) {
	if db == nil || interval <= 0 {
		return
	}
	appmetrics.SetDBPoolStats(db.Stats())
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				appmetrics.SetDBPoolStats(db.Stats())
			case <-stop:
				return
			}
		}
	}()
}

func formatObservedQuery(query string, args []any) string {
	return fmt.Sprintf("sql=%q args=%v", normalizeSQL(query), args)
}
