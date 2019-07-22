def label = "${env.NODE}"
def deployJob = "${env.DEPLOY_JOB}"
def DOCKER_REGISTRY = ""
def GIT_COMMIT
def GIT_COMMIT_SHORT

timestamps {
  node(label) {

    def app

    stage("Checkout") {
      checkout scm

      GIT_COMMIT = sh(returnStdout: true, script: "git rev-parse HEAD").trim()
      GIT_COMMIT_SHORT = GIT_COMMIT.substring(0, 8)
      println("Git revision is: ${GIT_COMMIT}, short revision is: ${GIT_COMMIT_SHORT}.")
    }

    def TAG = "${env.BUILD_TIMESTAMP}-${GIT_COMMIT_SHORT}"

    docker.withTool('Docker') {
      docker.withRegistry("${DOCKER_REGISTRY}", "${env.DOCKER_CREDS}") {

        stage("Build image") {
          app = docker.build("${env.DOCKER_IMAGE}", "--build-arg DOCKER_REGISTRY=${DOCKER_REGISTRY} .")
        }

        stage("Push image") {
          app.push("${TAG}")
        }

        currentBuild.description = "Branch: ${env.BRANCH}\nImage: ${env.DOCKER_IMAGE}:${TAG}"
      }
    }

    stage("Send e-mails") {
      if (currentBuild.currentResult == "SUCCESS") {
        emailext subject: '$DEFAULT_SUBJECT',
                body: '$DEFAULT_CONTENT',
                recipientProviders: [
                        [$class: 'RequesterRecipientProvider']
                ],
                replyTo: '$DEFAULT_REPLYTO',
                to: '$DEFAULT_RECIPIENTS'
      } else {
        emailext subject: '$DEFAULT_SUBJECT',
                body: '$DEFAULT_CONTENT',
                recipientProviders: [
                        [$class: 'RequesterRecipientProvider'],
                        [$class: 'DevelopersRecipientProvider']
                ],
                replyTo: '$DEFAULT_REPLYTO',
                to: '$DEFAULT_RECIPIENTS'
      }
    }

    stage("Deploy") {
      if (env.DEPLOY) {
        build job: "${deployJob}", parameters: [
          string(name: 'ENV', value: "${ENV}"),
          string(name: 'COMPONENT', value: "toggly-core"),
          string(name: 'IMAGE', value: "${DOCKER_REGISTRY}/${env.DOCKER_IMAGE}:${TAG}")]
      }
    }

    stage("Delete workspace") {
      cleanWs()
    }

  }
}
