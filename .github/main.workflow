workflow "push" {
  on = "push"
  resolves = ["docker-publish-latest"]
}

action "not-branch-del" {
  uses = "actions/bin/filter@master"
  args = "not deleted"
}

action "docker-auth" {
  needs = ["not-branch-del"]
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "docker-build" {
  needs = ["not-branch-del"]
  uses = "actions/docker/cli@master"
  args = "build -t blupig/httptest:dev --build-arg BUILD_BRANCH=${GITHUB_REF} --build-arg BUILD_COMMIT=${GITHUB_SHA} ."
}

action "docker-publish-dev" {
  needs = ["docker-auth", "docker-build"]
  uses = "actions/docker/cli@master"
  args = ["push blupig/httptest"]
}

action "branch-master" {
  needs = ["docker-publish-dev"]
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "docker-tag-latest" {
  needs = ["branch-master"]
  uses = "actions/docker/cli@master"
  args = "tag blupig/httptest:dev blupig/httptest"
}

action "docker-publish-latest" {
  needs = ["docker-tag-latest"]
  uses = "actions/docker/cli@master"
  args = ["push blupig/httptest"]
}
