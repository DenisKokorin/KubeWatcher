package handler

import (
	"context"
	"k8s-mon/internal/models"
	"k8s-mon/internal/monitor"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MonitorHandler struct {
	MonService monitor.Monitor
}

type MonitorService interface {
	Overview(ctx context.Context) (models.Cluster, error)
	Nodes(ctx context.Context) ([]*models.Node, error)
	Pods(ctx context.Context, namespace string) ([]*models.Pod, error)
	AddCluster(ctx context.Context, id int, name string) error
	DeleteCluster(ctx context.Context, id int) error
}

func NewMonitorHandler(m monitor.Monitor) *MonitorHandler {
	return &MonitorHandler{
		MonService: m,
	}
}

func (m *MonitorHandler) Overview(c echo.Context) error {
	cluster, err := m.MonService.GetClusterSummary(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, cluster)
}

func (m *MonitorHandler) Nodes(c echo.Context) error {
	// id := c.Param("id")
	// idInt, err := strconv.Atoi(id)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]any{
	// 		"error": "bad request",
	// 	})
	// }

	nodes, err := m.MonService.GetNodes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, nodes)
}

func (m *MonitorHandler) Pods(c echo.Context) error {
	namespace := c.QueryParam("namespace")
	// idInt, err := strconv.Atoi(id)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]any{
	// 		"error": "bad request",
	// 	})
	// }

	pods, err := m.MonService.GetPods(c.Request().Context(), namespace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, pods)
}

// func (m *MonitorHandler) AddCluster(c echo.Context) error {
// 	id, name := c.QueryParam("id"), c.QueryParam("name")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]any{
// 			"error": "bad request",
// 		})
// 	}

// 	err = m.MonService.AddCluster(c.Request().Context(), idInt, name)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]any{
// 			"error": "internal error",
// 		})
// 	}

// 	return c.JSON(http.StatusOK, id)
// }

// func (m *MonitorHandler) DeleteCluster(c echo.Context) error {
// 	id := c.Param("id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]any{
// 			"error": "bad request",
// 		})
// 	}

// 	err = m.MonService.DeleteCluster(c.Request().Context(), idInt)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]any{
// 			"error": "internal error",
// 		})
// 	}

// 	return c.JSON(http.StatusOK, id)
// }

func (m *MonitorHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}
