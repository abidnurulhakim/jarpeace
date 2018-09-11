 #!/bin/bash

# please run `ssh -R 80:localhost:{LOCAL_PORT} serveo.net`
# copy url "Forwarding HTTP traffic from {URL}" from output
# change "https://tersus.serveo.net/healthz" => URL that you get
# change "80:localhost:8888" => "80:localhost:{LOCAL_PORT}"

# this condition for check our tunnel is running. If not, it will start tunnel in background
if curl -s --head  --request GET https://tersus.serveo.net/healthz | grep "200" > /dev/null; then
	echo "OK"
else
	ssh -R 80:localhost:8888 serveo.net > /dev/null 2>&1 &
fi