# Design Doc Self Generator

Self-generates design docs.

## Running
1. Copy `.env.example` to `.env`

2. Add env. variables to: `.env`.

3. Run the app:
```
docker compose up
```

## GraphQL Playground

[http://localhost:8080/graphql](http://localhost:8080/graphql)

### Examples

Obtaining pull request data from GitHub by URLs:

```graphql
query {
  solutions {
    analysis {
      information {
        github {
          metrics {
            pullRequests(urls: ["https://github.com/graphql-go/graphql/pull/117"]) {
              url
              duration {
                inDays
                formattedIntervalDates
              }
              formattedContributors
            }
          }
        }
      }
    }
  }
}
```

Test output:
```json
{
  "data": {
    "solutions": [
      {
        "analysis": [
          {
            "information": [
              {
                "github": {
                  "metrics": {
                    "pullRequests": [
                      {
                        "duration": {
                          "formattedIntervalDates": "2016-03-07 09:07:29 +0000 UTC - 2016-05-30 01:52:47 +0000 UTC",
                          "inDays": 83
                        },
                        "formattedContributors": "- https://github.com/sogko</br>- https://github.com/coveralls</br>- https://github.com/pspeter3</br>- https://github.com/chris-ramon</br>- https://github.com/jvatic",
                        "url": "https://github.com/graphql-go/graphql/pull/117"
                      }
                    ]
                  }
                }
              }
            ]
          }
        ]
      }
    ]
  }
}
```

## Features

Contains the following features:
- [x] Data from GitHub.
- [x] Env Variables.
- [x] Config.
- [x] Auth.
- [x] JWT.
- [x] GraphQL.
- [x] PostgreSQL.
- [x] Type Safe SQL.
- [x] Docker Compose.
- [x] Live reload.
- [x] Admin.
- [x] Unit tests.
