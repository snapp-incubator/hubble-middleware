package api

import (
	"errors"
	"fmt"
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/auth"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/hubble-middleware/handler"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Claims struct {
	Username string
	Iat      int64
	Exp      int64
}

func main(cfg config.Config) {
	app := echo.New()

	clusterConfig, err := getClusterConfig()
	if err != nil {
		log.Errorf("", err)
		os.Exit(1)
	}

	projectsHandler := handler.NewProject(*clusterConfig)

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead,
			http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	app.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	projects := app.Group("/projects")

	projects.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			user, err := auth.Validate(token, cfg.Dex.Secret)
			if err != nil {
				return false, err
			}

			c.Set("user", user)

			return true, nil

		},
	}))
	projects.GET("", projectsHandler.Get)

	if err := app.Start(fmt.Sprintf(":%d", cfg.API.Port)); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("echo initiation failed: %s", err)
	}

	logrus.Println("API has been started :D")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func getClusterConfig() (*rest.Config, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err == nil {
		return clusterConfig, nil
	}

	kubeconfig, err := kubeconfigLocation()
	if err != nil {
		return nil, err
	}

	clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return clusterConfig, nil
}

func kubeconfigLocation() (string, error) {
	value, present := os.LookupEnv("KUBECONFIG")
	if present {
		fileExist, err := exists(value)
		if err != nil {
			return "", err
		}
		if fileExist {
			return value, nil
		}
	}
	return filepath.Join(os.Getenv("HOME"), ".kube", "config"), nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// Register API command.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "api",
			Short: "Run API to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
