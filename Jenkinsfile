pipeline {
  agent any
  stages {
    stage('Checkout Code') {
      parallel {
        stage('Checkout Code') {
          steps {
            git(url: 'https://github.com/MixedMachine/simple-signin-backend', branch: 'prod')
            sh 'go mod tidy'
          }
        }

        stage('Log') {
          steps {
            sh 'ls -la'
            sh 'go version'
            sh 'docker version'
          }
        }

      }
    }

    stage('Unit tests') {
      steps {
        echo 'Running Unit tests...'
        sh 'go test ./tests/unit/...'
      }
    }

    stage('Build images') {
      parallel {
        stage('Build images') {
          steps {
            echo 'Building docker images & pushing them to repo...'
          }
        }

        stage('Build resources') {
          steps {
            echo 'Building Databases & Storage resources...'
          }
        }

      }
    }

    stage('Run service') {
      steps {
        echo 'Running service with docker to run functional testing...'
      }
    }

    stage('Functional tests') {
      steps {
        echo 'Running functional tests with postman...'
      }
    }

    stage('Prod env set-up') {
      steps {
        echo 'Setting up production environment...'
      }
    }

  }
}