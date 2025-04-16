package admin

import (
	"net/http"

	"github.com/admin-golang/admin"

	"github.com/design-doc-self-generator/golang-scaffolding/domain/admin/dashboard"
	"github.com/design-doc-self-generator/golang-scaffolding/domain/admin/layout"
	"github.com/design-doc-self-generator/golang-scaffolding/domain/admin/signin"
	"github.com/design-doc-self-generator/golang-scaffolding/domain/admin/user"
)

type handlers struct {
	adminHandler http.Handler
}

func (h *handlers) GetAdmin() http.Handler {
	return h.adminHandler
}

func NewHandlers() (*handlers, error) {
	signInFormPage, err := signin.NewFormPage()
	if err != nil {
		return nil, err
	}

	pages := admin.Pages{
		signInFormPage,
		dashboard.NewPage(),
		user.NewList(),
	}

	adminHandler := admin.New(&admin.Config{
		DebugMode: false,
		UITheme:   admin.MaterialUI,
		Pages:     pages,
		Layout:    layout.NewLayout(),
	})

	return &handlers{adminHandler: adminHandler}, nil
}
