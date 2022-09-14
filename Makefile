.PHONY: update

update:
	@cd cmd; go run ./main.go $(MY_SPREADSHEET_IDS)
	@git add ./cmd/dist/.
	@git commit -m "updated at `date "+%Y/%m/%d %H:%M:%S"` by command"
