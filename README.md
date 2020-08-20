#SETTING UP EXAMPLES AND INSTALL INSTRUCTIONS FOR RabbitMQ on Kubernetes

After practicing the basic examples, follow the following steps to Install RMQ on cluster
- NOTE THE FOLLOWING PRODUCTION CHECKLIST
- https://www.rabbitmq.com/production-checklist.html


##Install RabbitMQ on Kubernetes Using The RabbitMQ Cluster Operator

- git clone git@github.com:rabbitmq/cluster-operator.git
- cd cluster-operator
- kubectl create -f config/namespace/base/namespace.yaml
- kubectl create -f config/crd/bases/rabbitmq.com_rabbitmqclusters.yaml
- kubectl -n rabbitmq-system create --kustomize config/rbac/
- kubectl -n rabbitmq-system create --kustomize config/manager/


