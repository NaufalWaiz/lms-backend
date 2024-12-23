package main

import (
	"log"
	"net/http"
	"project/config"
	"project/handler"
	"project/middleware"
	"project/postgres"
	"project/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db := postgres.Connect()
	defer db.Close()
	config.InitSupabase()

	forumService := service.NewForumService(db)
	forumHandler := handler.NewForumHandler(forumService)
	commentService := service.NewCommentService(db)
	commentHandler := handler.NewCommentHandler(commentService)
	authService := service.AuthService{DB: db}
	authHandler := handler.AuthHandler{AuthService: &authService}
	userService := service.UserService{DB: db}
	userHandler := handler.UserHandler{UserService: &userService}
	classService := service.ClassService{DB: db}
	classHandler := handler.ClassHandler{Service: &classService}
	materialService := service.MaterialService{DB: db}
	materialHandler := handler.MaterialHandler{Service: &materialService}
	assignmentService := service.AssignmentService{DB: db}
	assignmentHandler := handler.AssignmentHandler{Service: &assignmentService}
	rapotService := service.RapotService{DB: db}
	rapotHandler := handler.RapotHandler{Service: &rapotService}
	gradeService := service.GradeService{DB: db}
    gradeHandler := handler.GradeHandler{Service: &gradeService}

	router := mux.NewRouter()

	// JWT info Route
	router.Handle(
		"/get-token-claims",
		middleware.AuthMiddleware(http.HandlerFunc(forumHandler.GetJWTClaims)),
	).Methods("GET")

	// Forum routes
	router.Handle(
		"/forums",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(forumHandler.CreateForum)),
		),
	).Methods("POST")

	router.Handle(
		"/forums",
		middleware.AuthMiddleware(http.HandlerFunc(forumHandler.GetForums)),
	).Methods("GET")

	router.Handle(
		"/forums/{id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(forumHandler.DeleteForum)),
		),
	).Methods("DELETE")

	// Comment routes
	router.Handle(
		"/forums/{forumID}/comments",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(commentHandler.CreateComment)),
		),
	).Methods("POST")

	router.Handle(
		"/forums/{forum_id}/comments",
		middleware.AuthMiddleware(http.HandlerFunc(commentHandler.GetComments)),
	).Methods("GET")

	router.Handle(
		"/comments/{comment_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(commentHandler.DeleteComment)),
		),
	).Methods("DELETE")

	// Class Routes
	router.Handle(
		"/classes",
		middleware.AuthMiddleware(http.HandlerFunc(classHandler.GetClasses)),
	).Methods("GET")

	router.Handle(
		"/class/{id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(classHandler.GetClassByID)),
		),
	).Methods("GET")

	router.Handle(
		"/classes/count/{user_id}",
		middleware.AuthMiddleware(http.HandlerFunc(classHandler.CountClassesByUserID)),
	).Methods("GET")

	router.Handle(
		"/class",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(classHandler.CreateClass)),
		),
	).Methods("POST")

	router.Handle(
		"/class/{id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(classHandler.DeleteClass)),
		),
	).Methods("DELETE")

	router.Handle(
		"/class/{id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(classHandler.UpdateClass)),
		),
	).Methods("PUT")

	router.Handle(
		"/class/{class_id}/join",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(classHandler.JoinClass)),
		),
	).Methods("POST")

	router.Handle(
		"/class/{class_id}/members",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(classHandler.ManageClassMembers)),
		),
	).Methods("POST", "DELETE")

	router.Handle(
        "/class/{class_id}/members",
        middleware.AuthMiddleware(http.HandlerFunc(classHandler.GetMembers)),
    ).Methods("GET")

	router.Handle(
        "/classes/student/{student_id}",
        middleware.AuthMiddleware(http.HandlerFunc(classHandler.GetClassesByStudentID)),
    ).Methods("GET")

	// Material Routes
	router.Handle(
		"/materials/{class_id}",
		middleware.AuthMiddleware(http.HandlerFunc(materialHandler.GetMaterials)),
	).Methods("GET")

	router.Handle(
		"/material/{class_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(materialHandler.CreateMaterial)),
		),
	).Methods("POST")

	router.Handle(
		"/material/{id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(materialHandler.DeleteMaterial)),
		),
	).Methods("DELETE")

	router.Handle(
		"/{class_id}/material/{material_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(materialHandler.UpdateMaterial)),
		),
	).Methods("PUT")

	// Assignment Routes
	router.Handle(
		"/assignments/{class_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(assignmentHandler.GetAssignments)),
		),
	).Methods("GET")

	router.Handle(
		"/assignments/{class_id}/{user_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Siswa"})(http.HandlerFunc(assignmentHandler.GetAssignmentsByUserID)),
	),
	).Methods("GET")

	router.Handle(
		"/assignment/{class_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru", "Siswa"})(http.HandlerFunc(assignmentHandler.CreateAssignment)),
		),
	).Methods("POST")

	router.Handle(
		"/assignment/{id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(assignmentHandler.DeleteAssignment)),
		),
	).Methods("DELETE")

	router.Handle(
		"/{class_id}/assignment/{assignment_id}",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(assignmentHandler.UpdateAssignment)),
		),
	).Methods("PUT")
	
	router.Handle(
        "/assignments/count/{user_id}",
        middleware.AuthMiddleware(http.HandlerFunc(assignmentHandler.CountAssignmentsCreatedByUser)),
    ).Methods("GET")

	//grades routes
	router.Handle(
		"/grades",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware([]string{"Admin", "Guru"})(http.HandlerFunc(gradeHandler.CreateGrade)),
		),
	).Methods("POST")

	// Rapot Routes
	router.Handle(
        "/rapot/{user_id}",
        middleware.AuthMiddleware(http.HandlerFunc(rapotHandler.GetRapotByUserID)),
    ).Methods("GET")

	// Users routes
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/roles/count", userHandler.GetRoleCounts).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
