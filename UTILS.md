## Several utility scripts

#### Delete all tags in GitHub repository
```bash
#Delete local tags.
git tag -l | xargs git tag -d
#Fetch remote tags.
git fetch
#Delete remote tags.
git tag -l | xargs -n 1 git push --delete origin
#Delete local tasg.
git tag -l | xargs git tag -d
```

#### Clean Go modules cache
```bash
go clean -modcache
```

#### Run all tests in the project
```bash
go test ./...
```

