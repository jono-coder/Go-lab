package audit

import "time"

type Auditable struct {
	CreatedAt *time.Time `db:"created_at"`
	CreatedBy *uint      `db:"created_by"`
	UpdatedAt *time.Time `db:"updated_at"`
	UpdatedBy *uint      `db:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
