pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        echo 'Building the app'
      }
    }

    stage('Test') {
      steps {
        echo 'Testing the app'
      }
    }

    stage('Deploy ') {
      steps {
        input(message: 'Do you want to deploy', id: 'Ok')
        echo 'Deploying the app'
      }
    }

  }
}