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
        stage('GoLangCI-Lint'){
            steps{
                script{
                try{
                    echo "[*] Running Linter Gosec ..."
                    sh "golangci-lint run --disable-all -E gosec"
                }catch(err){}
                try{
                    echo "[*] Running Linter Deadcode ..."
                    sh "golangci-lint run --disable-all -E deadcode"
                }catch(err){}
                try{
                    echo "[*] Running Linter StaticCheck"
                    sh "golangci-lint run --disable-all -E staticcheck"
                }catch(err){}
                try{
                    echo "[*] Running Linter Unused"
                    sh "golangci-lint run --disable-all -E unused"
                }catch(err){}
                try{
                    echo "[*] Running Linter ErrCheck"
                    sh "golangci-lint run --disable-all -E errcheck"
                }catch(err){
                    echo "${err}"               }
                }
            }
        }
        stage('TruffleHog'){
            steps{
                script{
                    try{
                        echo "[*] Running truffleHog ..."
                        sh "trufflehog --regex --json --max_depth 1 --rules ${TARGET_DIR}/rules.json ${TARGET_REPO} > ${TARGET_DIR}/tfhog.result"
                    }
                    catch(err) {
                        
                    }
                    echo "[*] Scanning done ..."
                }
                
                echo "[*] Checking scan result ..."
                script{
                    TFHOG_RESULT = sh (
                        script: "jq . ${TARGET_DIR}/tfhog.result",
                        returnStdout: true
                    )

                }
                echo "${TFHOG_RESULT}"

                script{
                    if ( TFHOG_RESULT ){
                        echo "[*] Credential leaked ..."
                    }
                    else {
                        echo "[*] No credential leaked ..."
                    }
                }

                echo "[*] Remove existing scan result file (tfhog.result) ..."
                sh "rm ${TARGET_DIR}/tfhog.result"
            }
        }        
    }
}
