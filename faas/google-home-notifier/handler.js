"use strict";

const googlehome = require("google-home-notifier");
// const language = "ja";
// const language = "us";
const { GOOGLEHOME_IP } = process.env;

googlehome.ip(GOOGLEHOME_IP);
googlehome.device("Google-Home");


module.exports = (context, callback) => {
  const json = JSON.parse(context);
  googlehome.notify(json.message, res => {
    callback(undefined, { status: "ok" });
  });
};
