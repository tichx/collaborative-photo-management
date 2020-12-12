
const mongoose = require("mongoose");

const Schema = require('mongoose').Schema;

const channelSchema = new Schema({
    name: {type: String, required: true, unique: true},
    description: String,
    private: {type: Boolean, required: true},
    members: [],
    createdAt: {type: Date, default: Date.now},
    creator: {},
    editedAt: {type: Date, default: Date.now}
});

const messageSchema = new Schema({
    channelID: {type: String, required: true},
    body: {type: String, require: true},
    createdAt: {type: Date, default: Date.now},
    creator: {},
    editedAt: {type: Date, default: Date.now}
})

module.exports = {channelSchema, messageSchema}  