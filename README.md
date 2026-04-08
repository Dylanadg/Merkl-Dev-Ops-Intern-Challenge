# Merkl DevOps Intern Challenge

by Dylan Adghar

## AI Usage

**Model**: Claude Sonnet 4.6

**Tools**: claude.ai chat interface

**Why**: I used the model mainly to double‑check the syntax of some YAML elements and refresh a few Kubernetes concepts. It helped me recall best practices and ensure everything was implemented correctly

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


