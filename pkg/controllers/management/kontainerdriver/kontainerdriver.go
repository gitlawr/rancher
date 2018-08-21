package kontainerdriver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rancher/kontainer-engine/service"
	"github.com/rancher/kontainer-engine/types"
	"github.com/rancher/rancher/pkg/controllers/management/clusterprovisioner"
	"github.com/rancher/rancher/pkg/controllers/management/nodedriver"
	"github.com/rancher/types/apis/core/v1"
	corev1 "github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	driverNameLabel = "io.cattle.node_driver.name"
	DriverDir       = "./management-state/kontainer-drivers/"
)

func Register(management *config.ManagementContext) {
	lifecycle := &Lifecycle{
		dynamicSchemas:        management.Management.DynamicSchemas(""),
		dynamicSchemasLister:  management.Management.DynamicSchemas("").Controller().Lister(),
		namespaces:            management.Core.Namespaces(""),
		coreV1:                management.Core,
		kontainerDriverLister: management.Management.KontainerDrivers("").Controller().Lister(),
	}

	management.Management.KontainerDrivers("").AddLifecycle("mgmt-kontainer-driver-lifecycle", lifecycle)
}

type Lifecycle struct {
	dynamicSchemas        v3.DynamicSchemaInterface
	dynamicSchemasLister  v3.DynamicSchemaLister
	namespaces            v1.NamespaceInterface
	coreV1                corev1.Interface
	kontainerDriverLister v3.KontainerDriverLister
}

func (l *Lifecycle) Create(obj *v3.KontainerDriver) (*v3.KontainerDriver, error) {
	logrus.Infof("create kontainerdriver %v", obj.Name)

	v3.KontainerDriverConditionDownloaded.False(obj)

	if !obj.Spec.Active {
		return obj, nil
	}

	var path string
	if !obj.Spec.BuiltIn {
		// download file
		var err error
		path, err = download(obj)
		if err != nil {
			return nil, err
		}

		obj.Status.ExecutablePath = path
		v3.KontainerDriverConditionDownloaded.True(obj)

		logrus.Infof("kontainerdriver %v downloaded and registered at %v", obj.Name, path)
	}

	// if dynamic schema is set to false then do not add this driver's schema to the list
	if !obj.Spec.DynamicSchema {
		return obj, nil
	}

	err := l.createSchema(obj)

	return obj, err
}

func (l *Lifecycle) createSchema(obj *v3.KontainerDriver) error {
	driver := service.NewEngineService(
		clusterprovisioner.NewPersistentStore(l.namespaces, l.coreV1),
	)
	flags, err := driver.GetDriverCreateOptions(context.Background(), obj.Name, obj, v3.ClusterSpec{
		GenericEngineConfig: &v3.MapStringInterface{
			clusterprovisioner.DriverNameField: obj.Name,
		},
	})
	if err != nil {
		return fmt.Errorf("error getting driver create options: %v", err)
	}

	resourceFields := map[string]v3.Field{}
	for key, flag := range flags.Options {
		formattedName, field, err := toResourceField(key, flag)
		if err != nil {
			return fmt.Errorf("error formatting field name: %v", err)
		}

		resourceFields[formattedName] = field
	}

	// all drivers need a driverName field so kontainer-engine knows what type they are
	resourceFields[clusterprovisioner.DriverNameField] = v3.Field{
		Create: true,
		Update: true,
		Type:   "string",
		Default: v3.Values{
			StringValue: obj.Name,
		},
	}

	dynamicSchema := &v3.DynamicSchema{
		Spec: v3.DynamicSchemaSpec{
			ResourceFields: resourceFields,
		},
	}
	dynamicSchema.Name = getDynamicTypeName(obj)
	dynamicSchema.OwnerReferences = []v13.OwnerReference{
		{
			UID:        obj.UID,
			Kind:       obj.Kind,
			APIVersion: obj.APIVersion,
			Name:       obj.Name,
		},
	}
	dynamicSchema.Labels = map[string]string{}
	dynamicSchema.Labels[obj.Name] = obj.Spec.DisplayName
	dynamicSchema.Labels = map[string]string{}
	dynamicSchema.Labels[driverNameLabel] = obj.Spec.DisplayName
	_, err = l.dynamicSchemas.Create(dynamicSchema)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("error creating dynamic schema: %v", err)
	}

	return l.createOrUpdateKontainerDriverTypes(obj)
}

func (l *Lifecycle) createOrUpdateKontainerDriverTypes(obj *v3.KontainerDriver) error {
	nodedriver.SchemaLock.Lock()
	defer nodedriver.SchemaLock.Unlock()

	embeddedType := getDynamicTypeName(obj)
	fieldName := getDynamicFieldName(obj)

	nodeSchema, err := l.dynamicSchemasLister.Get("", "cluster")
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		resourceField := map[string]v3.Field{}
		resourceField[fieldName] = v3.Field{
			Create:   true,
			Nullable: true,
			Update:   true,
			Type:     embeddedType,
		}

		dynamicSchema := &v3.DynamicSchema{}
		dynamicSchema.Name = "cluster"
		dynamicSchema.Spec.ResourceFields = resourceField
		dynamicSchema.Spec.Embed = true
		dynamicSchema.Spec.EmbedType = "cluster"
		_, err := l.dynamicSchemas.Create(dynamicSchema)
		if err != nil {
			return err
		}
		return nil
	}

	shouldUpdate := false

	if nodeSchema.Spec.ResourceFields == nil {
		nodeSchema.Spec.ResourceFields = map[string]v3.Field{}
	}
	if _, ok := nodeSchema.Spec.ResourceFields[fieldName]; !ok {
		// if embedded we add the type to schema
		nodeSchema.Spec.ResourceFields[fieldName] = v3.Field{
			Create:   true,
			Nullable: true,
			Update:   true,
			Type:     embeddedType,
		}
		shouldUpdate = true
	}

	if shouldUpdate {
		_, err = l.dynamicSchemas.Update(nodeSchema)
		if err != nil {
			return err
		}
	}

	return nil
}

