package users

import "context"

func getCtxWithUser(ctx context.Context, createdBy *string) context.Context {
	var newCtx context.Context
	if createdBy != nil {
		newCtx = context.WithValue(ctx, "user_id", *createdBy)
	} else {
		newCtx = ctx
	}

	return newCtx
}
