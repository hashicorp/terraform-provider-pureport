#!/usr/bin/env groovy

@Library('jenkins-build-utils')_

def utils = new com.pureport.Utils()

pipeline {
    agent {
      docker {
        image 'golang:1.12'
      }
    }
    options {
        disableConcurrentBuilds()
    }
    environment {
        TF_LOG              = "INFO"
        GOPATH              = "/go"
        GOCACHE             = "/tmp/go/.cache"
        PUREPORT_ENDPOINT   = "https://dev1-api.pureportdev.com"
        PUREPORT_API_KEY    = "mKBkM3l1ScUHW"
        PUREPORT_API_SECRET = "JMzOfGAbLRcrNziGO"
        GOOGLE_CREDENTIALS  = credentials('terraform-google-credentials-id')
        GOOGLE_PROJECT      = "pureport-customer1"
        GOOGLE_REGION       = "us-west2"
    }
    parameters {
      booleanParam(
          name: 'RUN_ACCEPTANCE_TESTS',
          defaultValue: false,
          description: 'Should we run the acceptance tests as part of run?'
          )
    }
    stages {
        stage('Build') {
            steps {

                retry(3) {
                  sh "echo $GOOGLE_CREDENTIALS"
                  sh "make"
                }
            }
        }
        stage('Run Terraform Tests') {
            when {

                // This can take a long time so we may only want to do this on develop
                anyOf {
                  branch 'develop'
                  expression { return params.RUN_ACCEPTANCE_TESTS } 
                }
            }
            steps {
                sh "make testacc"
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
