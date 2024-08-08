pipeline {
  agent any
  stages {
    stage('Build') {
      parallel {
        stage('Build') {
          steps {
            echo 'Building the app'
          }
        }

        stage('Test') {
          steps {
            echo 'Testing'
          }
        }

        stage('TestLog') {
          steps {
            writeFile(file: 'Testlog.txt', text: 'A logger for test ')
          }
        }

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