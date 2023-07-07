package repository

import (
	"context"
	"database/sql"

	"github.com/ryanadiputraa/zenboard/domain"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"
)

type boardRepository struct {
	db *db.Queries
	tx db.Tx
}

func NewBoardRepository(db *db.Queries, tx db.Tx) domain.BoardRepository {
	return &boardRepository{
		db: db,
		tx: tx,
	}
}

func (r *boardRepository) Init(
	ctx context.Context,
	initBoard domain.InitBoardDTO,
	task1, task2, task3 domain.InitTaskDTO,
) (board domain.Board, err error) {
	err = r.tx.ExecTx(ctx, func(q *db.Queries) (err error) {
		user, err := q.GetBoardByOwnerID(ctx, initBoard.OwnerID)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if user.ID != "" {
			return
		}

		arg1 := db.CreateBoardParams{
			ID:          initBoard.ID,
			ProjectName: initBoard.ProjectName,
			OwnerID:     initBoard.OwnerID,
			CreatedAt:   initBoard.CreatedAt,
		}

		created, err := q.CreateBoard(ctx, arg1)
		if err != nil {
			return err
		}

		arg2 := db.InitTaskStatusParams{
			BoardID: created.ID,
			ID:      task1.ID,
			Order:   int32(task1.Order),
			Name:    task1.Name,
			ID_2:    task2.ID,
			Order_2: int32(task2.Order),
			Name_2:  task2.Name,
			ID_3:    task3.ID,
			Order_3: int32(task3.Order),
			Name_3:  task3.Name,
		}

		_, err = q.InitTaskStatus(ctx, arg2)
		if err != nil {
			return
		}

		board = domain.Board{
			ID:          created.ID,
			ProjectName: created.ProjectName,
			OwnerID:     created.OwnerID,
			CreatedAt:   created.CreatedAt,
		}

		return
	})

	return
}
