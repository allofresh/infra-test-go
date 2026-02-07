pipeline {
    agent {
        kubernetes {
            yaml '''
            spec:
              containers:
              - name: golang
                image: golang:1.22
                command: ['sleep']
                args: ['infinity']
            '''
            defaultContainer 'golang'
        }
    }

    stages {
        stage('Compile') {
            steps {
                sh 'make compile'
            }
        }

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
