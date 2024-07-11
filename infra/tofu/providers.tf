terraform {
    required_providers {
      kubernetes = {
        source = "hashicorp/kubernetes"
        version = "2.31.0"
      }
      
      helm = {
        source = "hashicorp/helm"
        version = "2.14.0"
      }

      kubectl = {
        source = "gavinbunney/kubectl"
        version = ">=1.7.0"
      }
    }
}