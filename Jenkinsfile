pipeline {
  agent {
    kubernetes {
      yaml """
            apiVersion: v1
            kind: Pod
            spec:
              nodeSelector:
                cloud.google.com/gke-spot: "true"
              tolerations:
              - key: "purpose"
                operator: "Equal"
                value: "gitlab-runner"
                effect: "NoSchedule"
              containers:
              - name: golang
                image: golang:1.22
                command: ['sleep']
                args: ['infinity']
                resources:
                  requests:
                    cpu: "1"
                    memory: "1Gi"
        """
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



