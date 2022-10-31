swag:
	swag fmt
	swag init --parseDependency --parseInternal

gosdk:
	make swag
	rm -rf $(ls -A ${CNSBACKEND_SDK_DIR_GO}| grep -v .git)
	openapi-generator generate -c ./code_gen/go_config.yaml -o ${CNSBACKEND_SDK_DIR_GO}
	cd ${CNSBACKEND_SDK_DIR_GO} && go mod tidy && go mod edit -module github.com/wangdayong228/cns-backend-sdk-go
