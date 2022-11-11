# repo-analytics
Analytics about the users who care about your repo

# Usage
To avoid throttling, you need a Github token, set as environment variable
```
export GITHUB_TOKEN="...
```

## Get emails of all stargazers
For a given repository, say [Jeadie/repo-analytics](https://github.com/Jeadie/repo-analytics)
```bash
 ./stars.sh Jeadie repo-analytics | xargs -I% -P0 ./user_email.sh %
```