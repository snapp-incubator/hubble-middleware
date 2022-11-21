package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/hubble-middleware/resp"
	"k8s.io/client-go/rest"

	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	ProjectHandler struct {
		k8sClusterConfig *rest.Config
	}

	jwtCustomClaims struct {
		jwt.Claims

		Username string `json:"username"`
	}
)

func NewProject(k8s *rest.Config) *ProjectHandler {
	return &ProjectHandler{k8sClusterConfig: k8s}
}

func (h *ProjectHandler) Get(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.ErrUnauthorized
	}

	claims := user.Claims.(*jwtCustomClaims)

	projects, err := h.getProjects(claims.Username)
	if err != nil {
		log.Errorf("Get Projects Error: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, resp.User{
		Username: claims.Username,
		Projects: projects,
	})
}

func (h *ProjectHandler) getProjects(username string) ([]string, error) {
	h.k8sClusterConfig.Impersonate.UserName = username
	projectClientset, err := projectv1.NewForConfig(h.k8sClusterConfig)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	res, err := projectClientset.Projects().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	//projects := make(map[string]struct{})
	projects := []string{}
	for _, item := range res.Items {
		projects = append(projects, item.ObjectMeta.Name)
	}

	return projects, err
}
