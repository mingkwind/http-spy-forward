#!/usr/bin/env python
# -*- coding:utf-8 -*-
from matplotlib.font_manager import json_dump
import tornado.ioloop
import tornado.web
import json
import sys


class TestHandler(tornado.web.RequestHandler):
    def get(self, uri):
        httpReq = {}
        httpReq["method"] = str(self.request.method)
        httpReq["uri"] = str(self.request.uri)
        httpReq["headers"] = dict(self.request.headers)
        httpReq["body"] = str(self.request.body)
        print(httpReq)
        self.write(json.dumps(httpReq, indent=4, separators=(',', ':')))

    def post(self, uri):
        httpReq = {}
        httpReq["method"] = str(self.request.method)
        httpReq["uri"] = str(self.request.uri)
        httpReq["headers"] = dict(self.request.headers)
        httpReq["body"] = str(self.request.body)
        print(httpReq)
        self.write(json.dumps(httpReq, indent=4, separators=(',', ':')))


def main():

    application = tornado.web.Application([
        (r"/(.*?)", TestHandler),
    ])

    application.listen(8080)
    tornado.ioloop.IOLoop.instance().start()


if __name__ == "__main__":
    main()
