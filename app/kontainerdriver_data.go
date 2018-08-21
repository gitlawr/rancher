package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/rancher/rancher/pkg/controllers/management/kontainerdriver"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func addKontainerDrivers(management *config.ManagementContext) error {
	// create binary drop location if not exists
	err := os.MkdirAll(kontainerdriver.DriverDir, 0777)
	if err != nil {
		return fmt.Errorf("error creating binary drop folder: %v", err)
	}

	creator := driverCreator{
		driversLister: management.Management.KontainerDrivers("").Controller().Lister(),
		drivers:       management.Management.KontainerDrivers(""),
	}

	if err := creator.add("import", true, false, ""); err != nil {
		return err
	}

	if err := creator.add("rancherKubernetesEngine", true, false, ""); err != nil {
		return err
	}

	if err := creator.add("googleKubernetesEngine", true, true, ""); err != nil {
		return err
	}

	if err := creator.add("azureKubernetesService", true, true, ""); err != nil {
		return err
	}

	if err := creator.add("amazonElasticContainerService", true, true, ""); err != nil {
		return err
	}

	return nil
}

type driverCreator struct {
	driversLister v3.KontainerDriverLister
	drivers       v3.KontainerDriverInterface
}

func (c *driverCreator) add(name string, builtin, dynamicSchema bool, url string) error {
	logrus.Infof("adding kontainer driver %v", name)

	driver, err := c.driversLister.Get("", name)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = c.drivers.Create(&v3.KontainerDriver{
				ObjectMeta: v1.ObjectMeta{
					Name:      strings.ToLower(name),
					Namespace: "",
				},
				Spec: v3.KontainerDriverSpec{
					DesiredURL:    url,
					DisplayName:   name,
					BuiltIn:       builtin,
					DynamicSchema: dynamicSchema,
					Active:        true,
				},
			})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating driver: %v", err)
			}
		} else {
			return fmt.Errorf("error getting driver: %v", err)
		}
	} else {
		driver.Spec.DesiredURL = url

		_, err = c.drivers.Update(driver)
		if err != nil {
			return fmt.Errorf("error updating driver: %v", err)
		}
	}

	return nil
}
