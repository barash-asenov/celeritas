package main

import "github.com/barash-asenov/celeritas"

type application struct {
	App *celeritas.Celeritas
}

func main() {
	c := initApplication()

	c.App.ListenAndServe()
}
