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

        stage('Env vars set-up') {
          steps {
            sh 'echo "LOG_LEVEL=$SAB_LOG_LEVEL" >> .env'
            sh 'echo "API_PORT=$SAB_API_PORT" >> .env'
            sh 'echo "DATABASE_USER=$SAB_DATABASE_USER" >> .env'
            sh 'echo "DATABASE_PASS=$SAB_DATABASE_PASS" >> .env'
            sh 'echo "DATABASE_HOST=$SAB_DATABASE_HOST" >> .env'
            sh 'echo "DATABASE_PORT=$SAB_DATABASE_PORT" >> .env'
            sh 'echo "DATABASE_NAME=$SAB_DATABASE_NAME" >> .env'
            sh 'echo "DATABASE_COLLECTION=$SAB_DATABASE_COLLECTION" >> .env'
            sh 'echo "REDIS_HOST=$SAB_REDIS_HOST" >> .env'
            sh 'echo "REDIS_PORT=$SAB_REDIS_PORT" >> .env'
            sh 'echo "REDIS_USER=$SAB_REDIS_USER" >> .env'
            sh 'echo "REDIS_PASS=$SAB_REDIS_PASS" >> .env'
            sh 'echo "REDIS_DB=$SAB_REDIS_DB" >> .env'
            sh 'echo "JWT_SECRET_KEY=$SAB_JWT_SECRET_KEY" >> .env'
            sh 'echo "" >> .env'
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