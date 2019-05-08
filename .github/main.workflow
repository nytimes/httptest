workflow "push" {
  on = "push"
  resolves = ["docker-publish-latest"]
}

action "docker-auth" {
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "docker-build" {
  uses = "actions/docker/cli@master"
  args = "build -t blupig/httptest:dev --build-arg BUILD_BRANCH=${GITHUB_REF} --build-arg BUILD_COMMIT=${GITHUB_SHA} ."
}

action "docker-publish-dev" {
  needs = ["docker-auth", "docker-build"]
  uses = "actions/docker/cli@master"
  args = ["push blupig/httptest"]
}

action "branch-master" {
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "docker-tag-latest" {
  needs = ["branch-master", "docker-publish-dev"]
  uses = "actions/docker/cli@master"
  args = "tag blupig/httptest:dev blupig/httptest"
}

action "docker-publish-latest" {
  needs = ["docker-tag-latest"]
  uses = "actions/docker/cli@master"
  args = ["push blupig/httptest"]
}
