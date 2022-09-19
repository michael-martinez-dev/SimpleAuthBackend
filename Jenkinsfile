pipeline {
  agent any
  stages {
    stage('Checkout Code') {
      steps {
        git(url: 'https://github.com/MixedMachine/simple-signin-backend', branch: 'prod')
      }
    }

    stage('Log') {
      steps {
        sh 'ls -la'
      }
    }

  }
}