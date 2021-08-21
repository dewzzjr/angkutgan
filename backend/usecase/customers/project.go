package customers

import (
	"context"
	"sort"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// UpdateProject update project by customer code
func (i *Customers) UpdateProject(ctx context.Context, code string, old, new []model.Project, actionBy int64) (err error) {
	// LOAD INSERT WITH ALL NEW
	insert := new
	// LOAD DELETE WITH EMPTY SLICE
	delete := make([]int64, 0)
	// COMPARE OLD AND NEW, THEN SAVE INDEX
	same := make([]int, 0)
	for _, o := range old {
		var match bool
		for j, n := range new {
			if o.Name == n.Name &&
				o.Location == n.Location {
				same = append(same, j)
				match = true
				break
			}
		}
		// DELETE ID WHEN OLD DONT MATCH FROM NEW VALUE
		if !match {
			delete = append(delete, o.ID)
		}
	}
	// SORT, REVERSE LOOP, REMOVE FROM BIGGEST INDEX
	sort.Ints(same)
	for i := range same {
		n := same[len(same)-1-i]
		insert[len(insert)-1], insert[n] = insert[n], insert[len(insert)-1]
		insert = insert[:len(insert)-1]
	}
	if err = i.database.InsertDeleteProject(ctx, code, insert, delete, actionBy); err != nil {
		err = errors.Wrapf(err, "InsertDeleteProject [%s]", code)
	}
	return
}
