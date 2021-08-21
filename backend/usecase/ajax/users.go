package ajax

import "context"

// IsValidUsername check if username either useable or not
func (u *Ajax) IsValidUsername(ctx context.Context, newUsername, oldUsername string) (ok bool, err error) {
	if newUsername == oldUsername {
		ok = true
		return
	}
	ok, err = u.database.IsValidUsername(ctx, newUsername)
	return
}
