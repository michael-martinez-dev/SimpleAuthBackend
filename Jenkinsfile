pipeline {
  agent any
  stages {
    stage('Checkout Code') {
      parallel {
        stage('Checkout Code') {
          steps {
            git(url: 'https://github.com/MixedMachine/SimpleAuthBackend', branch: 'prod')
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
            sh 'make image'
          }
        }

        stage('Build resources') {
          steps {
            echo 'Building Databases & Storage resources...'
            sh 'make db'
          }
        }

        stage('Log into Docker') {
          steps {
            sh 'docker login -u $DOCKER_HUB_USER -p $DOCKER_HUB_PW'
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

    stage('Docker Hub push') {
      steps {
        echo 'Pushing to Dockerhub...'
        sh 'make image-push'
      }
    }

    stage('Prod env set-up') {
      steps {
        echo 'Setting up production environment...'
      }
    }

  }
  post {
      always {
          sh 'make clean'
      }
      success {
          echo 'The Pipeline was successful! 🎉'
      }
      failure {
          echo'The Pipeline failed 😔'
      }
  }
}