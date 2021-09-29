pipeline {
    agent any
    tools {
    go 'go'
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
                echo "[*] Install Git"
                sh 'pip3 install gitpython'
            }
        }
        stage('Declarative Variable'){
            steps{
                script{
                    WORKSPACE = sh (
                        script: "pwd",
                        returnStdout: true
                    )
                    AUTHOR = sh (
                        script: "git log -1 --pretty=format:'%an <%ae>'",
                        returnStdout: true
                    )
                }
            }
        }
        stage('GoLangCI-Lint'){
            steps{
                script{
                    sh "export PATH=$PATH:/usr/local/go/bin"
                    try{
                        echo "[*] Running Linter"
                        sh "${GOLANGCI_DIR}/bin/golangci-lint run -c./.golangci.yaml --out-format json --new-from-rev=HEAD~ > errcheck.json"
                    }catch(err){}
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
                }
            }
        }
        stage('Create Reporting'){
            steps{
                echo '[*] Create report ...'
                script {
                    def now = new Date()
                    env.REPORT_TIME = now.format("dd-MM-YYYY_HH:mm:ss", TimeZone.getTimeZone('GMT+7'))
                }
                sh 'python3 ${TFHOG_DIR}/convert.py ${WORKSPACE} > ${WORKSPACE}/${REPORT_TIME}'
                sh 'cat ${REPORT_TIME}'
                script{
                    ISSUE_EXIST = sh(
                        script: "grep -o 'Issue #' ${REPORT_TIME}",
                        returnStdout: true
                    )
                }
                echo '[*] Remove report file ...'
                sh 'rm ${REPORT_TIME}'
            }
        }        
    }
    post{
        success {
            script{
                if(ISSUE_EXIST){
                    discordSend link: "${env.BUILD_URL}console", 
                                result: currentBuild.currentResult, 
                                title: "${env.JOB_NAME} #${env.BUILD_NUMBER}", 
                                webhookURL: "${env.DISCORD_WEBHOOK_URL}", 
                                description:"```fix\nalert ```"
                    sh "exit 0"
                }
            }
		}

		regression {
			discordSend link: "${env.BUILD_URL}console", result: currentBuild.currentResult, title: "${env.JOB_NAME}\n#${env.BUILD_NUMBER}", webhookURL: "${env.DISCORD_WEBHOOK_URL}"
			sh "exit 1"
		}
	}
}
