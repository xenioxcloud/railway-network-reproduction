export oxylabs_username="oxylabs_username" \
    oxylabs_password="oxylabs_password" \
    oxylabs_entry="oxylabs_entry"

cd src
go build -o dev_build && ./dev_build
cd ..
