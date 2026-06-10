package main

import (
	"github.com/iv-tunate/bizpapertrail/internal/utils"
	"github.com/labstack/echo"
)

func checkserverstatus(ctx echo.Context) error{
	return utils.SuccessResponse(ctx, 200, "Server running", "ok", nil)
}