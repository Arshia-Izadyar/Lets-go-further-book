package handlers

import (
	"clean_api/src/api/dto"
	"clean_api/src/api/filters"
	"clean_api/src/api/helpers"
	"clean_api/src/api/validators"
	"clean_api/src/services"
	"net/http"
)

type MovieHandler struct {
	service *services.MovieService
}

func NewMovieHandler() *MovieHandler {
	s := services.NewMovieService()
	return &MovieHandler{
		service: s,
	}
}

func (mh *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	req := &dto.CreateMovie{}
	err := helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)

		return
	}
	result, err := mh.service.Create(req)

	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	response := helpers.GenerateResponse(result, true, http.StatusCreated)
	err = helpers.WriteResponse(w, response)
	if err != nil {
		helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		return
	}
}

func (mh *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	req := &dto.UpdateMovie{}
	id, err := helpers.ReadParams(r)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}

	err = helpers.ReadRequestBody(w, r, &req)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}

	movie, err := mh.service.GetById(int32(id))
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	if req.Title == ""{
		req.Title = movie.Title
	}
	if req.Year == 0{
		req.Year = movie.Year
	}
	if req.Genres == nil {
		req.Genres = movie.Genres
	}
	if req.Runtime == nil {
		req.Runtime = &movie.Runtime
	}
	res, err := mh.service.Update(req, int32(id), int32(movie.Version))
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(res, true, http.StatusOK))
}

func (mh *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadParams(r)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	err = mh.service.Delete(id)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(nil, true, http.StatusNoContent))
}


func (mh *MovieHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadParams(r)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	res, err := mh.service.GetById(int32(id))
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}
	helpers.WriteResponse(w, helpers.GenerateResponse(res, true, http.StatusOK))

}

func (mh *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	
	filters := &filters.Filter{}

	qs := r.URL.Query()
	
	v := validators.NewValidator()
	title := helpers.ReadString(&qs, "title", "")
	genres := helpers.ReadCSV(&qs, "genres", []string{})
	filters.Page = helpers.ReadInt(&qs, "page", 1, v)
	filters.PageSize = helpers.ReadInt(&qs, "page_size", 5, v)
	filters.Sort = helpers.ReadString(&qs, "sort", "id")
	filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	res, mData,err := mh.service.GetAll(title, genres, *filters)
	if err != nil {
		response := helpers.GenerateResponseWithError(nil, false, http.StatusBadRequest, err)
		helpers.WriteResponse(w, response)
		return
	}

	helpers.WriteResponse(w, helpers.GenerateResponse(map[string]interface{}{
		"result":res,
		"info":mData,
	}, true, http.StatusOK))
}