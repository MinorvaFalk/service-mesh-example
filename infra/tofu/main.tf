provider "kubernetes" {
    config_path = var.kubernetes.config_path
    config_context = var.kubernetes.config_context
}

provider "helm" {
    kubernetes {
        config_path = var.kubernetes.config_path
        config_context = var.kubernetes.config_context
    }
}

provider "kubectl" {
    config_path = var.kubernetes.config_path
    config_context = var.kubernetes.config_context
}


# Deploy NSQ
resource "kubernetes_namespace_v1" "ns_nsq" {
    metadata {
        name = "nsq"
        annotations = {
          "linkerd.io/inject" = "enabled"
        }
    }
}

resource "helm_release" "nsq" {
    depends_on = [ kubernetes_namespace_v1.ns_nsq ]

    name = "nsq"
    namespace = "nsq"

    repository = "https://nsqio.github.io/helm-chart"
    chart = "nsq"
}

# Deploy app
data "kubectl_file_documents" "app" {
    content = file("${path.module}/../kubernetes/app.yaml")
}

resource "kubectl_manifest" "app" {
    depends_on = [ helm_release.nsq ]

    for_each = data.kubectl_file_documents.app.manifests
    yaml_body = each.value
}