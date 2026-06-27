environment       = "development"
region            = "us-east-1"
vpc_cidr          = "10.0.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b"]

system_node_instance = "t3.medium"
app_node_instance    = "t3.medium"
min_nodes            = 1
max_nodes            = 3
desired_nodes        = 1

db_instance_class   = "db.t3.medium"
db_engine_version   = "16.3"
db_backup_retention = 3
db_multi_az         = false

domain_name = "dev.sparkplatform.com"
