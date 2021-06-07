package codesnippetrepo

import (
	"database/sql"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"github.com/mp-hl-2021/code-swamp/internal/interface/memory/codesnippetrepo"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const queryCreateSnippet = `
	INSERT INTO snippets(
		code,
		language,                 
		lifetime,
	    isChecked,
	    message
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING id
`

const querySetCodeLinterMessage = `
	UPDATE snippets
	SET isChecked = $2,
	    message = $3
	WHERE snippets.id = $1
`

func (p *Postgres) SetCodeLinterMessage(sid uint, msg string) error {
	_, err := p.conn.Query(querySetCodeLinterMessage, sid, true, msg)
	return err
}

func (p *Postgres) CreateCodeSnippet(s codesnippet.CodeSnippet) (uint, error) {
	row := p.conn.QueryRow(queryCreateSnippet, s.Code, s.Lang, s.Lifetime, s.IsChecked, s.Message)
	var id uint
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

const queryCreateSnippetWithUser = `
	INSERT INTO snippets(
		code,
		uid,
		language,
		lifetime,
		isChecked,
	    message
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
`

func (p *Postgres) CreateCodeSnippetWithUser(s codesnippet.CodeSnippet, uid uint) (uint, error) {
	row := p.conn.QueryRow(queryCreateSnippetWithUser, s.Code, uid, s.Lang, s.Lifetime, s.IsChecked, s.Message)
	var id uint
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

const queryGetCodeSnippetById = `
	SELECT
		code,
		language,
	    isChecked,
	    message
	FROM snippets
	WHERE id = $1
`

func (p *Postgres) GetCodeSnippetById(sid uint) (codesnippet.CodeSnippet, error) {
	cs := codesnippet.CodeSnippet{}
	row := p.conn.QueryRow(queryGetCodeSnippetById, sid)
	err := row.Scan(&cs.Code, &cs.Lang, &cs.IsChecked, &cs.Message)
	if err != nil {
		if err == sql.ErrNoRows {
			return codesnippet.CodeSnippet{}, codesnippetrepo.ErrInvalidSnippedId
		}
		return codesnippet.CodeSnippet{}, err
	}
	return cs, nil
}

const queryGetMyCodeSnippetIds = `
	SELECT id
	FROM snippets
	WHERE uid = $1
`

func (p *Postgres) GetMyCodeSnippetIds(uid uint) ([]uint, error) {
	var ids []uint
	row, err := p.conn.Query(queryGetMyCodeSnippetIds, uid)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var id uint
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

const queryDeleteExpiredSnippets = `
	DELETE FROM snippets
	WHERE createdAt < now() - lifetime
`

func (p *Postgres) DeleteExpiredSnippets() error {
	_, err := p.conn.Query(queryDeleteExpiredSnippets)
	return err
}
