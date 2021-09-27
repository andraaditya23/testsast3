pipeline {
    agent any
    tools {
        go 'go-1.17'
    }
    parameters {
        credentials credentialType: 'com.dabsquared.gitlabjenkins.connection.GitLabApiTokenImpl', defaultValue: '1', description: 'gitlab_token', name: 'GITLAB_TOKEN', required: false
    }
    environment {
        GO111MODULE = "on"

        PROJECT_ID = "pharos-main"
        
        NAME = "backend-pipeline-security"
        ORG = "pharos"

        DOCKER_REGISTRY = "gcr.io"
        DOCKER_REGISTRY_URL = "https://gcr.io"
        DOCKER_REGISTRY_PROJECT_URL = "${DOCKER_REGISTRY}/${PROJECT_ID}"
        DOCKER_IMAGE_URL = "${DOCKER_REGISTRY_PROJECT_URL}/${NAME}"
        
        PIPELINE_BOT_EMAIL = "pharos.bot@gmail.com"
        PIPELINE_BOT_NAME = "Pharmalink Pipeline Bot"

        DISCORD_WEBHOOK_URL = "https://discord.com/api/webhooks/877591443986870313/0ALWAO9W7cSgo4LytxSYUJtSXDoRKm9dnQGp-fHWtKfcsS4YCgC7kUpQPApemhZBjOnf"

        TARGET_REPO = "https://${GITLAB_TOKEN}@gitlab.pharmalink.id/rnd/backend-pipeline-security"
        TARGET_DIR = "/var/jenkins_home/workspace/gitlab-scanner"

    }
    
    options {
        skipDefaultCheckout(true)
    }
    stages {
        stage('Checkout SCM') {
            steps {
                echo '> Checking out the source control ...'
                checkout scm
                echo "${GITLAB_TOKEN}"
            }
        }        
    }
}
