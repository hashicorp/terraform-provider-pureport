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
        TF_LOG=INFO
        PUREPORT_ENDPOINT="https://dev1-api.pureportdev.com"
        PUREPORT_API_KEY="mKBkM3l1ScUHW"
        PUREPORT_API_SECRET="JMzOfGAbLRcrNziGO"
    }
    stages {
        stage('Build') {
            steps {
                sh "make"
            }
        }
        stage('Run Terraform Tests') {
            when {

                // This can take a long time so we may only want to do this on develop
                branch 'develop'
            }
            steps {
                sh "make testacc"
            }
            post {
                cleanup {
                    // Make sure we cleanup an Networks or Connections that may be still hanging
                    // around here. The Terraform Test framework should always call delete on the
                    // plugin, but you never know.
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
