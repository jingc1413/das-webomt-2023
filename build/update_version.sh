set -e

if [ $1 = "clean" ]; then
sed -i "s/version:\(\s*\)\".*\"\(\s*\),/version:\1\"\"\2,/g" app-v3/src/settings.js
sed -i "s/build:\(\s*\)\".*\"\(\s*\),/build:\1\"\"\2,/g" app-v3/src/settings.js
sed -i "s/VERSION\(\s*\)=\(\s*\)\".*\"/VERSION\1=\2\"\"/g" gomt/core/info.go
sed -i "s/BUILD\(\s*\)=\(\s*\)\".*\"/BUILD\1=\2\"\"/g" gomt/core/info.go
exit 0
fi

export BUILD=`git rev-list HEAD -n 1 | cut -c 1-8 | tr "a-z" "A-Z"`
sed -i "s/version:\(\s*\)\".*\"\(\s*\),/version:\1\"$1\"\2,/g" app-v3/src/settings.js
sed -i "s/build:\(\s*\)\".*\"\(\s*\),/build:\1\"$BUILD\"\2,/g" app-v3/src/settings.js
sed -i "s/VERSION\(\s*\)=\(\s*\)\".*\"/VERSION\1=\2\"$1\"/g" gomt/core/info.go
sed -i "s/BUILD\(\s*\)=\(\s*\)\".*\"/BUILD\1=\2\"$BUILD\"/g" gomt/core/info.go

