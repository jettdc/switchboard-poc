from flask import Flask, request, Response, render_template, redirect, url_for

app = Flask(__name__)

@app.route("/")
def home():
    return render_template('base.html')

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
    id = 1
    token = request.args.get('auth_token')
    print ("token: ", token)
    if request.method == 'POST':
        if token == None:
            return redirect(url_for('login'))

        return redirect(url_for('track', id = id, token = token))
    return render_template('order.html', token=request.args.get('auth_token'))

@app.route("/track/<int:id>/<string:token>")
def track(id, token):
    return render_template('track.html', ID=id, TOKEN = token)

app.run(host="0.0.0.0", port=5000)


# this link is useful for url_for and redirect: https://stackoverflow.com/questions/26954122/how-can-i-pass-arguments-into-redirecturl-for-of-flask