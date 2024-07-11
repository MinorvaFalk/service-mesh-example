# Service Mesh Example

This is a implementation demo of notification service using distributed messaging and service mesh. Project consist of 2 application, **producer** and **consumer** which combined into single docker image.

**Producer** is an API for creating and sending notification job that consist of *title*, *body*, and *image url*. After creating job, it will be sent to NSQ first then consumer.

**Consumer** is a downstream service for processing job received from producer via NSQ. Currently consumer does not have any business logic, only printing the message received from NSQ. You can try to implement or modify business logic inside `/internal/handler/handler.go`.

## Requirements

- [Go](https://go.dev/doc/install) >=1.22.5
- [Docker](https://docs.docker.com/get-docker/)
- Kubernetes Cluster
  > Preferably local kubernetes cluster such as [kind](https://kind.sigs.k8s.io/) or [minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Flinux%2Fx86-64%2Fstable%2Fbinary+download)
- [Linkerd](https://linkerd.io/2.15/getting-started/)
- [OpenTofu](https://opentofu.org/docs/intro/install/) or [Terraform](https://developer.hashicorp.com/terraform/install?product_intent=terraform)

## How to Run

You need to prepare config file named `config.yaml`, for variables please refer [`example.config.yaml`](/example.config.yaml)

```bash
# running go app

make run-consumer
make run-producer


# building image

make build IMAGE_NAME=service_mesh IMAGE_TAG=1.0.0
# or if you dont have make installed
docker build -f infra/docker/Dockerfile -t service-mesh:0.0.0 .

# running docker image
docker run -itd \
--name service-mesh-producer \
-p 8080:8080 \
-v $(PWD)/config.yaml:/config.yaml \
service-mesh:0.0.0 /producer
```

For deploying on kubernetes, you need to modify certain code in [**infra/kubernetes/app.yaml**](/infra/kubernetes/app.yaml) such as `spec.containers.image` and modify `kind: ConfigMap`.

After modifying data for kubernetes, you need to create `values.tfvars` which contain variables from [**variables.tf**](/infra/tofu/variables.tf).

```bash
# values.tfvars example
kubernetes = {
    config_path = "~/.kube/config"
    config_context = "default"
}
```

Loading image into local kubernetes cluster

```bash
# for kind, reference https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster
kind load docker-image service-mesh:0.0.0
```

Install Linkerd into cluster [reference](https://linkerd.io/2.15/getting-started/)

```bash
linkerd check --pre
linkerd install --crds | kubectl apply -f -
linkerd install | kubectl apply -f -
linkerd check
```

After installing linkerd, we can automatically inject app by annotating namespace using this command. In this project we already included annotation inside kubernetes file and terraform file.

```bash
kubectl create namespace <your_namespace>
kubectl annotate namespace/<your_namespace> "linkerd.io/inject: enabled"
```

Before deploying app inside /infra/tofu, you need to configure few things inside [**infra/kubernetes/app.yaml**](./infra/kubernetes/app.yaml) such as

- `spec.containers.image` in Deployment
- `data.config.yaml` in ConfigMap

After configuring app.yaml, you need to create a `.tfvars` file for OpenTofu variable file. Please refer [**infra/tofu/variables.tf**](./infra/tofu/values.tf).

```bash
# Example variables.tfvars
kubernetes = {
    config_path = "~/.kube/config"
    config_context = "kind-kind"
}
```

```bash
# Show result before deploying
tofu plan -var-file=variables.tfvars

# Apply configuration
tofu apply -var-file=variables.tfvars -auto-approve

# Destroy
tofu destory -var-file=variables.tfvars
```