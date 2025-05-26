pipeline {
    agent any

    stages {
        stage('Lint') {
            steps {
                echo 'Running linter...'
                // Simulate success
                echo '✅ Lint passed!'
            }
        }

        stage('Tests') {
            steps {
                echo 'Running tests...'
                // Simulate success
                echo '✅ All tests passed!'
            }
        }

        stage('Security') {
            steps {
                echo 'Running security checks...'
                // Simulate success
                echo '✅ No security issues found!'
            }
        }
    }

    post {
        always {
            echo 'Build completed.'
        }
    }
}
