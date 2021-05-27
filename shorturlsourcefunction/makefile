GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
BINARY_NAME = cloud_function

# GCP cloud function deploy
PROJECTID = serverless-test-XXXXX
CONNECTOR = projects/$(PROJECTID)/locations/us-central1/connectors/serverless-connector
REGION = us-central1
TOPICID = short-url-source-topic

all: deploy

deploy:
	gcloud functions deploy ShortUrlSource \
		--runtime go113 \
		--trigger-topic $(TOPICID) \
		--region $(REGION) \
		--vpc-connector $(CONNECTOR)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm $(BINARY_NAME)