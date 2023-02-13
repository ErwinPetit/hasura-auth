#1/bin/sh

curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
  "email": "john.smith@nhost.io",
  "password": "Str0ngPassw#ord-94|%",
  "options": {
    "locale": "en",
    "defaultRole": "user",
    "allowedRoles": [
      "me",
      "user"
    ],
    "displayName": "John Smith",
    "metadata": {
      "firstName": "John",
      "lastName": "Smith"
    }
  }
}'  \
    http://localhost:4000/signup/email-password

 output=$(curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
      "email": "john.smith@nhost.io",
      "password": "Str0ngPassw#ord-94|%"
    }'\
        http://localhost:4000/signin/email-password)

token=$(echo $output | jq -r '.session.refreshToken')

hey -m POST \
    -n 500 \
    -H "Content-Type: application/json" \
    -d "{
  \"refreshToken\": \"$token\"
}" \
    http://localhost:4000/token
