# Binary name
BINARY= Zombie
VERSION = 0.9.9beta
# Builds the project
build:
		go build -ldflags "-s -w" -o ${BINARY} ./src/main.go

build-windows-64:
		CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -a -ldflags '-linkmode "external" -extldflags "-static"' -o ${BINARY}-windowsx64.exe src/main.go

build-linux-64:
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -a -ldflags '-linkmode external -extldflags -static' -o ${BINARY}-linux64 src/main.go

# Installs our project: copies binaries
install:
		go install
release-upx:
		# Clean
		#go clean
		rm -rf *.gz
		# Build for mac
		go build -ldflags "-s -w" -o ./bin/Zombie-mac64-${VERSION} ./src/main.go
		upx -2 ./bin/Zombie-mac64-${VERSION}
		# Build for linux
		#go clean

		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/Zombie-linux64-${VERSION} ./src/main.go
		upx -2 ./bin/Zombie-linux64-${VERSION}

		#go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/Zombie-linux32-${VERSION} ./src/main.go
		upx -2 ./bin/Zombie-linux32-${VERSION}
		# Build for win
		#go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/Zombie-win64-${VERSION}.exe  ./src/main.go
		upx -2 ./bin/Zombie-win64-${VERSION}.exe
		#go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/Zombie-win32-${VERSION}.exe ./src/main.go
		upx -2 ./bin/Zombie-win32-${VERSION}.exe
		#compress
		cp ./ReadMe.md ./bin/
		tar cvf Zombie-${VERSION}.tar.gz bin/*
# Cleans our projects: deletes binaries

release:
# Clean
		#go clean
		rm -rf *.gz
		# Build for mac
		go build -ldflags "-s -w" -o ./bin/Zombie-mac64-${VERSION} ./src/main.go
		# Build for linux
		#go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/Zombie-linux64-${VERSION} ./src/main.go
		#go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/Zombie-linux32-${VERSION} ./src/main.go
		# Build for win
		#go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/Zombie-win64-${VERSION}.exe  ./src/main.go
		#go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/Zombie-win32-${VERSION}.exe ./src/main.go
		#compress
		tar cvf Zombie.tar.gz bin/*

clean:
		go clean

.PHONY:  clean build