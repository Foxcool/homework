import os
import json

from flask import Flask, request, abort
import requests

app = Flask(__name__)

config = {
    'PROFILE_PATH': os.environ.get('PROFILE_PATH', ''),
    'PRIVATE_KEY': os.environ.get('PRIVATE_KEY', ''),
    'PUBLIC_KEY': os.environ.get('PUBLIC_KEY', ''),
}

print(config['PRIVATE_KEY'])
print(config['PUBLIC_KEY'])


def get_user_by_credentials(login, password):
    body = {
        "login": login,
        "password": password,
    }
    headers = {'Content-type': 'application/json', 'Accept': 'application/json'}
    r = requests.post('http://' + config['PROFILE_PATH'] + "/auth", json=body, headers=headers)

    data = r.json()

    if r.status_code != 200:
        abort(r.status_code, data['message'])

    userID = data['data'][0]['userID']

    headers = {'X-User-Id': userID}
    r = requests.get('http://' + config['PROFILE_PATH'] + '/users/' + userID, headers=headers)
    data = r.json()

    return data['data'][0]


def create_id_token(user_info):
    import jwt
    import datetime
    data = {
        "iss": "http://arch.homework",
        "exp": datetime.datetime.utcnow() + datetime.timedelta(minutes=15),
        "sub": user_info["ID"],
        "email": user_info["email"],
        "given_name": user_info["firstName"],
        "family_name": user_info["lastName"]
    }

    encoded = jwt.encode(data, config['PRIVATE_KEY'], algorithm='RS256', headers={'kid': '1'})
    return encoded.decode('utf-8')


@app.route("/login", methods=["POST"])
def login():
    request_data = request.get_json()

    user_info = get_user_by_credentials(request_data['login'], request_data['password'])
    if user_info:
        id_token = create_id_token(user_info)
        response = app.make_response({"IDtoken": id_token})
        return response
    else:
        abort(401)


@app.route("/.well-known/jwks.json")
def jwks():
    from authlib.jose import JsonWebKey
    from authlib.jose import JWK_ALGORITHMS
    jwk = JsonWebKey(algorithms=JWK_ALGORITHMS)
    key = jwk.dumps(config['PUBLIC_KEY'], kty='RSA')
    key['kid'] = '1'
    return {"keys": [key]}


@app.route("/health")
def health():
    return {"status": "OK"}


if __name__ == "__main__":
    app.run(host='0.0.0.0', port='80', debug=True)
