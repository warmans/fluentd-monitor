set -e

#help text
if [ "$1" == "help" ]; then
    echo "Usage ./package.sh [help|rpm|deb|solaris|puppet]"
    exit
fi

#validate fpm is installed
if [ -z "$(which fpm)" ]; then
    printf "error:\nPackaging requires effing package manager (fpm) to run.\nsee https://github.com/jordansissel/fpm\n"
    exit 1
fi

#set target package format
TARGET=$1
if [ -z "$TARGET" ]; then
    TARGET="rpm"
fi

echo "Building $TARGET package..."

#build in pkg directory
export DESTDIR=pkg

#run install using dest DESTDIR prefix
make install PREFIX=/usr

#clean
if [ -d "dist" ]; then
    rm -f dist/*.$TARGET
else
    mkdir dist
fi

#build RPM
fpm --rpm-os el6 -s dir -p dist -t $TARGET -n fluentd-monitor -v $(cat version) -C $DESTDIR .
