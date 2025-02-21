---
sources:
  stable:
    kind: jenkins
    depends_on:
      - weekly
    name: 'Get Latest Jenkins stable version and depends on {{ pipeline "sources.Weekly.kind" }}'
  weekly:
    kind: jenkins
    name: Get Latest Jenkins weekly version
conditions:
  stabledockerImage:
    name: 'Is docker image jenkins/jenkins:{{ pipeline "Sources.stable.kind" }} published?'
    kind: dockerImage
    sourceID: stable
    spec:
      image: jenkins/jenkins
      tag: '{{ source "weekly" }}-jdk11'
targets:
  imageTag:
    name: 'Update jenkins/jenkins docker'
    kind: yaml
    sourceID: stable
    spec:
      file: "charts/jenkins/values.yaml"
      key: "jenkins.controller.tag"
    scm:
      git:
        url: "git@github.com:olblak/charts.git"
        branch: master
        user: olblak
        email: me@olblak.com
