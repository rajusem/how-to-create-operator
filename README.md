# What is Kubernetes Operator 
A Kubernetes Operator is a method of packaging, deploying, and managing a Kubernetes application. It extends the Kubernetes API to create, configure, and manage instances of complex stateful applications.

Operators aim to automate routine operational tasks that are typically performed by human operators and administrators.


### Key characteristics/Components of a Kubernetes Operator

**Custom Resource Definitions (CRDs)**: Operators use Custom Resource Definitions to extend the Kubernetes API and define custom objects that represent instances of the application. These CRDs define the desired state of the application.

**Controller Logic**: The Operator includes a controller, which is a custom Kubernetes controller responsible for watching the custom resources and taking actions to ensure that the observed state matches the desired state.

**Automation**: The primary goal of an Operator is to automate the management of complex applications. This includes tasks such as deployment, scaling, updates, and failure recovery.

**Domain-Specific Knowledge**: Operators embed domain-specific knowledge about the application they manage. They understand the intricacies of the application, enabling them to make intelligent decisions and automate complex tasks.

**Lifecycle Management**: Operators handle the entire lifecycle of an application, from deployment to scaling, updates, and eventual decommissioning.

**Self-Healing**: Operators are designed to monitor the state of the application and take corrective actions in case of failures or deviations from the desired state.

**Example**: An Operator could be created for a database system like PostgreSQL. The Operator would define a custom resource for PostgreSQL instances, and the controller logic would ensure that the specified PostgreSQL instances are provisioned, scaled, and updated as per the desired configuration.

## Operator Capability Levels
Operators come in different maturity levels in regards to their lifecycle management capabilities for the application or workload they deliver. The capability models aims to provide guidance in terminology to express what features users can expect from an operator.
Refer [Capability Levels](https://sdk.operatorframework.io/docs/overview/operator-capabilities/) for more info. 

# How to create/run operator

## Prerequisites:

Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and have access to a cluster.
Install [Operator SDK](https://sdk.operatorframework.io/docs/installation/).

## Steps to create operator 
### Step 1: Create a New Operator Project and initialize operator
```
mkdir my-operator
cd my-operator
```
Initializes a new Operator project with a basic directory structure.
```
operator-sdk init --domain=example.com --repo=github.com/example/my-operator
```

### Step 2: Create a New API for Your Operator
```
operator-sdk create api --group=app --version=v1alpha1 --kind=MyApp
```
This generates the necessary code for your custom resource (CR) and its controller.
below files gets generated 

#### Couple of files we interested are as follow. 
- `api/v1alpha1/myapp_types.go` - It holds schema for the myapps API
- `controllers/myapp_controller.go` - Holds logic for the Reconcile. This will be called when kubernetes resource with `kind: MyApp` gets created (Object of CRD).   
- `config/samples/app_v1alpha1_myapp.yaml` - Holds sample object for the CRD

### Step 3: Update Schema file and controller logic.
Here we will try to automate namespace creation as well as deploying cronjob based on given configruation. 

- For the same, replace `api/v1alpha1/myapp_types.go` with [myapp_types.go](my-operator/api/v1alpha1/myapp_types.go)  where we have defined the `MyAppSpec` struct, which will be used to get information on namespace and cronjob. 


- Replace file `controllers/myapp_controller.go` with [myapp_controller.go](my-operator/controllers/myapp_controller.go) to update reconciliation loop logic to create namespace and cronjob.

- Run below command to generate Kubernetes manifests or resource definition files, including those for Custom Resource Definitions (CRDs)
    ```
    make generate
    make manifests
    ```

- Replace file `config/samples/app_v1alpha1_myapp.yaml` with [app_v1alpha1_myapp.yaml](my-operator/config/samples/app_v1alpha1_myapp.yaml) to create sample definition of the CRD Object.  

## Step to deploy operator and usage 

### Option 1: Run operator from console. 

#### Step 1.1 Login into cluster. 
Make sure you have cluster running and you can able to access it from your console. 

#### Step 1.2 Install the CRDs into the cluster:
```sh
make install
or
kubectl apply -f config/crd/bases/
```

#### Step 1.3 Once installed you can verify it on cluster by running following command. 
```sh
kubectl get crd | grep myapp
kubectl describe crd myapps.app.example.com
```

#### Step 1.4 Run your controller:
```sh
make run
```
This will run in the foreground, so switch to a new terminal

#### Step 1.5 Install Instances of Custom Resources:
```sh
kubectl apply -f config/samples/app_v1alpha1_myapp.yaml
```
Note: You can update the image you wanted to deploy, currently its pointed to the simple app describe [here](sample-python-app/README.md). 

#### Step 1.6 Operator should create namespace and cronjob under given namespace in 4.5 step. You can verify the same in your cluster.  

#### Step 1.7 Delete the created cronjob from the cluster and you can see after sometime reconsiler will create it again.(See the logs of `make run` console window) 


Refer to the official [Operator SDK documentation](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/) for more detailed information and advanced features.
