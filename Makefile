include .env

# found 'watcher' at https://github.com/canthefason/go-watcher
# that wasn't working as expected, so found and switched 'fswatch' at: https://github.com/emcrisostomo/fswatch

default: build \
	local-api

build:
	sam build
	@cp ./config/defaults.yml ./.aws-sam/build/PDFGeneratorFunction/

deploy:	build \
	moveDefaults \
	dev-cloud

moveDefaults: 
	aws s3 cp ./config/defaults.yml s3://$(S3_STORAGE_BUCKET)/public/ && \
	aws s3api put-object-tagging \
  --bucket $(S3_STORAGE_BUCKET) \
  --key public/defaults.yml \
  --tagging '{"TagSet": [{"Key": "public", "Value": "true"}]}' && \
	echo "default.yml uploaded and tagged"

# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-start-api.html
local-api:
	sam local start-api --env-vars env.json --profile $(PROFILE)

local-invoke:
	sam local invoke --env-vars env.json --profile $(PROFILE)

dev-cloud:
	sam sync --stack-name $(STACK_NAME) --profile $(PROFILE) \
	--s3-prefix $(AWS_DEPLOYMENT_PREFIX) \
	--parameter-overrides \
		ParamCertificateArn=$(CERTIFICATE_ARN) \
		ParamCustomDomainName=$(CUSTOM_DOMAIN_NAME) \
		ParamHostedZoneId=$(HOSTED_ZONE_ID) \
		ParamKMSKeyID=$(KMS_KEY_ID) \
		ParamSSMPath=$(SSM_PARAM_PATH) \
		ParamStorageBucket=${S3_STORAGE_BUCKET}

dev-cloud-watch:
	sam sync --stack-name $(STACK_NAME) --watch --profile $(PROFILE) \
	--s3-prefix $(AWS_DEPLOYMENT_PREFIX) \
	--parameter-overrides \
		ParamCertificateArn=$(CERTIFICATE_ARN) \
		ParamCustomDomainName=$(CUSTOM_DOMAIN_NAME) \
		ParamHostedZoneId=$(HOSTED_ZONE_ID) \
		ParamKMSKeyID=$(KMS_KEY_ID) \
		ParamSSMPath=$(SSM_PARAM_PATH) \
		ParamStorageBucket=${S3_STORAGE_BUCKET}

tail-logs:
	sam logs -n PDFGeneratorFunction --profile $(PROFILE) \
	--stack-name $(STACK_NAME) --tail

tail-logs-trace:
	sam logs -n PDFGeneratorFunction --profile $(PROFILE) \
	--stack-name $(STACK_NAME) --tail --include-traces

validate:
	sam validate
	
watch:
	fswatch -o ./ | xargs -n1 -I{} sam build

test:
	go test -v