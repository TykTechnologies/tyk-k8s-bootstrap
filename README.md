# Tyk-K8S-Bootstrap

This is a standalone app meant to help with bootstrap and deletion of the [tyk-stack](https://github.com/TykTechnologies/tyk-charts/tree/main/tyk-stack) when installed
via Helm Charts

> [!NOTE]
This app is needed only for [Tyk Self-Managed](https://tyk.io/docs/tyk-on-premises/) deployment!
<br/>[Tyk Open Source](https://tyk.io/docs/apim/open-source/) doesn't have a special bootstrap and in [Tyk Cloud](https://tyk.io/docs/tyk-cloud/) it is done for you (being a SaaS).

## What it does?

`tyk-k8s-bootstrap` offers three applications functioning as [Chart Hooks](https://helm.sh/docs/topics/charts_hooks/) in Helm charts.

- `bootstrap-pre-install` is a binary functioning as a `pre-install` hook, validating the Tyk Dashboard License key.
- `bootstrap-post-install` is a binary functioning as a `post-install` hook, bootstrapping the Tyk Dashboard by 
setting up an organization and an admin user. Additionally, it generates Kubernetes secrets utilized by 
[Tyk Operator](https://tyk.io/docs/tyk-operator/) and [Tyk Enterprise Portal](https://tyk.io/docs/tyk-stack/tyk-developer-portal/enterprise-developer-portal/install-tyk-enterprise-portal/)
- `bootstrap-pre-delete` is a binary functioning as a `pre-delete` hook, responsible for system cleanup.

## Environment Variables

| Environment Variable                           | Description                                                                                                                                                                                                           |
|------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| TYK_K8SBOOTSTRAP_INSECURESKIPVERIFY            | enables InsecureSkipVerify options in HTTP requests sent to Tyk -<br/> might be useful for Tyk Dashboard with self-signed certs                                                                                       |
| TYK_K8SBOOTSTRAP_BOOTSTRAPDASHBOARD            | controls bootstrapping Tyk Dashboard or not.                                                                                                                                                                          |
| TYK_K8SBOOTSTRAP_BOOTSTRAPPORTAL               | controls bootstrapping Tyk Classic Portal or not.                                                                                                                                                                     |
| TYK_K8SBOOTSTRAP_OPERATORKUBERNETESSECRETNAME  | corresponds to the Kubernetes secret name that will be created for Tyk Operator.<br/> Set it to an empty to string to disable bootstrapping Kubernetes secret for Tyk Operator.                                       |
| TYK_K8SBOOTSTRAP_DEVPORTALKUBERNETESSECRETNAME | corresponds to the Kubernetes secret name that will be created for Tyk Developer Enterprise Portal.<br/> Set it to an empty to string to disable bootstrapping Kubernetes secret for Tyk Developer Enterprise Portal. |
| TYK_K8SBOOTSTRAP_K8S_DASHBOARDSVCURL           | corresponds to the URL of Tyk Dashboard.                                                                                                                                                                              |
| TYK_K8SBOOTSTRAP_K8S_DASHBOARDSVCPROTO         | corresponds to Tyk Dashboard Service Protocol (either http or https). By default, it is http.                                                                                                                         |
| TYK_K8SBOOTSTRAP_K8S_RELEASENAMESPACE          | corresponds to the namespace where Tyk is deployed via Helm Chart.                                                                                                                                                    |
| TYK_K8SBOOTSTRAP_K8S_DASHBOARDDEPLOYMENTNAME   | corresponds to the name of the Tyk Dashboard Deployment, which is being used to restart<br/> Dashboard pod after bootstrapping.                                                                                       |
| TYK_K8SBOOTSTRAP_TYK_ADMIN_SECRET              | corresponds to the secret that will be used in Admin APIs.                                                                                                                                                            |
| TYK_K8SBOOTSTRAP_TYK_ADMIN_FIRSTNAME           | corresponds to the first name of the admin being created.                                                                                                                                                             |
| TYK_K8SBOOTSTRAP_TYK_ADMIN_LASTNAME            | corresponds to the last name of the admin being created.                                                                                                                                                              |
| TYK_K8SBOOTSTRAP_TYK_ADMIN_EMAILADDRESS        | corresponds to the email address of the admin being created.                                                                                                                                                          |
| TYK_K8SBOOTSTRAP_TYK_ADMIN_PASSWORD            | corresponds to the password of the admin being created.                                                                                                                                                               |
| TYK_K8SBOOTSTRAP_TYK_ADMIN_AUTH                | corresponds to Tyk Dashboard API Access Credentials of the admin user, and it will be used in Authorization <br/>header of the HTTP requests that will be sent to Tyk for bootstrapping.                              |
| TYK_K8SBOOTSTRAP_TYK_ORG_NAME                  | corresponds to the name for your organization that is going to be bootstrapped in Tyk.                                                                                                                                |
| TYK_K8SBOOTSTRAP_TYK_ORG_CNAME                 | corresponds to the Organisation CNAME which is going to bind the Portal to.                                                                                                                                           |
| TYK_K8SBOOTSTRAP_TYK_ORG_ID                    | corresponds to the organisation ID that is being created.                                                                                                                                                             |
| TYK_K8SBOOTSTRAP_TYK_DASHBOARDLICENSE          | corresponds to the license key of Tyk Dashboard.                                                                                                                                                                      |

## Required RBAC roles for the app to work inside the Kubernetes cluster
- delete
- list

## Useful debug/test tips/commands:

### Load images to Kind Cluster

After making your changes to applications, running the following command loads your local changes into KinD cluster with `tykio/tyk-k8s-boostrap{pre-post}-{delete-install}:testing` image.

```bash
$ ./hack/load_images.sh
```

### KinD with a local image repository

If you want to create a k8s kind cluster that also has a local repository where
you can push the images generated by the Makefile just run the "local_registry.sh" script.
After, the commands below help you with building and pushing the containers to the local repository.


```bash
(rm bin/bootstrapapp-post || true) && make build-bootstrap-post && docker build -t localhost:5001/bootstrap-tyk-post:$bsVers -f ./.container/image/bootstrap-post/Dockerfile ./bin && docker push localhost:5001/bootstrap-tyk-post:$bsVers
```
```bash
(rm bin/bootstrapapp-pre-delete || true) && make build-bootstrap-pre-delete && docker build -t localhost:5001/bootstrap-tyk-pre-delete:$bsVers -f ./.container/image/bootstrap-pre-delete/Dockerfile ./bin & docker push localhost:5001/bootstrap-tyk-pre-delete:$bsVers
```

The "hack" folder comes with a job (job.yaml) that can be applied directly together
with the role.yaml (which contains the ServiceAccount and ClusterRoleBinding) 
into a namespace running an empty tyk stack for debugging/development purposes.
