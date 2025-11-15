package di

import (
	"log"

	maxbot "github.com/rectid/max-bot-api-client-go"

	"maxBot/internal/repository"
	"maxBot/internal/service"
)

// Services контейнер зависимостей для всех сервисов и внешних зависимостей
type Services struct {
	AdminService       service.AdminService
	ApplicationService service.VolunteerApplicationService
	CategoryService    service.CategoryService
	EventService       service.EventService
	ImageService       service.EventMediaService
	OrganizerService   service.OrganizerService
	UserService        service.UserService
	VolunteerService   service.VolunteerService
	API                *maxbot.Api
}

// NewServices инициализирует все сервисы и возвращает контейнер зависимостей
func NewServices(api *maxbot.Api, repo *repository.Repository) *Services {
	if repo == nil {
		log.Fatal("repository is required")
	}

	queries := repo.Queries()

	adminService := service.NewAdminService(queries)
	applicationService := service.NewVolunteerApplicationService(queries)
	categoryService := service.NewCategoryService(queries)
	eventService := service.NewEventService(queries)
	imageService := service.NewEventMediaService(queries)
	organizerService := service.NewOrganizerService(queries, api)
	userService := service.NewUserService(queries)
	volunteerService := service.NewVolunteerService(queries)

	return &Services{
		AdminService:       adminService,
		ApplicationService: applicationService,
		CategoryService:    categoryService,
		EventService:       eventService,
		ImageService:       imageService,
		OrganizerService:   organizerService,
		UserService:        userService,
		VolunteerService:   volunteerService,
		API:                api,
	}
}
