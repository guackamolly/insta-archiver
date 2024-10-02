package http

import (
	"github.com/guackamolly/insta-archiver/internal/model"
	"github.com/labstack/echo/v4"
)

func onCustomError(
	ectx echo.Context,
	err *model.Error,
) error {
	return ectx.File(fallback)
}
