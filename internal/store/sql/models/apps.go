package models

import "github.com/vishenosik/sso/internal/entities"

type App struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	Secret string `db:"secret"`
}

func AppToEntities(app *App) *entities.App {
	if app == nil {
		return new(entities.App)
	}
	return &entities.App{
		ID:     app.ID,
		Name:   app.Name,
		Secret: app.Secret,
	}
}

func AppFromEntities(app *entities.App) *App {
	if app == nil {
		return new(App)
	}
	return &App{
		ID:     app.ID,
		Name:   app.Name,
		Secret: app.Secret,
	}
}