func toResourceField(name string, flag *types.Flag) (string, v3.Field, error) {
	field := v3.Field{
		Create: true,
		Update: true,
		Type:   "string",
	}

	name, err := toLowerCamelCase(name)
	if err != nil {
		return name, field, err
	}

	field.Description = flag.Usage

	if flag.Type == types.StringType {
		field.Default.StringValue = flag.Value

		if flag.Password {
			field.Type = "password"
		} else {
			field.Type = "string"
		}

		if flag.Default != nil {
			field.Default.StringValue = flag.Default.DefaultString
		}
	} else if flag.Type == types.IntType {
		field.Type = "int"

		if flag.Default != nil {
			field.Default.StringValue = strconv.Itoa(int(flag.Default.DefaultInt))
		}
	} else if flag.Type == types.BoolType || flag.Type == types.BoolPointerType {
		field.Type = "boolean"

		if flag.Default != nil {
			field.Default.BoolValue = flag.Default.DefaultBool
		}
	} else if flag.Type == types.StringSliceType {
		field.Type = "array[string]"

		if flag.Default != nil {
			field.Default.StringSliceValue = flag.Default.DefaultStringSlice.Value
		}
	} else {
		return name, field, fmt.Errorf("unknown type of flag %v: %v", flag, reflect.TypeOf(flag))
	}

	return name, field, nil
}

func toLowerCamelCase(nodeFlagName string) (string, error) {
	flagNameParts := strings.Split(nodeFlagName, "-")
	flagName := flagNameParts[0]
	for _, flagNamePart := range flagNameParts[1:] {
		flagName = flagName + strings.ToUpper(flagNamePart[:1]) + flagNamePart[1:]
	}
	return flagName, nil
}

func binaryPath(obj *v3.KontainerDriver) string {
	return DriverDir + obj.Name
}

func download(obj *v3.KontainerDriver) (string, error) {
	path := binaryPath(obj)
	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("error creating binary file %v", err)
	}
	defer out.Close()

	err = os.Chmod(path, 0777)
	if err != nil {
		return "", fmt.Errorf("error changing file permissions: %v", err)
	}

	resp, err := http.Get(obj.Spec.DesiredURL)
	if err != nil {
		return "", fmt.Errorf("error on http get: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %v", err)
	}

	obj.Status.ActualURL = obj.Spec.DesiredURL

	return path, nil
}

func (l *Lifecycle) Updated(obj *v3.KontainerDriver) (*v3.KontainerDriver, error) {
	logrus.Infof("update kontainerdriver %v", obj.Name)
	if obj.Spec.BuiltIn {
		return obj, nil
	}

	if !obj.Spec.DynamicSchema {
		return obj, nil
	}

	// redownload file if url changed or not downloaded
	if obj.Spec.DesiredURL != obj.Status.ActualURL || v3.KontainerDriverConditionDownloaded.IsFalse(obj) {
		path, err := download(obj)
		if err != nil {
			return nil, err
		}

		obj.Status.ExecutablePath = path
		obj.Status.ActualURL = obj.Spec.DesiredURL
		v3.KontainerDriverConditionDownloaded.True(obj)
	}

	var err error
	if obj.Spec.Active {
		err = l.createSchema(obj)
	} else {
		if err = l.dynamicSchemas.Delete(getDynamicTypeName(obj), &v13.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return nil, fmt.Errorf("error deleting schema: %v", err)
		}

		if err = l.removeFieldFromCluster(obj); err != nil {
			return nil, err
		}
	}

	return obj, err
}

func getDynamicTypeName(obj *v3.KontainerDriver) string {
	return obj.Name + "config"
}

func getDynamicFieldName(obj *v3.KontainerDriver) string {
	return obj.Spec.DisplayName + "Config"
}

func (l *Lifecycle) Remove(obj *v3.KontainerDriver) (*v3.KontainerDriver, error) {
	logrus.Infof("remove kontainerdriver %v", obj.Name)

	// delete file
	if err := os.Remove(binaryPath(obj)); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error deleting driver binary: %v", err)
	}

	v3.KontainerDriverConditionDownloaded.False(obj)

	if err := l.dynamicSchemas.Delete(getDynamicTypeName(obj), &v13.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("error deleting dynamic schema: %v", err)
	}

	if obj.Spec.DynamicSchema {
		if err := l.removeFieldFromCluster(obj); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (l *Lifecycle) removeFieldFromCluster(obj *v3.KontainerDriver) error {
	nodedriver.SchemaLock.Lock()
	defer nodedriver.SchemaLock.Unlock()

	fieldName := getDynamicFieldName(obj)

	nodeSchema, err := l.dynamicSchemasLister.Get("", "cluster")
	if err != nil {
		return fmt.Errorf("error getting schema: %v", err)
	}

	delete(nodeSchema.Spec.ResourceFields, fieldName)

	if _, err = l.dynamicSchemas.Update(nodeSchema); err != nil {
		return fmt.Errorf("error removing schema from cluster: %v", err)
	}

	return nil
}
