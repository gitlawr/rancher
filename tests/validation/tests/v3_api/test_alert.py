import pytest
import requests

from .common import *  # NOQA
from .alert_common import MockReceiveAlert

dingtalk_config = {
    "type"  : "/v3/schemas/dingtalkConfig",
    "url"   : "http://127.0.0.1:4050/dingtalk/",
}

microsoft_teams_config = {
    "type" : "/v3/schemas/msTeamsConfig",
    "url"  : "http://127.0.0.1:4050/microsoftTeams",
}

MOCK_RECEIVER_ALERT_PORT = 4050
CLUSTER_ALERTING_APP = "cluster-alerting"
ALERTMANAGER_CLUSTER_ALERTING_WORKLOAD  = "alertmanager-cluster-alerting"  	
WEBHOOK_RECEIVER_WORKLOAD = "webhook-receiver-cluster-alerting"

@pytest.fixture(scope="module")
def mock_receiver_alert():
    server = MockReceiveAlert(port=MOCK_RECEIVER_ALERT_PORT)
    server.start()
    yield server
    server.shutdown_server()

def test_add_notifier(remove_resource, mock_receiver_alert):
    client, cluster = get_global_admin_client_and_cluster()
    nodes = client.list_node(clusterId=cluster.id)
    assert len(nodes.data) > 0
    node_id = nodes.data[0].id

    # Add the notifier dingtalk and microsoftTeams 
    notifier_dingtalk = client.create_notifier(name="dingtalk",
                                                clusterId=cluster.id,
                                                dingtalkConfig=dingtalk_config
                                              )

    notifier_microsoft_teams = client.create_notifier(name="microsoftTeams",
                                                 clusterId=cluster.id,
                                                 msteamsConfig=microsoft_teams_config
                                                )  

    client.action(obj=notifier_microsoft_teams,
                action_name="send",
                msteamsConfig=microsoft_teams_config)

    client.action(obj=notifier_dingtalk,
                action_name="send",
                dingtalkConfig=dingtalk_config)

    # Add the alertGroup
    all_recipients = []

    recipient1 = {
        "notifierId"   :  notifier_dingtalk.id,
        "notifierType" :  "dingtalk",
        "type"         :  "/v3/schemas/recipient"
    }

    recipient2 = {
        "notifierId"   :  notifier_microsoft_teams.id,
        "notifierType" :  "msteams",
        "type"         :  "/v3/schemas/recipient"
    }

    all_recipients.append(recipient1)
    all_recipients.append(recipient2)

    cluster_alert_group = client.create_cluster_alert_group(name="testGroup",
                                                           clusterId=cluster.id,
                                                           groupIntervalSeconds=180,
                                                           groupWaitSeconds=30,
                                                           recipients=all_recipients
                                                           ) 

    # Add the alertRule  
    node_rule = {
        "condition"    :  "mem",
        "cpuThreshold" :  70,
        "memThreshold" :  1,
        "nodeId"       :  node_id,
        "type"         :  "/v3/schemas/nodeRule"
    }

    cluster_alert_rule = client.create_cluster_alert_rule(name="testRule",
                                                         clusterId=cluster.id,
                                                         groupId=cluster_alert_group.id,
                                                         nodeRule=node_rule
                                                        ) 

    # Make sure the webhook-receive is work 
    system_project = client.list_project(name="System",
                                         clusterId=cluster.id).data[0]

    project_client = get_project_client_for_token(system_project, ADMIN_TOKEN)
    wait_for_app_to_active(project_client, CLUSTER_ALERTING_APP)
    alerting_app = project_client.list_app(name=CLUSTER_ALERTING_APP).data[0]

    assert alerting_app.answers["webhook-receiver.enabled"] == "true"

    # Make sure the workload is active
    webhook_receiver = project_client.list_workload(name=WEBHOOK_RECEIVER_WORKLOAD).data[0]
    assert webhook_receiver.state == "active"

    alertmanager_cluster_alerting = project_client.list_workload(name=ALERTMANAGER_CLUSTER_ALERTING_WORKLOAD).data[0]
    assert alertmanager_cluster_alerting.state == "active"

    # Remove the alertGroup
    remove_resource(cluster_alert_group)

    # Remove the notifiers
    remove_resource(notifier_dingtalk)
    remove_resource(notifier_microsoft_teams) 