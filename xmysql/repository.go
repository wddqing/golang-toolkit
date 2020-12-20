package xmysql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	stmtDB *StmtDB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{NewStmtDB(db)}
}

func (r *Repository) Exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	stmt, err := r.stmtDB.GetStmt(sqlStr)
	if err != nil && err != ErrCanNotPrepared {
		return nil, err
	}
	if stmt == nil {
		return r.stmtDB.Exec(sqlStr, args)
	}

	return stmt.Exec(args...)
}

func (r *Repository) ExecNamed(sqlStr string, args interface{}) (sql.Result, error) {
	stmt, err := r.stmtDB.GetNamedStmt(sqlStr)
	if err != nil && err != ErrCanNotPrepared {
		return nil, err
	}
	if stmt == nil {
		return r.stmtDB.ExecNamed(sqlStr, args)
	}

	return stmt.Exec(args)
}

func (r *Repository) Select(sqlStr string, obj interface{}, args ...interface{}) error {
	stmt, err := r.stmtDB.GetStmt(sqlStr)
	if err != nil && err != ErrCanNotPrepared {
		return err
	}
	if stmt == nil {
		return r.stmtDB.Select(sqlStr, obj, args)
	}

	return stmt.Select(obj, args...)
}

func (r *Repository) Get(sqlStr string, obj interface{}, args ...interface{}) error {
	stmt, err := r.stmtDB.GetStmt(sqlStr)
	if err != nil && err != ErrCanNotPrepared {
		return err
	}
	if stmt == nil {
		return r.stmtDB.Get(sqlStr, obj, args)
	}

	return stmt.Get(obj, args...)
}
