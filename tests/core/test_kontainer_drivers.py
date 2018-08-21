import sys

from .conftest import wait_for_condition, wait_until

DRIVER_URL = "https://github.com/rancher/kontainer-engine-" + \
             "example-driver/releases/download/v0.1.0/kontainer-engine-" + \
             "example-driver-" + sys.platform


def test_builtin_drivers_are_present(admin_mc):
    admin_mc.client.reload_schema()
    types = admin_mc.client.schema.types

    assert 'azurekubernetesserviceconfig' in types
    assert 'googlekubernetesengineconfig' in types
    assert 'amazonelasticcontainerserviceconfig' in types


def test_can_add_driver(admin_mc, remove_resource):
    kd = admin_mc.client.create_kontainerDriver(
        name="gke-driver-1",
        displayName="gkeDriver1",
        dynamicSchema=True,
        active=True,
        desirdUrl=DRIVER_URL
    )
    remove_resource(kd)

    wait_for_condition('Downloaded', 'True', admin_mc.client, kd)

    verify_driver_in_types(admin_mc.client, kd)


def test_enabling_driver_exposes_schema(admin_mc, remove_resource):
    kd = admin_mc.client.create_kontainerDriver(
        name="gke-driver-2",
        displayName="gkeDriver2",
        dynamicSchema=True,
        active=False,
        desirdUrl=DRIVER_URL
    )
    remove_resource(kd)

    verify_driver_not_in_types(admin_mc.client, kd)

    kd.active = True
    admin_mc.client.update_by_id_kontainerDriver(kd.id, kd)

    wait_for_condition('Downloaded', 'True', admin_mc.client, kd)

    verify_driver_in_types(admin_mc.client, kd)


def test_disabling_driver_hides_schema(admin_mc, remove_resource):
    kd = admin_mc.client.create_kontainerDriver(
        name="gke-driver-3",
        displayName="gkeDriver3",
        dynamicSchema=True,
        active=True,
        desirdUrl=DRIVER_URL
    )
    remove_resource(kd)

    wait_for_condition('Downloaded', 'True', admin_mc.client, kd)

    verify_driver_in_types(admin_mc.client, kd)

    kd.active = False
    admin_mc.client.update_by_id_kontainerDriver(kd.id, kd)

    verify_driver_not_in_types(admin_mc.client, kd)


def test_can_delete_driver(admin_mc, remove_resource):
    kd = admin_mc.client.create_kontainerDriver(
        name="gke-driver-4",
        displayName="gkeDriver4",
        dynamicSchema=True,
        active=True,
        desirdUrl=DRIVER_URL
    )
    remove_resource(kd)

    wait_for_condition('Downloaded', 'True', admin_mc.client, kd)

    verify_driver_in_types(admin_mc.client, kd)

    admin_mc.client.delete(kd)

    verify_driver_not_in_types(admin_mc.client, kd)


def verify_driver_in_types(client, kd):
    def check():
        client.reload_schema()
        types = client.schema.types
        return kd.id + 'config' in types

    wait_until(check)
    assert kd.id + 'config' in client.schema.types


def verify_driver_not_in_types(client, kd):
    def check():
        client.reload_schema()
        types = client.schema.types
        return kd.id + 'config' not in types

    wait_until(check)
    assert kd.id + 'config' not in client.schema.types
