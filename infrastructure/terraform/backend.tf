terraform {
  backend "s3" {
    bucket         = "spark-terraform-state"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "spark-terraform-locks"
  }
}
