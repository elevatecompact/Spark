environment        = "production"
region             = "us-east-1"
secondary_region   = "us-west-2"
vpc_cidr           = "10.0.0.0/16"
secondary_vpc_cidr = "10.1.0.0/16"

instance_types   = ["c5.large", "c5.xlarge"]
min_nodes        = 5
max_nodes        = 20
desired_nodes    = 8

system_node_instance = "t3.large"
app_node_instance    = "c5.large"

db_instance_class      = "db.r6g.xlarge"
db_read_replica_class  = "db.r6g.large"
db_read_replica_count  = 2
db_engine_version      = "16.3"
db_backup_retention    = 30
db_multi_az            = true

domain_name = "sparkplatform.com"
