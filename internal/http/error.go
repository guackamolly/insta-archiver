package http

import (
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func onCustomError(
	ectx echo.Context,
	err *model.Error,
) error {
	ectx.Logger().Errorf("processing error: %v\n", err)

	return ectx.File(fallback)
}
