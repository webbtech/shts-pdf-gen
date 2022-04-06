include .env

# found 'watcher' at https://github.com/canthefason/go-watcher
# that wasn't working as expected, so found and switched 'fswatch' at: https://github.com/emcrisostomo/fswatch

.PHONY: build

build:
	sam build
	@cp ./config/defaults.yml ./.aws-sam/build/PDFGeneratorFunction/

# https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-start-api.html
local-api: build
	sam local start-api --env-vars env.json --profile $(PROFILE)
	#  -p 3100

local-invoke:
	sam local invoke --env-vars env.json --profile $(PROFILE)

dev-cloud:
	sam  sync --stack-name $(STACK_NAME) --profile $(PROFILE)
	# --parameter-overrides \
		# ParamMailRecipient=$(MAIL_RECIPIENT) \
		# ParamMailSender=$(MAIL_SENDER)

dev-cloud-watch:
	sam  sync --stack-name $(STACK_NAME) --watch --profile $(PROFILE)
	# --parameter-overrides \
		# ParamMailRecipient=$(MAIL_RECIPIENT) \
		# ParamMailSender=$(MAIL_SENDER)

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