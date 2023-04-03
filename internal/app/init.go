package app

import (
	forumRepository "technopark_db_forum/internal/forum/repository"
	"technopark_db_forum/internal/forum/usecase"
	forumUsecase "technopark_db_forum/internal/forum/usecase"
	postRepository "technopark_db_forum/internal/posts/repository"
	postsUsecase "technopark_db_forum/internal/posts/usecase"
	serviceHandler "technopark_db_forum/internal/service/delivery"
	serviceRepository "technopark_db_forum/internal/service/repository"
	serviceUsecase "technopark_db_forum/internal/service/usecase"
	threadRepository "technopark_db_forum/internal/thread/repository"
	threadUsecase "technopark_db_forum/internal/thread/usecase"
	userRepository "technopark_db_forum/internal/users/repository"
	userUsecase "technopark_db_forum/internal/users/usecase"
	"technopark_db_forum/pkg/logger"

	"github.com/labstack/echo/v4"
	"technopark_db_forum/internal/forum/delivery"
	postsHandler "technopark_db_forum/internal/posts/delivery"
	threadHandler "technopark_db_forum/internal/thread/delivery"
	usersHandler "technopark_db_forum/internal/users/delivery"
)

type Server struct {
	Echo *echo.Echo

	forumUsecase   usecase.ForumUsecase
	usersUsecase   userUsecase.UsersUsecase
	postsUsecase   postsUsecase.PostUsecase
	threadUsecase  threadUsecase.ThreadUsecase
	serviceUsecase serviceUsecase.ServiceUsecase

	forumHandler   delivery.ForumHandler
	usersHandler   usersHandler.UserHandler
	postsHandler   postsHandler.PostHandler
	threadHandler  threadHandler.ThreadHandler
	serviceHandler serviceHandler.ServiceHandler
}

func (s *Server) init(URL string) {
	s.makeEchoLogger()
	s.makeUseCase(URL)
	s.makeHandlers()
	s.makeRouter()
}

func (s *Server) Start(host, URL string) error {
	s.init(URL)
	return s.Echo.Start(host)
}

func makeAddress(host, port string) string {
	return host + ":" + port
}

func (s *Server) makeUseCase(URL string) {
	usersRepo, err := userRepository.NewPostgres(URL)
	if err != nil {
		s.Echo.Logger.Error(err)
	}
	threadRepo, err := threadRepository.NewPostgres(URL)
	if err != nil {
		s.Echo.Logger.Error(err)
	}
	postRepo, err := postRepository.NewPostgres(URL)
	if err != nil {
		s.Echo.Logger.Error(err)
	}
	forumRepo, err := forumRepository.NewPostgres(URL)
	if err != nil {
		s.Echo.Logger.Error(err)
	}
	service, err := serviceRepository.NewPostgres(URL)
	if err != nil {
		s.Echo.Logger.Error(err)
	}

	s.usersUsecase = userUsecase.NewUserUsecase(usersRepo)
	s.threadUsecase = threadUsecase.NewThreadUsecase(threadRepo, usersRepo, forumRepo)
	s.postsUsecase = postsUsecase.NewPostUsecase(postRepo, usersRepo, threadRepo, forumRepo)
	s.forumUsecase = forumUsecase.NewUserUsecase(forumRepo, usersRepo)
	s.serviceUsecase = serviceUsecase.NewServiceUsecase(service)
}

func (s *Server) makeHandlers() {
	s.serviceHandler = serviceHandler.NewServiceHandler(s.serviceUsecase)
	s.forumHandler = delivery.NewForumHandler(s.forumUsecase)
	s.usersHandler = usersHandler.NewUserHandler(s.usersUsecase)
	s.postsHandler = postsHandler.NewPostHandler(s.postsUsecase)
	s.threadHandler = threadHandler.NewThreadHandler(s.threadUsecase)
}

func (s *Server) makeEchoLogger() {
	s.Echo.Logger = logger.GetInstance()
}

func (s *Server) makeRouter() {
	v1 := s.Echo.Group("/api")
	v1.Use(logger.Middleware())

	v1.GET("/service/status", s.serviceHandler.GetStatus)
	v1.POST("/service/clear", s.serviceHandler.Clear)

	v1.POST("/forum/create", s.forumHandler.CreateForum)
	v1.GET("/forum/:slug/details", s.forumHandler.GetForum)
	v1.GET("/forum/:slug/users", s.forumHandler.GetForumUsers)

	v1.GET("/forum/:slug/threads", s.threadHandler.GetThreadMsgs)

	v1.POST("/user/:nickname/create", s.usersHandler.CreateUser)
	v1.GET("/user/:nickname/profile", s.usersHandler.GetUser)
	v1.POST("/user/:nickname/profile", s.usersHandler.UpdateUser)

	v1.POST("/thread/:slug_or_id/create", s.postsHandler.CreatePosts)
	v1.GET("/thread/:slug_or_id/details", s.threadHandler.GetThread)
	v1.POST("/thread/:slug_or_id/details", s.threadHandler.UpdateThread)
	v1.GET("/thread/:slug_or_id/posts", s.postsHandler.GetThreadPosts)
	v1.POST("/thread/:slug_or_id/vote", s.threadHandler.CreateVote)

	v1.GET("/post/:id/details", s.postsHandler.GetPost)
	v1.POST("/post/:id/details", s.postsHandler.UpdatePost)

	v1.POST("/forum/:slug/create", s.threadHandler.CreateThread)
}

func New(echo *echo.Echo) *Server {
	return &Server{
		Echo: echo,
	}
}
