GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
BINARY_NAME = cloud_function

# GCP cloud function deploy
CONNECTOR = projects/serverless-test-XXXXX/locations/us-central1/connectors/serverless-connector
REGION = us-central1
REDISHOST = 10.0.0.1
REDISPORT = 6379
SHORTURLBASE = https://us-central1-serverless-test-XXXXX.cloudfunctions.net/Redirect/
PROJECTID = serverless-test-XXXXX
ISONGCP = True
TOPICID = short-url-source-topic

all: test deploy_register deploy_redirect

test:
	$(GOTEST) -v ./...

deploy_register:
	gcloud functions deploy Register \
		--runtime go113 \
		--trigger-http \
		--region $(REGION) \
		--vpc-connector $(CONNECTOR) \
		--set-env-vars REDISHOST=$(REDISHOST),REDISPORT=$(REDISPORT),SHORTURLBASE=$(SHORTURLBASE),ISONGCP=$(ISONGCP),PROJECTID=$(PROJECTID)

deploy_redirect:
	gcloud functions deploy Redirect \
		--runtime go113 \
		--trigger-http \
		--region $(REGION) \
		--vpc-connector $(CONNECTOR) \
		--set-env-vars REDISHOST=$(REDISHOST),REDISPORT=$(REDISPORT),ISONGCP=$(ISONGCP),PROJECTID=$(PROJECTID),TOPICID=$(TOPICID)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm $(BINARY_NAME)