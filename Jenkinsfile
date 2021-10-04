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

        TFHOG_DIR = '/usr/local/trufflehog'
        GOLANGCI_DIR = '/usr/local/golangci-lint'
        DEPENDENCY_CHECK_DIR = '/usr/local/dependency-check/6.3.1'
        scannerHome = tool 'SonarQube';
    }
    
    options {
        skipDefaultCheckout(true)
        ansiColor('xterm')
    }
    stages {
        stage('Checkout SCM') {
            steps {
                echo '> Checking out the source control .... '
                script{
                    def GIT = checkout scm
                    env.TARGET_REPO = GIT.GIT_URL
                }
            }
        }
        stage('Installing Library'){
            steps{
                echo "[*] Install Git"
                sh '{ pip3 install gitpython; } 2>/dev/null'
            }
        }
        stage('Declarative Variable'){
            steps{
                script{
                    WORKSPACE = sh (
                        script: "{ pwd; } 2>/dev/null",
                        returnStdout: true
                    )
                    AUTHOR = sh (
                        script: "{ git log -1 --pretty=format:'%an <%ae>'; } 2>/dev/null",
                        returnStdout: true
                    )
                }
            }
        }
        stage('GoLangCI-Lint'){
            steps{
                script{
                    sh "{ export PATH=$PATH:/usr/local/go/bin; } 2>/dev/null"
                    try{
                        echo "[*] Running Linter"
                        sh "{ ${GOLANGCI_DIR}/bin/golangci-lint run -c./.golangci.yaml --out-format json --new-from-rev=HEAD~ > golangci-report.json; } 2>/dev/null"
                    }catch(err){}
                }
            }
        }
        stage('TruffleHog'){
            steps{
                script{
                    try{
                        echo "[*] Running truffleHog ...."
                        withCredentials([gitUsernamePassword(credentialsId: 'gitlab-pipeline-bot', gitToolName: 'git-tool')]) {
                        sh "{ ${TFHOG_DIR}/bin/trufflehog --regex --json --max_depth 1 --rules ${TFHOG_DIR}/rules.json ${TARGET_REPO} > tfhog-report.json; } 2>/dev/null"
}

                    }
                    catch(err) {
                        
                    }
                }
            }
        }
        stage('SonarQube Analysis') {
            steps{
                withSonarQubeEnv() {
                    sh "${scannerHome}/bin/sonar-scanner"
                }
            }
        }
        stage('Create Reporting'){
            steps{
                echo '[*] Create report ....'
                script {
                    def now = new Date()
                    env.REPORT_TIME = now.format("dd-MM-YYYY_HH:mm:ss", TimeZone.getTimeZone('GMT+7'))

                    sh '{ python3 ${TFHOG_DIR}/convert.py ${WORKSPACE} > ${WORKSPACE}/${REPORT_TIME}; } 2>/dev/null'
                    sh '{ cat ${REPORT_TIME}; } 2>/dev/null'
                    
                    ISSUE_COUNT = sh(
                        script: "{ grep -o 'Found IssuE' ${REPORT_TIME} | wc -l; } 2>/dev/null",
                        returnStdout: true
                    ).trim().toString()
                    echo "${ISSUE_COUNT}"
                }               
                echo '[*] Remove report file ...'
                sh '{ rm ${REPORT_TIME}; } 2>/dev/null'
            }
        }        
    }
    post{
        success {
            script{
                if(ISSUE_COUNT != '0'){
                    discordSend link: "${env.BUILD_URL}console", 
                    result: currentBuild.currentResult, 
                    title: "${env.JOB_NAME} #${env.BUILD_NUMBER}\n>> click for details ...", 
                    webhookURL: "${env.DISCORD_WEBHOOK_URL}", 
                    description:"```yaml\nTimestamp  : ${REPORT_TIME}\nAuthor     : ${AUTHOR}\nIssue      : ${ISSUE_COUNT}\n```SonarQube  : [here](http://34.126.163.106:9000/dashboard?id=research-test)"
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
