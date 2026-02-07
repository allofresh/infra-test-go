pipeline {
    agent any

    stages {
        stage('Vet') {
            steps {
                sh 'make vet'
            }
        }

        stage('Unit Test') {
            steps {
                sh 'make unit-test'
            }
        }

        stage('Race Test') {
            steps {
                sh 'make race-test'
            }
        }

        stage('Coverage') {
            steps {
                sh 'make coverage'
            }
        }
    }
}



