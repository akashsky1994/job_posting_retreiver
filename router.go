package main

import (
	"akashsky1994/job_retreiver/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BuiltInRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		// Seek, verify and validate JWT tokens
		// router.Use(jwtauth.Verifier(tokenAuth))
		// router.Use(jwtauth.Authenticator)
		// router.Use(newRelicMiddleWare)
		// router.Use(pageLimitMiddleWare)

		router.Get("/jobs/{category_id}", func(res http.ResponseWriter, req *http.Request) {
			category_id := chi.URLParam(req, "category_id")
			err := handler.FetchBuiltInJobs(category_id)
			if err != nil {
				panic(err)
			}
			message := map[string]string{"message": "Fetching Successful"}
			RespondwithJSON(res, http.StatusOK, message)

			// http.Redirect(res, req, filepath, http.StatusOK)
			// res.WriteHeader(http.StatusOK)
			// res.Header().Set("Content-Type", "application/octet-stream")
			// res.Write(fileBytes)
			// return
		})

	})
	return router
}
