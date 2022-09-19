pipeline {
  agent any
  stages {
    stage('Checkout Code') {
      steps {
        git(url: 'https://github.com/MixedMachine/simple-signin-backend', branch: 'prod')
        sh 'export PATH="/go/bin:/usr/local/go/bin:$PATH"'
      }
    }

    stage('Log') {
      parallel {
        stage('Log') {
          steps {
            sh 'ls -la'
            sh 'ls -la /'
            sh 'ls -la /usr/local/'
          }
        }

        stage('Unit Tests') {
          steps {
            sh '''go version

&& go test ./...'''
          }
        }

      }
    }

  }
}