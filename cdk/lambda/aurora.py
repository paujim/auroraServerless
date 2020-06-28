import boto3

client = boto3.client('rds-data')


def on_event(event, context):
    print(event)
    request_type = event['RequestType']
    if request_type == 'Create':
        return on_create(event)
    if request_type == 'Update':
        return on_update(event)
    if request_type == 'Delete':
        return on_delete(event)
    raise Exception("Invalid request type: %s" % request_type)


def on_create(event):
    props = event["ResourceProperties"]
    print("create new resource with props %s" % props)
    physical_id = "CR-"+props["AuroraArn"]
    client.execute_statement(
        resourceArn=props["AuroraArn"],
        secretArn=props["SecretArn"],
        sql=props["Sql"],
    )

    return {'PhysicalResourceId': physical_id}


def on_update(event):
    physical_id = event["PhysicalResourceId"]
    props = event["ResourceProperties"]
    print("update resource %s with props %s" % (physical_id, props))


def on_delete(event):
    physical_id = event["PhysicalResourceId"]
    print("delete resource %s" % physical_id)
