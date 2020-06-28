import json
from aws_cdk import (
    aws_ec2 as ec2,
    aws_rds as rds,
    aws_secretsmanager as secretsmanager,
    core,
)


class ServerlessAuroraStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        vpc = ec2.Vpc(
            scope=self,
            id="aurora-VPC",
            cidr="10.10.0.0/16"
        )

        templated_secret = secretsmanager.Secret(
            scope=self,
            id="templated-secret",
            generate_secret_string=secretsmanager.SecretStringGenerator(
                secret_string_template=json.dumps(
                    {"username": "testuser"}),
                generate_string_key="password",
                exclude_punctuation=True,
            )
        )

        db_subnets = []
        for sn in vpc.private_subnets:
            db_subnets.append(sn.subnet_id)

        cfn_db_subnets = rds.CfnDBSubnetGroup(
            scope=self,
            id="DB-subnet-group",
            db_subnet_group_description="subnet group",
            subnet_ids=db_subnets,
        )

        cfn_cluster = rds.CfnDBCluster(
            scope=self,
            id="db-cluster",
            db_cluster_identifier="serverless-cluster",
            master_username=templated_secret.secret_value_from_json(
                "username").to_string(),
            master_user_password=templated_secret.secret_value_from_json(
                "password").to_string(),
            engine="aurora",
            engine_mode="serverless",
            scaling_configuration=rds.CfnDBCluster.ScalingConfigurationProperty(
                auto_pause=True,
                min_capacity=4,
                max_capacity=8,
                seconds_until_auto_pause=1000,
            ),
            deletion_protection=False,
            db_subnet_group_name=cfn_db_subnets.ref,
        )

        core.CfnOutput(
            scope=self,
            id="aurora-arn",
            value=f"arn:aws:rds:{core.Aws.REGION}:{core.Aws.ACCOUNT_ID}:cluster:{cfn_cluster.ref}"
        )
