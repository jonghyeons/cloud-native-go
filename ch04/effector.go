package ch04

import "context"

type Effector func(ctx context.Context) (string, error)
