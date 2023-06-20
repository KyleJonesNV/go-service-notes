run build script

`sh scripts/build.sh`

Update lambda

`aws lambda update-function-code --function-name go-service-notes \
    --zip-file fileb://./deployment.zip \
    --region eu-west-1`

Update lambda using S3 bucket

`aws lambda update-function-code --function-name go-service-notes \
    --s3-bucket notes-service-packt \
    --s3-key deployment.zip \
    --region eu-west-1`


https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/insertTopic

curl -sX POST https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/insertTopic -d '{"userID": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "something interesting"}'

curl -sX POST https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/getAllForUser -d '{"userID": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf"}'



curl -sX DELETE https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/deleteNote -d '{"userID": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "Fun stuff", "note": "Building fun stuff"}'