go build -o build/wordDeJourMac
GOOS=windows GOARCH=386; go build -o build/wordDeJour.exe
GOOS=linux  GOARCH=386; go build -o build/wordDeJourLinux

 ./build/wordDeJourMac


