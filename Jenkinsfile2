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
	post {
		success {
			build job: 'k8s-blue-sapphire-staging', parameters: [
				string(name: 'PROJECT_NAME', value: "${env.NAME}"),
				string(name: 'PROJECT_VERSION', value: "${env.VERSION}")
			], wait: false

			discordSend link: env.BUILD_URL, result: currentBuild.currentResult, title: "${env.JOB_NAME} #${env.BUILD_NUMBER}", webhookURL: "${env.DISCORD_WEBHOOK_URL}"
			sh "exit 0"
		}

		regression {
			discordSend link: env.BUILD_URL, result: currentBuild.currentResult, title: "${env.JOB_NAME} #${env.BUILD_NUMBER}", webhookURL: "${env.DISCORD_WEBHOOK_URL}"
			sh "exit 1"
		}
	}
}
