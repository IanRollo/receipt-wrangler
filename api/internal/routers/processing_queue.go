package routers

import (
	"receipt-wrangler/api/internal/handlers"
	"receipt-wrangler/api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// BuildProcessingQueueRouter serves the processing queue endpoints (swagger tag
// ProcessingQueue), mounted at /api/processingQueue.
func BuildProcessingQueueRouter() *chi.Mux {
	processingQueueRouter := chi.NewRouter()

	processingQueueRouter.Use(middleware.UnifiedAuthMiddleware)
	processingQueueRouter.Get("/group/{groupId}", handlers.GetProcessingQueueForGroup)
	processingQueueRouter.Get("/task/{queue}/{taskId}/sourceFile", handlers.GetProcessingQueueTaskSourceFile)
	processingQueueRouter.Get("/task/{queue}/{taskId}/sourceFile/download", handlers.DownloadProcessingQueueTaskSourceFile)
	processingQueueRouter.Post("/task/{queue}/{taskId}/retry", handlers.RetryProcessingQueueTask)

	return processingQueueRouter
}
