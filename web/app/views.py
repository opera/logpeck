# -*- encoding:utf-8 -*-

from flask import Flask
from flask import abort
from flask import json
from flask import make_response
from flask import redirect
from flask import render_template
from flask import request
from flask import url_for
from tinydb import TinyDB, Query
import base64
import datetime
import flask
import logging
import os
import re
import requests
import sys
import time

app = Flask(__name__)

if not app.debug:
    app.logger.addHandler(logging.StreamHandler())
    app.logger.setLevel(logging.INFO)

reload(sys)
sys.setdefaultencoding('utf8')

g_db = TinyDB('logpeck_db.json')
g_table_pecker = g_db.table('peckers')


def log(msg):
    prefix = '[' + time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time())) + '] '
    print(prefix + msg)


def http_post(url):
    log(url)
    try:
        r = requests.post(url=url, data='', headers={'Content-Type': 'text/html'})
        print(r.status_code, r.json())
        return r.json()
    except Exception,e:
        log('post err: ' + str(e))
        return ""
    sys.stdout.flush()


def logger(prefix):
    """
    Wall Clocks
    """
    def real_decorator(fn):
        """
        the real decorator
        """
        from functools import wraps
        # http://stackoverflow.com/a/309000/1498303
        @wraps(fn)
        def wrapper(*args, **kwargs):
            """
            the wrapper
            """
            # time.time() returns the time in seconds since the epoch as a floating point number
            start_timestamp = long(time.time() * 1000) # ms
            result = fn(*args, **kwargs)
            end_timestamp = long(time.time() * 1000)
            app.logger.info(datetime.datetime.fromtimestamp(time.time()).strftime('%Y-%m-%d %H:%M:%S')
                            + ' ' + prefix + ' ' + str(end_timestamp - start_timestamp) + ' ms @ ' + fn.__name__)
            return result
        return wrapper
    return real_decorator


@app.route('/list-pecktasks', methods=['POST'])
@logger('request')
def list_pecktasks():
    url = 'http://' + request.args['addr'] + '/peck_task/list'
    ret = http_post(url)
    return flask.jsonify(**ret)


@app.route('/list-peckers', methods=['POST'])
@logger('request')
def list_peckers():
    all_peckers = g_table_pecker.all()
    peckers = dict()
    for node in all_peckers:
        for k, v in node.iteritems():
            peckers[v] = True
    return flask.jsonify(**peckers)


@app.route('/add-pecker', methods=['POST'])
@logger('request')
def add_pecker():
    pecker = request.args['addr']
    g_table_pecker.insert({'pecker':pecker})
    return 'Add success'


@app.route('/remove-pecker', methods=['POST'])
@logger('request')
def remove_pecker():
    pecker = request.args['addr']
    q = Query()
    elements = g_table_pecker.search(q.pecker == pecker)
    eids = []
    for e in elements:
        eids.append(e.eid)
    g_table_pecker.remove(eids=eids)
    return "Remove finish"


@app.route('/pecker')
@logger('request')
def pecker():
    return render_template(
            'pecker.html',
            address=request.args['addr'])


@app.route('/')
@logger('request')
def index():
    return render_template('index.html')

if __name__ == "__main__":
    #url='http://127.0.0.1:7117/peck_task/list'
    #http_post(url)
    app.run(host = '0.0.0.0', port = 7119, debug = True)
