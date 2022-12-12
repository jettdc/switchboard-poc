from flask import Flask, request, Response, render_template, redirect, url_for
import requests
import random
import os
import json

app = Flask(__name__)

authAddr = os.getenv("AUTH_ADDRESS") or "localhost:8081"
switchboardAddr = os.getenv("SWITCHBOARD_ADDRESS") or "localhost:8080"
deliveryAddr = os.getenv("DELIVERY_ADDRESS") or "localhost:12345"
storeAddr = os.getenv("STORE_ADDRESS") or "localhost:54321"

@app.route("/")
def home():
    return redirect(url_for('order'))

@app.route("/login", methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        URL = f"http://{authAddr}/login"
        DATA = {
            'username': request.form['username'],
            'password': request.form['password']
        }
        token_resp = requests.post(url = URL, json = DATA)
        
        # Given invalid credentials, prompt user to login again
        if(token_resp.status_code != 200):
            return redirect(url_for('login'))
            
        json_response = json.loads(token_resp.content.decode())
        token = json_response['token']
        return redirect(url_for('order', auth_token = token))

    return render_template('login.html')

@app.route("/order", methods=['GET', 'POST'])
def order():
    id = random.randrange(1, 1000)
    token = request.args.get('auth_token')
    if request.method == 'POST':
        if not token:
            token = " "
        return redirect(url_for('track', id = id, token = token))
    return render_template('order.html', token=request.args.get('auth_token'))

@app.route("/track/<int:id>/<string:token>")
def track(id, token):
    # TODO: also pass host into template render
    os.system(f"curl {storeAddr}/store/{id}/events && curl {deliveryAddr}/delivery/{id}/events &")
    return render_template('track.html', id = id, token = token, switchboardAddr = switchboardAddr)

app.run(host="0.0.0.0", port=5000)


# this link is useful for url_for and redirect: https://stackoverflow.com/questions/26954122/how-can-i-pass-arguments-into-redirecturl-for-of-flask