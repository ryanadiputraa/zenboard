package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/zenboard/internal/domain"
)

type boardRepository struct {
	db *sqlx.DB
}

func NewBoardRepository(db *sqlx.DB) domain.BoardRepository {
	return &boardRepository{
		db: db,
	}
}

func (r *boardRepository) Init(ctx context.Context, initBoard domain.InitBoardDTO, task1, task2, task3 domain.InitTaskDTO,
) (board domain.Board, err error) {
	tx := r.db.MustBeginTx(ctx, &sql.TxOptions{})

	// check if user board already exists
	var existingBoard domain.Board
	err = tx.Get(&existingBoard, "SELECT id FROM boards WHERE owner_id = $1 LIMIT 1", initBoard.OwnerID)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return
	}
	if existingBoard.ID != "" {
		tx.Rollback()
		return
	}

	// init user board
	if err = tx.QueryRowx(`INSERT INTO boards ( id, project_name, picture, owner_id, created_at
		) VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		initBoard.ID, initBoard.ProjectName, "", initBoard.OwnerID, initBoard.CreatedAt,
	).StructScan(&board); err != nil {
		tx.Rollback()
		return
	}

	// init default board tasks
	if _, err = tx.Exec(`INSERT INTO task_status ( id, "order", name, board_id
		) VALUES
		($2, $3, $4, $1),
		($5, $6, $7, $1),
		($8, $9, $10, $1)
		RETURNING *`,
		initBoard.ID,
		task1.ID, task1.Order, task1.Name,
		task2.ID, task2.Order, task2.Name,
		task3.ID, task3.Order, task3.Name,
	); err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}

func (r *boardRepository) FetchByOwnerID(ctx context.Context, ownerID string) (boards []domain.Board, err error) {
	err = r.db.Select(&boards, `SELECT * FROM boards WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	return
}
