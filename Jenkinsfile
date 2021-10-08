pipeline {
    agent any
    tools {
        go 'go1.17'
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
        GITLEAKS_DIR = '/usr/local/gitleaks'

        GCS_BUCKET = 'pharmalink-id-build-logs'
    }
    
    options {
        skipDefaultCheckout(true)
        ansiColor('xterm')
    }
    stages {
        stage('Checkout SCM') {
            steps {
                echo '> Checking out the source control ....'
                script{
                    def GIT = checkout scm
                    env.TARGET_REPO = GIT.GIT_URL
                }
            }
        }
        stage('Installing Library'){
            steps{
                echo "[*] Install Git .."
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
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE'){
                    sh "{ export PATH=$PATH:/usr/local/go/bin; } 2>/dev/null"
                    echo "[*] Running Linter"
                    sh "{ ${GOLANGCI_DIR}/bin/golangci-lint run -c./.golangci.yaml --out-format json --new-from-rev=HEAD~ > golangci-report.json; } 2>/dev/null"
                }               
            }
        }
        stage('TruffleHog'){
            steps{
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE'){
                    echo "[*] Running truffleHog ...."
                    withCredentials([gitUsernamePassword(credentialsId:'gitlab-pipeline-bot',gitToolName: 'git-tool')]) {
                        sh "{ ${TFHOG_DIR}/bin/trufflehog --regex --json --max_depth 1 --rules ${TFHOG_DIR}/rules.json ${TARGET_REPO} > tfhog-report.json; } 2>/dev/null"
    }
                }
                    
            }
        }
        stage('SonarQube Analysis') {
            steps{
                script{
                    def scannerHome = tool 'SonarQube';
                    withSonarQubeEnv() {
                        sh "{ ${scannerHome}/bin/sonar-scanner; } 2>/dev/null"
                    }                    
                }
            }
        }
        stage('Gitleaks'){
            steps{
                echo '[*] Running Gitleaks ...'
                sh "{ ${GITLEAKS_DIR}/bin/gitleaks -p ${WORKSPACE} --config-path=${GITLEAKS_DIR}/gitleaks.toml --no-git -v -q > gitleaks-report.json; } 2>/dev/null"
            }
        }
        stage('Create Reporting'){
            steps{
                echo '[*] Create report ...'
                script {
                    def now = new Date()
                    env.REPORT_TIME = now.format("dd-MM-YYYY HH:mm:ss", TimeZone.getTimeZone('GMT+7'))

                    sh '{ python3 ${TFHOG_DIR}/convert.py --path ${WORKSPACE} --out ${REPORT_TIME} > ${WORKSPACE}/${REPORT_TIME}; } 2>/dev/null'
                    sh '{ cat ${REPORT_TIME}; } 2>/dev/null'
                    
                    ISSUE_COUNT = sh(
                        script: "{ grep -o 'Found IssuE' ${REPORT_TIME} | wc -l; } 2>/dev/null",
                        returnStdout: true
                    ).trim().toString()
                    echo "[*] Total Issue : ${ISSUE_COUNT}"
                }               
            }
        }
        stage('Upload Logs to GCS') {
            steps {
               step([$class: 'ClassicUploadStep', credentialsId: 'pharmalink-id', bucket: "gs://${env.GCS_BUCKET}", pattern: '${REPORT_TIME}.pdf'])
            }
        }
        stage('Compile') {
			steps {
				echo '> Building executable ...'
				sh 'make build'
			}
		}
		stage('Version') {
			steps {
				script {
					env.VERSION = sh(script: "jx-release-version", returnStdout: true).trim()
				}
				withCredentials([gitUsernamePassword(credentialsId: 'gitlab-pipeline-bot', gitToolName: 'git-tool')]) {
					sh "git config user.email '${env.PIPELINE_BOT_EMAIL}'"
					sh "git config user.name '${env.PIPELINE_BOT_NAME}'"
					sh "git tag -fa v${env.VERSION} -m '${env.VERSION}'"
					sh "git push origin v${env.VERSION}"
				}
			}
		}
		stage('Dockerize') {
			steps {
				script {
					echo '> Creating image ...'
					def dockerImage = docker.build("${PROJECT_ID}/${NAME}")
					echo '> Pushing image ...'
					docker.withRegistry("${DOCKER_REGISTRY_URL}", "gcr:pharmalink-id") {
						dockerImage.push("${env.VERSION}")
					}
				}
			}
		}
		stage('Helm Charts') {
			steps {
				echo '> Changing repository name value ...'
				sh "sed -i 's#repository: draft#repository: gcr.io/${ORG}-main/${NAME}#g' charts/values.yaml"
				echo '> Changing version value ...'
				sh "sed -i 's/tag: dev/tag: ${env.VERSION}/g' charts/values.yaml"
				echo '> Packing helm chart ...'
				sh "cd charts && helm package . --version=${env.VERSION}"
				echo '> Uploading chart ...'
				sh "cd charts && curl --data-binary '@${env.NAME}-${env.VERSION}.tgz' http://chartmuseum:8080/api/${env.ORG}/charts"
				echo '> Removing uploaded chart package ...'
				sh "rm charts/${env.NAME}-${env.VERSION}.tgz"
			}
		}        
    }
    post{
        success {
            build job: 'k8s-blue-sapphire-staging', parameters: [
				string(name: 'PROJECT_NAME', value: "${env.NAME}"),
				string(name: 'PROJECT_VERSION', value: "${env.VERSION}")
			], wait: false
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
