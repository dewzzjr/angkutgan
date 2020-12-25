package items

import (
	"context"
	"sort"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/pkg/errors"
)

// UpdatePriceRent update price rent by item code
func (i *Items) UpdatePriceRent(ctx context.Context, code string, old, new []model.PriceRent) (err error) {
	// LOAD INSERT WITH ALL NEW
	insert := new
	// LOAD DELETE WITH EMPTY SLICE
	delete := make([]int64, 0)
	// COMPARE OLD AND NEW, THEN SAVE INDEX
	same := make([]int, 0)
	for _, o := range old {
		var match bool
		for j, n := range new {
			if o.Description == n.Description &&
				o.Duration == n.Duration &&
				o.TimeUnit == n.TimeUnit &&
				o.Value == n.Value {
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
	if err = i.database.InsertDeleteRentPrice(ctx, code, insert, delete); err != nil {
		err = errors.Wrapf(err, "InsertDeleteRentPrice [%s]", code)
	}
	return
}
