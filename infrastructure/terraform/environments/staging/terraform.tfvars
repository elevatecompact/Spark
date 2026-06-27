environment = "staging"
region      = "us-east-1"
vpc_cidr    = "10.1.0.0/16"

instance_types   = ["t3.large"]
min_nodes        = 2
max_nodes        = 5
desired_nodes    = 3

db_instance_class   = "db.t3.large"
db_engine_version   = "16.3"
db_backup_retention = 14
db_multi_az         = false

domain_name = "staging.sparkplatform.com"

system_node_instance = "t3.medium"
app_node_instance    = "t3.large"
