## Backend markdown notes app

This is a personal project I built to experiment with [Next.js](https://nextjs.org/) and [AWS](https://aws.amazon.com/) services.

This application sets up a REST API for fetching and updating notes for the application.

The backend is built with go as this is my current most used backend langauge.


## Frontend repo

The frontend of the application can be found [here](https://github.com/KyleJonesNV/frontend-notes)


## Access the API

The API is currently exposing the following endpoints
<ol>
  <li>/getAllForUser</li>
  <li>/insertTopic</li>
  <li>/deleteTopic</li>
  <li>/insertNote</li>
  <li>/getAllNotes</li>
  <li>/deleteNote</li>
</ol> 


The API is hosted in AWS here:

`https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging`

I have also setup a test user with id: 

`1d7ee7f0-36f5-4e33-a766-26981e62d9cf`

If you have Curl installed you can test getting all topics for the test user using:

```
curl -sX POST https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/getAllForUser -d '{"id": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf"}'
```

Insert a new topic with:

```
curl -sX POST https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/insertTopic -d '{"userID": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "something interesting"}'
```

Delete the new topic with:

```
curl -sX DELETE https://ifhrxwl601.execute-api.eu-west-1.amazonaws.com/staging/deleteTopic -d '{"userID": "1d7ee7f0-36f5-4e33-a766-26981e62d9cf", "title": "something interesting"}'
```


## Improvements / things I would like to do next

<ol>
  <li>Export API to swagger file for easier communication</li>
  <li>Better error handling for incorrect inputs</li>
  <li>More tests for each endpoint</li>
  <li>Support for running locally with docker</li>
</ol> 
