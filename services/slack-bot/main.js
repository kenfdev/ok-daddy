"use strict";

const fs = require("fs");
const request = require("request");
const { RtmClient, CLIENT_EVENTS, RTM_EVENTS } = require("@slack/client");
const { DADDYHOME_OAUTH_TOKEN_KEY_NAME, GOOGLEHOME_URL } = process.env;
const { CHANNEL_ID } = process.env;

// An access token (from your Slack app or custom integration - usually xoxb)
const token = fs
  .readFileSync(`/run/secrets/${DADDYHOME_OAUTH_TOKEN_KEY_NAME}`)
  .toString();

// TODO: This should be made dynamically (not hard coded)
const userMap = {
  U1SH5NPV3: "Dad"
};

// Cache of data
const appData = {};

// Initialize the RTM client with the recommended settings. Using the defaults for these
// settings is deprecated.
const rtm = new RtmClient(token, {
  dataStore: false,
  useRtmConnect: true
});

// The client will emit an RTM.AUTHENTICATED event on when the connection data is avaiable
// (before the connection is open)
rtm.on(CLIENT_EVENTS.RTM.AUTHENTICATED, connectData => {
  // Cache the data necessary for this app in memory
  appData.selfId = connectData.self.id;
  console.log(`Logged in as ${appData.selfId} of team ${connectData.team.id}`);
});

// The client will emit an RTM.RTM_CONNECTION_OPENED the connection is ready for
// sending and recieving messages
rtm.on(CLIENT_EVENTS.RTM.RTM_CONNECTION_OPENED, () => {
  console.log(`Ready`);
});

rtm.on(RTM_EVENTS.MESSAGE, message => {
  // For structure of `message`, see https://api.slack.com/events/message
  console.log("Incoming Message:", message);
  // Skip messages that are from a bot or my own user ID
  if (
    (message.subtype && message.subtype === "bot_message") ||
    (!message.subtype && message.user === appData.selfId)
  ) {
    return;
  }

  const regex = new RegExp(`^<@${appData.selfId}> (.*)`);
  const match = regex.exec(message.text);

  const userName = userMap[message.user];
  if (match != null && userName) {
    const msg = `Incoming message from ${userName}. ${match[1]}`;
    console.log("message", msg);
    request(
      {
        url: GOOGLEHOME_URL,
        method: "POST",
        json: true,
        body: {
          message: msg
        }
      },
      (err, res, body) => {
        if (err) {
          console.error(err);
          return;
        }
        console.log("request succeeded", res.body);
      }
    );
  }
});

// Start the connecting process
rtm.start();

// express
const express = require("express");
const bodyParser = require("body-parser");
const app = express();
app.use(bodyParser.json());

app.post("/api/message", (req, res) => {
  const msg = req.body.parameters.message;
  rtm.sendMessage(msg, CHANNEL_ID, err => {
    res.json({ result: "The message has been successfully sent." });
  });
});

app.listen(8080, () => console.log("Listening on port 8080!"));
