from flask import Flask, request, Response, render_template

app = Flask(__name__)

@app.route("/")
def home():
    return render_template('base.html')

@app.route("/login")
def login():
    pass

@app.route("/order")
def order():
    pass

@app.route("/track/<int:id>")
def track(id):
    return render_template('track.html', ID=id)

app.run(host="0.0.0.0", port=5000)