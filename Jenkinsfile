#!/usr/bin/env groovy

@Library('jenkins-build-utils')_

def utils = new com.pureport.Utils()

def version = "1.0.0"
def plugin_name = "terraform-provider-pureport"

pipeline {
    agent {
      docker {
        image 'golang:1.12'
      }
    }
    options {
        disableConcurrentBuilds()
    }
    parameters {
      booleanParam(
          name: 'ACCEPTANCE_TESTS_RUN',
          defaultValue: false,
          description: 'Should we run the acceptance tests as part of run?'
          )
      booleanParam(
          name: 'ACCEPTANCE_TESTS_LOG_TO_FILE',
          defaultValue: true,
          description: 'Should debug logs be written to a separate file?'
          )
      choice(
          name: 'ACCEPTANCE_TESTS_LOG_LEVEL',
          choices: ['WARN', 'ERROR', 'DEBUG', 'INFO', 'TRACE'],
          description: 'The Terraform Debug Level'
          )
      choice(
          name: 'ACC_TEST_ENVIRONMENT',
          choices: ['default', 'Production', 'Dev1'],
          description: 'The environment to deploy Terraform Acceptance Tests'
          )
    }
    environment {
        GOPATH                = "/go"
        GOCACHE               = "/tmp/go/.cache"
    }
    stages {
        stage('Configure') {
            steps {
                script {

                    // Setup the test environment
                    def environment = params.ACC_TEST_ENVIRONMENT
                    def provider_version = ""

                    provider_version += "v${version}"

                    // Only add the build version for the develop branch
                    if (env.BRANCH_NAME == "develop") {
                      provider_version += "-b${env.BUILD_NUMBER}"
                    }


                    // If the environment is specified to be the default,
                    // use the branch name to determine the environment
                    if (params.ACC_TEST_ENVIRONMENT == "default") {

                      switch (env.BRANCH_NAME) {

                      case ~/release\/.*/:
                        environment = "Production"

                      default:
                        environment = "Dev1"
                      }
                    }

                    plugin_name += "_${provider_version}"

                    env.PUREPORT_ACC_TEST_ENVIRONMENT = environment
                    env.PROVIDER_VERSION = provider_version
                }
            }
        }
        stage('Build') {
            steps {

                retry(3) {
                  sh "make"
                  sh "PROVIDER_VERSION=${env.PROVIDER_VERSION} make plugin"
                  sh "chmod +x ${plugin_name}"

                  archiveArtifacts(
                      artifacts: "${plugin_name}"
                      )
                }
            }
        }
        stage('Run Terraform Tests') {
            when {

                // This can take a long time so we may only want to do this on develop
                anyOf {
                  branch 'develop'
                  branch 'release/*'
                  expression { return params.ACCEPTANCE_TESTS_RUN }
                }
            }
            environment {
                TF_LOG                = "${params.ACCEPTANCE_TESTS_LOG_LEVEL}"
                TF_LOG_PATH           = "${params.ACCEPTANCE_TESTS_LOG_TO_FILE ? 'tf_log.log' : '' }"
                TF_IN_AUTOMATION      = "true"
                TF_INPUT              = "false"

                GOOGLE_CREDENTIALS    = credentials('terraform-google-credentials-id')
                GOOGLE_PROJECT        = "pureport-customer1"
                GOOGLE_REGION         = "us-west2"

                AWS_DEFAULT_REGION    = "us-east-1"
                AWS_ACCESS_KEY_ID     = credentials('terraform-acc-test-aws-key-id')
                AWS_SECRET_ACCESS_KEY = credentials('terraform-acc-test-aws-secret')

                ARM_CLIENT_ID         = credentials('terraform-acc-test-azure-client-id')
                ARM_CLIENT_SECRET     = credentials('terraform-acc-test-azure-client-secret')
                ARM_SUBSCRIPTION_ID   = credentials('terraform-acc-test-azure-subscription-id')
                ARM_TENANT_ID         = credentials('terraform-acc-test-azure-tenant-id')
                ARG_USE_MSI           = true
            }
            stages {

                stage('in Dev1') {
                    when {
                      expression { return env.PUREPORT_ACC_TEST_ENVIRONMENT == "Dev1" }
                    }
                    environment {
                      PUREPORT_ENDPOINT     = "https://dev1-api.pureportdev.com"
                      PUREPORT_API_KEY      = credentials('terraform-pureport-dev1-api-key')
                      PUREPORT_API_SECRET   = credentials('terraform-pureport-dev1-api-secret')
                    }
                    steps {
                        script {
                            sh "make testacc"
                        }
                    }
                }

                stage('in Production') {
                    when {
                      expression { return env.PUREPORT_ACC_TEST_ENVIRONMENT == "Production" }
                    }
                    environment {
                      PUREPORT_ENDPOINT     = "https://api.pureport.com"
                      PUREPORT_API_KEY      = credentials('terraform-testacc-prod-key-id')
                      PUREPORT_API_SECRET   = credentials('terraform-testacc-prod-secret')
                    }
                    steps {
                        script {
                            sh "make testacc"
                        }
                    }
                }
            }
            post {
                always {

                    archiveArtifacts(
                        allowEmptyArchive: true,
                        artifacts: 'pureport/tf_log.log'
                    )
                }
            }
        }
        stage('Copy plugin to Nexus') {
            when {

                // This can take a long time so we may only want to do this on develop
                anyOf {
                  branch 'develop'
                  branch 'release/*'
                }
            }
            steps {
                script {
                    withCredentials([
                        usernamePassword(
                          credentialsId: 'nexus_credentials',
                          usernameVariable: 'nexusUsername',
                          passwordVariable: 'nexusPassword'
                          )
                    ]) {

                      def nexus_url = "https://nexus.dev.pureport.com/repository/terraform-provider-pureport/${env.BRANCH_NAME}/"

                      sh "curl -v -u ${nexusUsername}:${nexusPassword} --upload-file ${plugin_name} ${nexus_url}"

                      // Set the description text for the job
                      currentBuild.description = "Version: ${plugin_name}"

                    }
                }
            }
        }
    }
    post {
        success {
            slackSend(color: '#30A452', message: "SUCCESS: <${env.BUILD_URL}|${env.JOB_NAME}#${env.BUILD_NUMBER}>")
        }
        unstable {
            slackSend(color: '#DD9F3D', message: "UNSTABLE: <${env.BUILD_URL}|${env.JOB_NAME}#${env.BUILD_NUMBER}>")

            script {
                utils.sendUnstableEmail()
            }
        }
        failure {
            slackSend(color: '#D41519', message: "FAILED: <${env.BUILD_URL}|${env.JOB_NAME}#${env.BUILD_NUMBER}>")
            script {
                utils.sendFailureEmail()
            }
        }
    }
}
