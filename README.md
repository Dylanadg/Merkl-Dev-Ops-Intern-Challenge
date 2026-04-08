# Merkl DevOps Intern Challenge

by Dylan Adghar

## AI Usage

**Model**: Claude Sonnet 4.6

**Tools**: claude.ai chat interface

**Why**:  I used Claude mainly for Part 2 (Operator SDK) to better understand the operator-sdk workflow and clarify Go syntax I was not fully comfortable with. I also used it to help structure and improve the clarity of this README.

## Part 1 / Kubernetes core concepts

Use minikube as my local Kubernetes cluster.

### Task 1.1 / Namespace & basic workload

Creates the `intern-assessment` namespace, an `nginx:1.25` Deployment with 2 replicas, and a NodePort Service on port 80.

```bash
kubectl apply -f task1-workload.yaml

namespace/intern-assessment created
deployment.apps/nginx-deployment created
service/nginx-service created
```

Verification commands used to check the created resources:

```bash
kubectl get namespace intern-assessment
NAME                STATUS   AGE
intern-assessment   Active   22s

kubectl get deployment -n intern-assessment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   2/2     2            2           44s

kubectl get pods -n intern-assessment
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-6f9664446b-98n42   1/1     Running   0          53s
nginx-deployment-6f9664446b-f5r6g   1/1     Running   0          53s
```

Access the service:

```bash
minikube service nginx-service -n intern-assessment
┌───────────────────┬───────────────┬─────────────┬───────────────────────────┐
│     NAMESPACE     │     NAME      │ TARGET PORT │            URL            │
├───────────────────┼───────────────┼─────────────┼───────────────────────────┤
│ intern-assessment │ nginx-service │ 80          │ http://192.168.49.2:30080 │
└───────────────────┴───────────────┴─────────────┴───────────────────────────┘
```

---

### Task 1.2 / ConfigMap & environment variable injection

Creates a ConfigMap `app-config` with `APP_ENV=staging` and injects it into the Deployment as an environment variable.

```bash
kubectl apply -f task2-configmap.yaml

configmap/app-config created
deployment.apps/nginx-deployment configured
```

Verification:

```bash
kubectl get configmap app-config -n intern-assessment
NAME         DATA   AGE
app-config   1      57s

kubectl describe configmap app-config -n intern-assessment
Name:         app-config
Namespace:    intern-assessment
Data
====
APP_ENV:
----
staging

kubectl exec -n intern-assessment nginx-deployment-96c78cb99-58csb -- env | grep APP_ENV
APP_ENV=staging
```

---

### Task 1.3 / Resource requests & limits

Updates the Deployment with CPU and memory requests/limits on the nginx container.

```bash
kubectl apply -f task3-resources.yaml

deployment.apps/nginx-deployment configured
```

Verification:

```bash
kubectl describe pod -n intern-assessment
    Limits:
      cpu:     100m
      memory:  128Mi
    Requests:
      cpu:     50m
      memory:  64Mi
```

---

### Task 1.4 / Liveness & readiness probes

Adds HTTP health checks to the nginx container.

```bash
kubectl apply -f task4-probes.yaml

deployment.apps/nginx-deployment configured
```

Verification:

```bash
kubectl describe pod -n intern-assessment
    Liveness:   http-get http://:80/ delay=15s timeout=1s period=20s #success=1 #failure=3
    Readiness:  http-get http://:80/ delay=5s timeout=1s period=10s #success=1 #failure=3
```

## Part 2 / Operator SDK
 
### Task 2.1 / Scaffold
 
```bash
cd hello-operator
operator-sdk init --domain intern.dev --repo github.com/Dylanadg/hello-operator
operator-sdk create api --group apps --version v1alpha1 --kind HelloApp --resource --controller
```
 
### Task 2.2 / CRD spec
 
After editing `api/v1alpha1/helloapp_types.go`:
 
```bash
make generate && make manifests
```
 
### Task 2.3 / Reconcile loop
 
Implemented in `internal/controller/helloapp_controller.go`.
 
### Task 2.4 / Deploy & verify
 
Install the CRD and run the controller:
 
```bash
make install
make run
```
 
Apply the sample HelloApp in a second terminal:
 
```bash
kubectl apply -f config/samples/helloapp-sample.yaml
 
helloapp.apps.intern.dev/hello-sample created
```
 
Verification:
 
```bash
kubectl get deployment -n intern-assessment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
hello-sample       2/2     2            2           10s
nginx-deployment   2/2     2            2           165m
 
kubectl get pods -n intern-assessment
NAME                                READY   STATUS    RESTARTS   AGE
hello-sample-5bcfb88cc5-2b6st       1/1     Running   0          30s
hello-sample-5bcfb88cc5-zl5nv       1/1     Running   0          30s
nginx-deployment-656fcf7d54-88l9q   1/1     Running   0          140m
nginx-deployment-656fcf7d54-xtqdn   1/1     Running   0          140m
 
kubectl logs hello-sample-5bcfb88cc5-2b6st -n intern-assessment
Hello from my first operator!
 
kubectl get helloapp hello-sample -n intern-assessment -o yaml | grep availableReplicas
  availableReplicas: 2
```
