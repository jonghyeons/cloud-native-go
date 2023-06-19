package ch04

import "context"

type Circuit func(context.Context) (string, error)
