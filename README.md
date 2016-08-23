# Cloud Foundry App Instance Usage CLI Plugin

Cloud Foundry plugin extension to view the instance count for an application that is running in a Cloud Foundry deployment.

## Install

```
$ go get github.com/cglean/app-instances
$ cf install-plugin $GOPATH/bin/app-instances
```

## Usage

**SAMPLE OUTPUT**

```
$ cf app-instances 'application-name'

Following is the output 
1
```

## Uninstall

```
$ cf uninstall-plugin app-instances
```

## Jenkins Usage
```
export APP_NAME=the-app-name
INSTANCES=$(cf app-instances $APP_NAME)
cf scale $APP_NAME -i $INSTANCES
```

