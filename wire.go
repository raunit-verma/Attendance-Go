//go:build wireinject
// +build wireinject

package main

import (
	"attendance/api/auth"
	rh "attendance/api/restHandler"
	"attendance/api/router"
	repo "attendance/repository"
	svc "attendance/services"

	"github.com/go-pg/pg"
	"github.com/google/wire"
)

func InitializeApp(*pg.DB) *router.MUXRouterImpl {
	wire.Build(
		auth.NewAuthTokenImpl, wire.Bind(new(auth.AuthToken), new(*auth.AuthTokenImpl)),
		rh.NewAddNewUserImpl, wire.Bind(new(rh.AddNewUserHandler), new(*rh.AddNewUserImpl)),
		rh.NewFetchStatusImpl, wire.Bind(new(rh.FetchStatusHandler), new(*rh.FetchStatusImpl)),
		rh.NewGetClassAttendanceImpl, wire.Bind(new(rh.GetClassAttendanceHandler), new(*rh.GetClassAttendanceImpl)),
		rh.NewGetStudentAttendanceImpl, wire.Bind(new(rh.GetStudentAttendanceHandler), new(*rh.GetStudentAttendanceImpl)),
		rh.NewGetTeacherAttendanceImpl, wire.Bind(new(rh.GetTeacherAttendanceHandler), new(*rh.GetTeacherAttendanceImpl)),
		rh.NewHomeImpl, wire.Bind(new(rh.HomeHandler), new(*rh.HomeImpl)),
		rh.NewLoginImpl, wire.Bind(new(rh.LoginHandler), new(*rh.LoginImpl)),
		rh.NewPunchInOutImpl, wire.Bind(new(rh.PunchInOutHandler), new(*rh.PunchInOutImpl)),
		repo.NewRepositoryImpl, wire.Bind(new(repo.Repository), new(*repo.RepositoryImpl)),
		svc.NewAddNewUserServiceImpl, wire.Bind(new(svc.AddNewUserService), new(*svc.AddNewUserServiceImpl)),
		svc.NewFetchStatusImpl, wire.Bind(new(svc.FetchStatusService), new(*svc.FetchStatusImpl)),
		svc.NewGetClassAttendanceImpl, wire.Bind(new(svc.GetClassAttendanceService), new(*svc.GetClassAttendanceImpl)),
		svc.NewGetStudentAttendanceServiceImpl, wire.Bind(new(svc.GetStudentAttendanceService), new(*svc.GetStudentAttendanceServiceImpl)),
		svc.NewGetTeacherAttendanceServiceImpl, wire.Bind(new(svc.GetTeacherAttendanceService), new(*svc.GetTeacherAttendanceServiceImpl)),
		svc.NewHomeServiceImpl, wire.Bind(new(svc.HomeService), new(*svc.HomeServiceImpl)),
		svc.NewPunchInOutServiceImpl, wire.Bind(new(svc.PunchInOutService), new(*svc.PunchInOutServiceImpl)),
		router.NewMUXRouterImpl)
	return &router.MUXRouterImpl{}
}
