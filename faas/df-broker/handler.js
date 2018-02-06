"use strict";

/**
https://github.com/dialogflow/dialogflow-nodejs-client
 */

const fs = require("fs");
const apiai = require("apiai");
const uuidv4 = require("uuid/v4");
const { DF_CLIENT_ACCESS_TOKEN_KEY_NAME } = process.env;

const token = fs
  .readFileSync(`/run/secrets/${DF_CLIENT_ACCESS_TOKEN_KEY_NAME}`)
  .toString();

var app = apiai(token);

module.exports = (context, callback) => {
  const json = JSON.parse(context);
  const uuid = uuidv4();
  var request = app.textRequest(json.text, {
    sessionId: uuid
  });

  request.on("response", function(response) {
    const { result } = response;

    callback(undefined, {
      metadata: result.metadata,
      fulfillment: result.fulfillment
    });
  });

  request.on("error", function(error) {
    console.error(error);
  });

  request.end();
};
