# Campaigns design doc

Why are [campaigns](../user/campaigns/index.md) designed the way they are?

Goals:

- Declarative (not imperative) API
- 

## Inspired by Kubernetes

We've found Kubernetes to be a good source of inspiration for the campaigns API design, because both Kubernetes and campaigns help you manage a distributed system. Here's how Kubernetes concepts (for a Kubernetes [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)) map to Sourcegraph campaigns concepts:

<table>
  <tr>
    <td/>
    <th>Kubernetes <a href="https://kubernetes.io/docs/concepts/workloads/controllers/deployment/">Deployment</a></th>
    <th>Sourcegraph campaign</th>
  </tr>
  <tr>
    <th>What underlying thing does this API manage?</th>
    <td>Pods running on many (possibly unreliable) nodes</td>
    <td>Branches and changesets on many repositories that can be rate-limited and externally modified (and our authorization can change)</td>
  </tr>
  <tr>
    <th>Spec YAML</th>
    <td>
      <pre><code><em># File: foo.Deployment.yaml</em>
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:<div style="padding:0.5rem;margin:0 -0.5rem;background-color:rgba(0,0,255,0.25)"><strong>  # Evaluate this to enumerate instances of...</strong>
  replicas: 2</div>
<div style="padding:0.5rem;margin:0 -0.5rem;background-color:rgba(255,0,255,0.25)"><strong>  # ...this template.</strong>
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80</div></pre></code>
    </td>
    <td>
      <pre><code><em># File: hello-world.campaign.yaml</em>
name: hello-world
description: Add Hello World to READMEs

<div style="padding:0.5rem;margin:0 -0.5rem;background-color:rgba(0,0,255,0.25)"><strong># Evaluate this to enumerate instances of...</strong>
on:
  - repositoriesMatchingQuery: file:README.md

steps:
  - run: echo Hello | tee -a $(find -name '*.md')
    container: alpine:3
</div>
<div style="padding:0.5rem;margin:0 -0.5rem;background-color:rgba(255,0,255,0.25)"><strong># ...this template.</strong>
changesetTemplate:
  title: Hello World
  body: My first campaign!
  branch: hello-world
  commit:
    message: Append Hello to .md files
  published: false</pre></code>
    </td>
  </tr>
  <tr>
    <th>How desired state is computed</th>
    <td>
      <ol>
        <li>Evaluate <code>replicas</code>, etc. (blue) to determine pod count and other template inputs</li>
        <li>Instantiate <code>template</code> (pink) once for each pod to produce PodSpecs</li>
      </ol>
    </td>
    <td>
      <ol>
        <li>Evaluate <code>on</code>, <code>steps</code> (blue) to determine list of patches</li>
        <li>Instantiate <code>changesetTemplate</code> (purple) once for each patch to produce ChangesetSpecs
</li>
      </ol>
    </td>
  </tr>
  <tr>
    <th>Desired state consists of...</th>
    <td>
      <ul>
        <li>DeploymentSpec file (the YAML above)</li>
        <li>List of PodSpecs (template instantiations)</li>
      </ul>
    </td>
    <td>
      <ul>
        <li>CampaignSpec file (the YAML above)</li>
        <li>List of ChangesetSpecs (template instantiations)</li>
      </ul>
    </td>
  </tr>
  <tr>
    <th>Where is the desired state computed?</th>
    <td>The deployment controller (part of the Kubernetes cluster) consults the DeploymentSpec and continuously computes the desired state.</td>
    <td>
      <p>The <a href="https://github.com/sourcegraph/src-cli">Sourcegraph CLI</a> (running on your local machine, not on the Sourcegraph server) consults the campaign spec and computes the desired state when you invoke <code>src campaign apply</code>.</p>
      <p><strong>Difference vs. Kubernetes:</strong> A campaign's desired state is computed locally, not on the server. It requires executing arbitrary commands, which is not yet supported by the Sourcegraph server. See campaigns known issue "<a href="../user/campaigns/index.md#server-execution">Campaign steps are run locally...</a>".</p>
    </td>
  </tr>
  <tr>
    <th>Reconciling desired state vs. actual state</th>
    <td>The "deployment controller" reconciles the resulting PodSpecs against the current actual PodSpecs (and does smart things like rolling deploy).</td>
    <td>The "campaign controller" (i.e., our backend) reconciles the resulting ChangesetSpecs against the current actual changesets (and does smart things like gradual roll-out/publishing and auto-merging when checks pass).</td>
  </tr>
</table>

asdf


- [Object Management](https://kubernetes.io/docs/concepts/overview/working-with-objects/object-management/)
- [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)

- [Architecture](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/architecture/architecture.md)
- [Kubernetes Design and Development Explained](https://thenewstack.io/kubernetes-design-and-development-explained/)