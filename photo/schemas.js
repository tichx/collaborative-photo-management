const Schema = require('mongoose').Schema

const photoSchema = new Schema({
    id: {type:Schema.Types.ObjectId, required:true, unique:true, auto:true},
    url: {type:String, required:true, unique:true},
    description: String,
    tags: {type:[{_id: false,id:Schema.Types.ObjectId}]},
    createdAt: {type:Date, required:true},
    creator: {type:{id:Number, email:String}},
    editedAt: Date,
})

const tagSchema = new Schema({
    id: {type:Schema.Types.ObjectId, required:true, unique:true, auto:true},
    name: {type:String, required:true},
    members: {type:[{_id: false,id:Number, email:String}]},
    // photos: {type:[{_id: false,id:Schema.Types.ObjectId, url: String}]},
    createdAt: {type:Date, required:true},
    creator: {id:Number, email:String},
    editedAt: Date,
})


module.exports = {photoSchema, tagSchema}