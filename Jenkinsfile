pipeline {
    agent any

    parameters {
        string(name: 'BRANCH', defaultValue: 'development', description: 'Git branch to test')
        choice(name: 'TEST_TAG', choices: ['smoke', 'ui', 'regression'], description: 'Test tag to run')
        string(name: 'API_BASE_URL', defaultValue: 'https://demoqa.com', description: 'API Base URL')
        string(name: 'UI_BASE_URL', defaultValue: 'https://demoqa.com', description: 'UI Base URL')
        string(name: 'USERNAME', defaultValue: 'useruser', description: 'Username for login')
        string(name: 'PASSWORD', defaultValue: 'P@ssw0rd', description: 'Password for login')
        string(name: 'HEADLESS', defaultValue: 'true', description: 'Run browser in headless mode (true/false)')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: "*/${params.BRANCH}"]],
                    extensions: [],
                    userRemoteConfigs: [[
                        url: '<your-rep-ip>',
                        credentialsId: 'rep-git-creds'
                    ]]
                ])
            }
        }

        stage('Debug Docker') {
            steps {
                sh 'which docker || echo "❌ docker CLI not found"'
                sh 'docker --version || echo "❌ docker not executable"'
                sh 'docker info | head -n 3 || echo "❌ docker daemon not reachable"'
            }
        }

        stage('Run Go Tests in Docker') {
            steps {
                sh '''
                    echo "🔧 Running Go UI tests in Docker container..."

                    docker run --rm \
                      -v $(pwd):/app \
                      -v $(pwd)/allure-results:/app/allure-results \
                      -w /app \
                      -e USERNAME="${USERNAME}" \
                      -e PASSWORD="${PASSWORD}" \
                      -e API_BASE_URL="${API_BASE_URL}" \
                      -e UI_BASE_URL="${UI_BASE_URL}" \
                      -e HEADLESS="${HEADLESS}" \
                      golang:1.23 \
                      bash -c "
                        apt-get update && apt-get install -y chromium chromium-driver \
                        && go mod download \
                        && go test -v -timeout=15m ./tests/... \
                             -args \
                               -username=\$USERNAME \
                               -password=\$PASSWORD \
                               -api-base-url=\$API_BASE_URL \
                               -ui-base-url=\$UI_BASE_URL \
                               -headless=\$HEADLESS \
                               -allure-results-dir=/app/allure-results
                      "
                '''
            }
        }

        stage('Allure Report Ready') {
            steps {
                echo "✅ Allure results saved to ./allure-results"
                echo "📊 Live report: http://<your-server-ip>:5050"
            }
        }
    }

    post {
        always {
            echo "✅ Pipeline finished for branch: ${params.BRANCE}"
        }
        failure {
            echo "❌ Pipeline failed! Check Docker access and test logs."
        }
    }
}