package app

import (
	"errors"
	"fmt"
	"gacha-master/exception"
	"gacha-master/helper"
	"gacha-master/model/web"
	"net/http"
)

type errorData struct {
	Message string `json:"message"`
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ErrorHandler(writer, request, err)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if actualErr, ok := err.(error); ok {
		if notFoundError(writer, request, actualErr) {
			return
		}
		if userError(writer, request, actualErr) {
			return
		}
		if conflictError(writer, request, actualErr) {
			return
		}
		if badRequestError(writer, request, actualErr) {
			return
		}
		internalServerError(writer, request, actualErr)
	} else {
		internalServerError(writer, request, errors.New(fmt.Sprintf("%v", err)))
	}
}

func writeErrorResponse(writer http.ResponseWriter, statusCode int, status string, message string) {
	writer.WriteHeader(statusCode)

	webResponse := web.WebResponse{
		Code:   statusCode,
		Status: status,
		Data:   errorData{message},
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err error) bool {
	var notFoundErr *exception.NotFoundError
	if errors.As(err, &notFoundErr) {
		writeErrorResponse(writer, http.StatusNotFound, "NOT FOUND", notFoundErr.Error())
		return true
	}
	return false
}

func userError(writer http.ResponseWriter, request *http.Request, err error) bool {
	var userErr *exception.UserError
	if errors.As(err, &userErr) {
		writeErrorResponse(writer, http.StatusOK, "OK", userErr.Error())
		return true
	}
	return false
}

func conflictError(writer http.ResponseWriter, request *http.Request, err error) bool {
	var userErr *exception.ConflictError
	if errors.As(err, &userErr) {
		writeErrorResponse(writer, http.StatusConflict, "CONFLICT", userErr.Error())
		return true
	}
	return false
}

func badRequestError(writer http.ResponseWriter, request *http.Request, err error) bool {
	var userErr *exception.BadRequestError
	if errors.As(err, &userErr) {
		writeErrorResponse(writer, http.StatusBadRequest, "BAD REQUEST", userErr.Error())
		return true
	}
	return false
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err error) {
	writeErrorResponse(writer, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err.Error())
}
