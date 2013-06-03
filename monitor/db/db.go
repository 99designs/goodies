// Decorates any database driver with a logger
package db

import "database/sql"
import "github.com/99designs/goodies/monitor"
import "database/sql/driver"
import "time"

type MonitorResult struct {
	Description string
	Duration    time.Duration
}

type DatabaseMonitor interface {
	OnEvent(MonitorResult)
}

func RegisterNewDriver(newname string, driver driver.Driver, mf DatabaseMonitor) {
	sql.Register(newname, MonitoredDriver{delegate: driver, monitorFunc: mf})
}

type MonitoredDriver struct {
	delegate    driver.Driver
	monitorFunc DatabaseMonitor
}

func (md MonitoredDriver) Open(name string) (driver.Conn, error) {
	conn, err := md.delegate.Open(name)
	return MonitoredConn{conn, md.monitorFunc}, err
}

type MonitoredConn struct {
	delegate    driver.Conn
	monitorFunc DatabaseMonitor
}

func (c MonitoredConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.delegate.Prepare(query)
	return MonitoredStmt{stmt, query, c.monitorFunc}, err
}

func (c MonitoredConn) Close() error {
	return c.delegate.Close()
}

func (c MonitoredConn) Begin() (driver.Tx, error) {
	return c.delegate.Begin()
}

type MonitoredTx struct {
	delegate    driver.Tx
	monitorFunc DatabaseMonitor
}

func (c MonitoredTx) Commit() error {
	return c.delegate.Commit()
}

func (c MonitoredTx) Rollback() error {
	return c.delegate.Rollback()
}

type MonitoredStmt struct {
	delegate    driver.Stmt
	query       string
	monitorFunc DatabaseMonitor
}

func (c MonitoredStmt) Close() error {
	return c.delegate.Close()
}

func (c MonitoredStmt) NumInput() int {
	return c.delegate.NumInput()
}

func (c MonitoredStmt) Exec(args []driver.Value) (driver.Result, error) {
	var result driver.Result
	var err error
	duration := monitor.Duration(func() {
		result, err = c.delegate.Exec(args)
	})

	c.monitorFunc.OnEvent(MonitorResult{
		Description: c.query,
		Duration:    duration,
	})

	return result, err
}

func (c MonitoredStmt) Query(args []driver.Value) (driver.Rows, error) {
	var rows driver.Rows
	var err error
	duration := monitor.Duration(func() {
		rows, err = c.delegate.Query(args)
	})

	c.monitorFunc.OnEvent(MonitorResult{
		Description: c.query,
		Duration:    duration,
	})

	return rows, err
}
