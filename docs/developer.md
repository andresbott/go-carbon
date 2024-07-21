# Running locally

### The Go service

The Go service can be started locally with:
```
make run
```
Access the backend on  http://localhost:8085

Optionally you can build the SPA and package into the go app with
```
make package-ui
```

### Node dev server

The nodejs dev server can be started locally and proxy request to the local go service.

1. first stat the GO backend service
2. start the frontend 
```
cd webui
make run
```
Access the SPA on  http://localhost:5173/spa

# Validations

### Code linting

Run with 
```
make lint
```
### License check

Depends on https://github.com/elastic/go-licence-detector and is used to ensure
that only valid licenses are part of the project, including dependencies.

```
make license-check 
```

# Release

The process to release a new version is as follows:

1. make sure you have checked out the main branch and git is a clean state
2. get a GitHub Token from: https://github.com/settings/tokens/new?scopes=repo,write:packages
3. run `make release version="vx.y.z"`

The make target will test, build, tag and release to GitHub all binaries related to the new release.

# CI/CD
There are GitHub actions to validate all commits pushed, check [.github](.github)