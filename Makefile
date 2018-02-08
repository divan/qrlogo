test:
	go run cmd/qrlogo/main.go -i=example/color.png -o=example/qr_logo.png -k=true -size=512 http://www.163.com
	go run cmd/qrlogo/main.go -b=example/bg.png   -o=example/qr_bg.png -ox=50 -oy=100 http://www.163.com