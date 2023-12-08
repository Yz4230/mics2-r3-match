build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build .

deploy:
	gcloud compute scp ./shogi-match instance-1:~

deploy-ced:
	scp ./shogi-match ced:~/mics2-r3/source/shogi-match
