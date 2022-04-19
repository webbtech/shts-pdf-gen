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
	aws s3 cp ./config/defaults.yml s3://$(S3_STORAGE_BUCKET)/public/

# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-start-api.html
local-api:
	sam local start-api --env-vars env.json --profile $(PROFILE)

local-invoke:
	sam local invoke --env-vars env.json --profile $(PROFILE)

dev-cloud:
	sam sync --stack-name $(STACK_NAME) --profile $(PROFILE) \
	--s3-prefix $(AWS_DEPLOYMENT_PREFIX) \
	--parameter-overrides \
		ParamKMSKeyID=$(KMS_KEY_ID) \
		ParamSSMPath=$(SSM_PARAM_PATH) \
		ParamStorageBucket=${S3_STORAGE_BUCKET}

dev-cloud-watch:
	sam sync --stack-name $(STACK_NAME) --watch --profile $(PROFILE) \
	--s3-prefix $(AWS_DEPLOYMENT_PREFIX) \
	--parameter-overrides \
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

# TODO: add command for running go tests