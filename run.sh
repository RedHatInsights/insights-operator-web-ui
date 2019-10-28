go build

if [ $? -eq 0 ]
then
    echo "Build ok"
    ./insights-operator-web-ui
else
    echo "Build failed"
fi
