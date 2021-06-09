GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
BINARY_NAME = cloud_function

# GCP cloud function deploy
PROJECTID = serverless-test-XXXXX
CONNECTOR = projects/$(PROJECTID)/locations/us-central1/connectors/serverless-connector
REGION = us-central1
REDISHOST = 10.35.95.115
REDISPORT = 6379
SHORTURLBASE = https://us-central1-$(PROJECTID).cloudfunctions.net/Redirect/
ISONGCP = True
TOPICID = short-url-source-topic

all: test deploy_register deploy_redirect deploy_register_with_auth

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

deploy_register_with_auth:
	gcloud functions deploy RegisterWithAuth \
		--runtime go113 \
		--trigger-http \
		--region $(REGION) \
		--vpc-connector $(CONNECTOR) \
		--set-env-vars REDISHOST=$(REDISHOST),REDISPORT=$(REDISPORT),SHORTURLBASE=$(SHORTURLBASE),ISONGCP=$(ISONGCP),PROJECTID=$(PROJECTID)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm $(BINARY_NAME)