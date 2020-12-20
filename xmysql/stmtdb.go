package xmysql

import (
	"database/sql"
	"errors"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type StmtDB struct {
	db *sqlx.DB

	mu         sync.RWMutex
	stmts      map[string]*sqlx.Stmt
	namedStmts map[string]*sqlx.NamedStmt
}

func NewStmtDB(db *sqlx.DB) *StmtDB {
	return &StmtDB{
		db:         db,
		stmts:      make(map[string]*sqlx.Stmt),
		namedStmts: make(map[string]*sqlx.NamedStmt),
	}
}

func (r *StmtDB) DB() *sqlx.DB {
	return r.db
}

func (r *StmtDB) Exec(sqlStr string, args []interface{}) (sql.Result, error) {
	return r.db.Exec(sqlStr, args...)
}

func (r *StmtDB) Select(sqlStr string, obj interface{}, args []interface{}) error {
	return r.db.Select(obj, sqlStr, args...)
}

func (r *StmtDB) Get(sqlStr string, obj interface{}, args []interface{}) error {
	return r.db.Get(obj, sqlStr, args...)
}

func (r *StmtDB) ExecNamed(sqlStr string, args interface{}) (sql.Result, error) {
	return r.db.NamedExec(sqlStr, args)
}

func (r *StmtDB) GetStmt(sqlPrepared string) (*sqlx.Stmt, error) {
	r.mu.RLock()
	stmt, ok := r.stmts[sqlPrepared]
	r.mu.RUnlock()

	if !ok {
		var err error
		stmt, err = r.prepare(sqlPrepared)
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}

func (r *StmtDB) GetNamedStmt(sqlPrepared string) (*sqlx.NamedStmt, error) {
	r.mu.RLock()
	stmt, ok := r.namedStmts[sqlPrepared]
	r.mu.RUnlock()
	if !ok {
		var err error
		stmt, err = r.prepareNamed(sqlPrepared)
		if err != nil {
			return nil, err
		}
	}
	return stmt, nil
}

var (
	ErrCanNotPrepared = errors.New("can't prepare in condition")
)

func (r *StmtDB) prepare(sqlPrepared string) (*sqlx.Stmt, error) {
	if strings.Contains(sqlPrepared, "IN") {
		return nil, ErrCanNotPrepared
	}

	stmt, err := r.db.Preparex(sqlPrepared)
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	if newStmt, ok := r.stmts[sqlPrepared]; ok {
		stmt.Close()
		stmt = newStmt
	} else {
		r.stmts[sqlPrepared] = stmt
	}
	r.mu.Unlock()
	return stmt, nil
}

func (r *StmtDB) prepareNamed(sqlPrepared string) (*sqlx.NamedStmt, error) {
	stmt, err := r.db.PrepareNamed(sqlPrepared)
	if err != nil {
		return nil, err
	}
	r.mu.Lock()
	if newStmt, ok := r.namedStmts[sqlPrepared]; ok {
		stmt.Close()
		stmt = newStmt
	} else {
		r.namedStmts[sqlPrepared] = stmt
	}
	r.mu.Unlock()

	return stmt, nil
}
