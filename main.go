package main

import (
	"log"
	"net/http"

	"apteka/internal/config"
	"apteka/internal/database"
	"apteka/internal/handlers"
	"apteka/internal/middleware"
	"apteka/internal/services"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}
	db, err := database.Connect(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		panic(err)
	}
	audit := &services.AuditService{
		DB: db,
	}

	stock := &services.StockService{
		DB: db,
	}

	prescription := &services.PrescriptionService{
		DB: db,
	}

	sale := &services.SaleService{
		DB:           db,
		Stock:        stock,
		Prescription: prescription,
		Audit:        audit,
	}

	authService := &services.AuthService{
		DB: db,
	}

	authHandler := &handlers.AuthHandler{
		Auth:      authService,
		JWTSecret: cfg.JWTSecret,
	}

	doctor := &handlers.DoctorHandler{
		Prescription:   prescription,
		ExpirationDays: cfg.Prescription.ExpirationDays,
	}

	employee := &handlers.EmployeeHandler{
		Sale: sale,
		DB:   db,
	}

	patient := &handlers.PatientHandler{
		DB: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.Handle("/api/doctor/create", middleware.Auth(cfg.JWTSecret)(middleware.RequireRole("doctor")(http.HandlerFunc(doctor.CreatePrescription))))
	mux.Handle("/api/employee/sell", middleware.Auth(cfg.JWTSecret)(middleware.RequireRole("employee")(http.HandlerFunc(employee.SellMedicine))))
	mux.HandleFunc("/api/patient/prescription", patient.GetPrescription)
	mux.HandleFunc("/api/patient/otc", patient.GetOTC)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(
			"PATH:",
			r.URL.Path,
		)
		w.Header().Set(
			"Content-Type",
			"text/html; charset=utf-8",
		)
		http.ServeFile(
			w,
			r,
			"./web/index.html",
		)
	},
	)
	server := &http.Server{Addr: cfg.Server.Host + ":" + cfg.Server.Port, Handler: mux}
	log.Println("running")
	log.Fatal(server.ListenAndServe())
}
