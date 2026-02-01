# Setting up ACK S3 Controller on a Non-EKS Cluster

This guide explains how to install the AWS Controllers for Kubernetes (ACK) S3 controller on a standard Kubernetes cluster (e.g., K3s, Minikube, or bare metal) and configure it with AWS credentials.

## Prerequisites
- A running Kubernetes cluster.
- `helm` installed.
- AWS Credentials (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`) with permissions to manage S3 buckets.

## 1. Install the ACK S3 Controller

We will use Helm to install the controller.

```bash
# Logout from ECR Public if needed (sometimes required to pull public images without auth issues)
docker logout public.ecr.aws

# Login to Helm registry (optional, usually public access works)
aws ecr-public get-login-password --region us-east-1 | helm registry login --username AWS --password-stdin public.ecr.aws

# Install the S3 Chart
export SERVICE=s3
export RELEASE_VERSION=$(curl -sL https://api.github.com/repos/aws-controllers-k8s/${SERVICE}-controller/releases/latest | jq -r '.tag_name | ltrimstr("v")')
export AWS_REGION=ap-south-1

helm install --create-namespace -n ack-system "ack-$SERVICE-controller" \
  oci://public.ecr.aws/aws-controllers-k8s/$SERVICE-chart \
  --version "$RELEASE_VERSION" \
  --set aws.region="$AWS_REGION"
```

## 2. Configure AWS Credentials

Since this is a **non-EKS** cluster, we cannot use IRSA (IAM Roles for Service Accounts) directly without setting up an OIDC provider. The simplest method for dev/non-EKS is utilizing a **Shared Credentials File** mounted as a Secret.

### Step 2.1: Create Credentials File

Create a file named `creds.conf` with your AWS credentials:

```ini
[default]
aws_access_key_id = YOUR_ACCESS_KEY_ID
aws_secret_access_key = YOUR_SECRET_ACCESS_KEY
region = ap-south-1
```

### Step 2.2: Create Kubernetes Secret

```bash
kubectl create secret generic aws-creds -n ack-system --from-file=credentials=./creds.conf
```

### Step 2.3: Configure Controller to use Secret

Upgrade the Helm chart to inject these credentials as environment variables.

```bash
helm upgrade --install -n ack-system "ack-$SERVICE-controller" \
  oci://public.ecr.aws/aws-controllers-k8s/$SERVICE-chart \
  --version "$RELEASE_VERSION" \
  --set aws.region="$AWS_REGION" \
  --set aws.credentials.secretName=aws-creds \
  --set aws.credentials.secretKey=credentials \
  --set aws.credentials.profile=default
```

## 3. Verify Installation

Check if the controller pod is running:

```bash
kubectl get pods -n ack-system
```

## 4. Create an S3 Bucket

Apply the sample manifest:

```bash
kubectl apply -f s3-bucket.yaml
```

Verify creation:

```bash
kubectl get buckets
```

## 5. Setup Kro (Kubernetes Resource Orchestrator)

We use Kro to define abstract application resources.

### Step 5.1: Install Kro

```bash
helm install kro kro \
  --repo https://kro-run.github.io/kro \
  --namespace kro \
  --create-namespace
```

### Step 5.2: Define the `Application` Resource

Apply the `ResourceGraphDefinition` which tells Kro to create a new `Application` CRD.

```bash
kubectl apply -f kro-spec.yaml
```

Wait a few seconds for Kro to generate the CRD. You can check if it exists:

```bash
kubectl get crd applications.kro.run
```

### Step 5.3: Create an Application Instance

Now you can create your application, which will automatically trigger the creation of the underlying S3 bucket via ACK.

```bash
kubectl apply -f s3.yaml
```

## 6. Troubleshooting

### "applications.kro.run not found"
If the CRD is not found, the Kro controller might have failed to generate it. Check the logs:
```bash
# Get the Kro pod name
kubectl get pods -n kro

# Check logs
kubectl logs -n kro <kro-pod-name>
```

### "no matches for kind..."
Ensure you have applied the `kro-spec.yaml` first and waited for the CRD to be established.


