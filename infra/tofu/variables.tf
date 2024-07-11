variable "kubernetes" {
    type = object({
      config_path = string
      config_context = string
    })

    default = {
      config_path = "~/.kube/config"
      config_context = "default"
    }
}