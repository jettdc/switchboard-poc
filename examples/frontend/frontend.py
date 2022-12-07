from flask import Flask, request, Response, render_template, redirect, url_for
import random
import os

app = Flask(__name__)

@app.route("/")
def home():
    return redirect(url_for('order'))

@app.route("/login", methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        print (request.form['username'])# - verified
        print (request.form['password'])
        # call loginService and pass in the above two, saves token
        temp_token = "200"
        return redirect(url_for('order', auth_token = temp_token)) #also return token

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
    # TODO: get host from env variable and insert into f string
    # TODO: also pass host into template render
    os.system(f"curl localhost:54321/store/{id}/events && curl localhost:12345/delivery/{id}/events &")
    return render_template('track.html', id = id, token = token)

app.run(host="0.0.0.0", port=5000)


# this link is useful for url_for and redirect: https://stackoverflow.com/questions/26954122/how-can-i-pass-arguments-into-redirecturl-for-of-flask