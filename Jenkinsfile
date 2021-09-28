pipeline {
    agent any
    
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

        TARGET_REPO = "https://oauth2:hvE2MzrZzH6wnFyEDcjS@gitlab.pharmalink.id/rnd/backend-pipeline-security"
        TFHOG_DIR = '/usr/local/trufflehog'
        GOLANGCI_DIR = '/usr/local/golangci-lint'
    }
    
    options {
        skipDefaultCheckout(true)
    }
    stages {
        stage('Checkout SCM') {
            steps {
                echo '> Checking out the source control ...'
                checkout scm
            }
        }
        stage('Installing Library'){
            steps{
                echo '[*] Installing TruffleHog ...'
                sh 'pip3 install trufflehog'
                sh 'locate trufflehog && locate truffleHog'
            }
        }
        stage('Declarative Variable'){
            steps{
                script{
                    WORKSPACE = sh (
                        script: "pwd",
                        returnStdout: true
                    )
                }
            }
        }
        stage('GoLangCI-Lint'){
            steps{
                script{
                    try{
                        echo "[*] Running Linter ErrCheck"
                        sh "${GOLANGCI_DIR}/bin/golangci-lint run --disable-all -E errcheck"
                    }catch(err){
                        echo "${err}"               
                    }
                }
            }
        }
        stage('TruffleHog'){
            steps{
                script{
                    try{
                        echo "[*] Running truffleHog ..."
                        sh "${TFHOG_DIR}/bin/trufflehog --regex --json --max_depth 1 --rules ${TFHOG_DIR}/rules.json ${TARGET_REPO} > tfhog.json"
                    }
                    catch(err) {
                        
                    }
                    sh 'ls -la'
                    sh 'cat tfhog.json'
                    echo "[*] Scanning done ..."
                }
            }
        }
        stage('Create Reporting'){
            steps{
                echo '[*] Create report ...'
                script {
                    def now = new Date()
                    env.FILENAME = now.format("dd-MM-YYYY_HH:mm:ss", TimeZone.getTimeZone('GMT+7'))
                }
                sh 'python3 ${TFHOG_DIR}/convert.py ${WORKSPACE}'
            }
        }        
    }
    post{
        success {
			discordSend link: env.BUILD_URL, result: currentBuild.currentResult, title: "${env.JOB_NAME} #${env.BUILD_NUMBER}", webhookURL: "${env.DISCORD_WEBHOOK_URL}"
			sh "exit 0"
		}

		regression {
			discordSend link: env.BUILD_URL, result: currentBuild.currentResult, title: "${env.JOB_NAME} #${env.BUILD_NUMBER}", webhookURL: "${env.DISCORD_WEBHOOK_URL}"
			sh "exit 1"
		}
	}
}
