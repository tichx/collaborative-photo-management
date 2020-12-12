const Schema = require('mongoose').Schema

const channelSchema = new Schema({
    id: {type:Schema.Types.ObjectId, required:true, unique:true, auto:true},
    name: {type:String, required:true, unique:true},
    description: String,
    private: Boolean,
    members: {type:[{_id: false,id:Number, email:String}]},
    createdAt: {type:Date, required:true},
    creator: {type:{id:Number, email:String}},
    editedAt: Date,
})

const messageSchema = new Schema({
    id: {type:Schema.Types.ObjectId, required:true, unique:true, auto:true},
    channelID: {type:Schema.Types.ObjectId, required:true},
    body: {type:String, required:true},
    createdAt: {type:Date, required:true},
    creator: {id:Number, email:String},
    editedAt: Date,
})


module.exports = {channelSchema, messageSchema}