// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"attendance/api/auth"
	"attendance/api/restHandler"
	"attendance/api/router"
	"attendance/repository"
	"attendance/services"
	"github.com/go-pg/pg"
)

// Injectors from wire.go:

func InitializeApp(db *pg.DB) *router.MUXRouterImpl {
	repositoryImpl := repository.NewRepositoryImpl(db)
	homeServiceImpl := services.NewHomeServiceImpl(repositoryImpl)
	homeImpl := restHandler.NewHomeImpl(homeServiceImpl)
	authTokenImpl := auth.NewAuthTokenImpl(repositoryImpl)
	loginImpl := restHandler.NewLoginImpl(repositoryImpl, authTokenImpl)
	addNewUserServiceImpl := services.NewAddNewUserServiceImpl(repositoryImpl)
	addNewUserImpl := restHandler.NewAddNewUserImpl(addNewUserServiceImpl)
	punchInOutServiceImpl := services.NewPunchInOutServiceImpl(repositoryImpl)
	punchInOutImpl := restHandler.NewPunchInOutImpl(punchInOutServiceImpl)
	getTeacherAttendanceServiceImpl := services.NewGetTeacherAttendanceServiceImpl(repositoryImpl)
	getTeacherAttendanceImpl := restHandler.NewGetTeacherAttendanceImpl(getTeacherAttendanceServiceImpl)
	getClassAttendanceImpl := services.NewGetClassAttendanceImpl(repositoryImpl)
	restHandlerGetClassAttendanceImpl := restHandler.NewGetClassAttendanceImpl(getClassAttendanceImpl)
	getStudentAttendanceServiceImpl := services.NewGetStudentAttendanceServiceImpl(repositoryImpl)
	getStudentAttendanceImpl := restHandler.NewGetStudentAttendanceImpl(getStudentAttendanceServiceImpl)
	fetchStatusImpl := services.NewFetchStatusImpl(repositoryImpl)
	restHandlerFetchStatusImpl := restHandler.NewFetchStatusImpl(fetchStatusImpl)
	muxRouterImpl := router.NewMUXRouterImpl(homeImpl, loginImpl, addNewUserImpl, punchInOutImpl, getTeacherAttendanceImpl, restHandlerGetClassAttendanceImpl, getStudentAttendanceImpl, restHandlerFetchStatusImpl)
	return muxRouterImpl
}