#!/bin/bash
echo "mode: set" > acc.out
for Dir in $(go list ./... | grep -v "vendor");
do
    returnval=`go test -coverprofile=profile.out $Dir`
    echo ${returnval}
    if [[ ${returnval} != *FAIL* ]]
    then
        if [ -f profile.out ]
        then
            cat profile.out | grep -v "mode: set" >> acc.out 
        fi
    else
        exit 1
    fi  

done
if [ -n "$COVERALLS_TOKEN" ]
then
    goveralls -coverprofile=acc.out -repotoken=$COVERALLS_TOKEN -service=travis-pro
fi  

rm -rf ./profile.out
rm -rf ./acc.out
