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
        TARGET_DIR = '/usr/local/trufflehog/'
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
        stage('GoLangCI-Lint'){
            steps{
                script{
                    sh 'alias go="/usr/local/go/bin/go"'
                try{
                    echo "[*] Running Linter ErrCheck"
                    sh "golangci-lint run --disable-all -E errcheck --out-format json --new-from-rev=HEAD~ > ${TARGET_DIR}/rawJson/errcheck.json"
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
                        sh "${TARGET_DIR}/bin/trufflehog --regex --json --max_depth 1 --rules ${TARGET_DIR}/rules.json ${TARGET_REPO} > ${TARGET_DIR}/rawJson/tfhog.json"
                    }
                    catch(err) {
                        
                    }
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
                sh 'python3 ${TARGET_DIR}/convert.py > ${TARGET_DIR}/beautyJson/${FILENAME}'
                script{
                    env.FILE_CONTENT = sh (
                        script: "cat ${TARGET_DIR}/beautyJson/${FILENAME}"
                    )
                }
                echo '${FILE_CONTENT}'
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
