source:
  kind: maven
  postfix: "-jdk11"
  spec:
    owner: "maven"
    url: "repo.jenkins-ci.org"
    repository: "releases"
    groupID: "org.jenkins-ci.main"
    artifactID: "jenkins-war"
conditions:
  docker:
    name: "Docker Image Published on Registry"
    kind: dockerImage
    spec:
      image: "jenkins/jenkins"
targets:
  imageTag:
    name: "Docker Image"
    kind: yaml
    prefix: "test-"
    postfix: "-jdk13"
    spec:
      file: "charts/jenkins/values.yaml"
      key: "jenkins.master.imageTag"
    scm:
      github:
        user: "updatecli"
        email: "updatecli@olblak.com"
        owner: "olblak"
        repository: "charts"
        token: {{ requiredEnv "GITHUB_TOKEN" }}
        username: "olblak"
        branch: "master"