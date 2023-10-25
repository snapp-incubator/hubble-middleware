package handler

import (
	"context"
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/auth"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/hubble-middleware/resp"
	"k8s.io/client-go/rest"

	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	groupv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	ProjectHandler struct {
		k8sUserConfig rest.Config
		k8sAppConfig  rest.Config
	}
)

func NewProject(k8s rest.Config) *ProjectHandler {
	return &ProjectHandler{k8sUserConfig: k8s, k8sAppConfig: k8s}
}

func (h *ProjectHandler) Get(c echo.Context) error {
	user, ok := c.Get("user").(auth.User)
	if !ok {
		log.Error("Unauthorized user")
		return echo.ErrUnauthorized
	}

	groups, err := h.getUserGroups(user.Username)
	if err != nil {
		log.Errorf("Get Groups Error: %s", err)
		return echo.ErrInternalServerError
	}

	projects, err := h.getUserProjects(user.Username, groups)
	if err != nil {
		log.Errorf("Get Projects Error: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, resp.User{
		Username: user.Username,
		Projects: projects,
	})
}

func (h *ProjectHandler) getUserProjects(username string, groups []string) ([]string, error) {
	h.k8sUserConfig.Impersonate.UserName = username

	if len(groups) > 0 {
		h.k8sUserConfig.Impersonate.Groups = groups
	}

	projectClientset, err := projectv1.NewForConfig(&h.k8sUserConfig)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	res, err := projectClientset.Projects().List(context.Background(), metav1.ListOptions{})
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

func (h *ProjectHandler) getUserGroups(username string) ([]string, error) {
	groupClientset, err := groupv1.NewForConfig(&h.k8sAppConfig)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	res, err := groupClientset.Groups().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	groups := []string{}
	for _, item := range res.Items {
		if slices.Contains(item.Users, username) {
			groups = append(groups, item.ObjectMeta.Name)
		}
	}

	return groups, err
}
