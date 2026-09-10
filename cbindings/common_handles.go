package main

/*
#include "headers/common.h"
*/
import "C"
import (
	"context"
	"fmt"

	"go.chrastecky.dev/fio-api/fio"
)

func getCommonHandles(client C.ClientHandle, ctx C.ContextHandle) (fio.Client, context.Context, error) {
	clientGo, err := getHandleObj[fio.Client](handle(client))
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting client object out of handle: %w", err)
	}

	ctxGo, err := getHandleObj[*contextHandle](handle(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting context object out of handle: %w", err)
	}

	return clientGo, ctxGo.ctx, nil
}
