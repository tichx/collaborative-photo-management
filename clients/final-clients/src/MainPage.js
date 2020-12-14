import React, { useEffect, useState } from 'react';
import ReactDOM from 'react-dom';

const MainPage = (props) => {
    //data structure for getting images
    const [imgDataList,setImgDataList] = useState([]);
    //data structure for getting tags
    const [tagDataList,setTagDataList] = useState([]);
    // const [tagText, setTagText] = useState("");
    const [tag,setTag] = useState(-1);

    // Fetch all img data
    useEffect(() => {
        fetch("https://api.xutiancheng.me/v1/photos")
        .then(resp => resp.json())
        .then(data => {setImgDataList(data)})
    });
    // Fetch all tag data
    useEffect(() => {
        fetch("https://api.xutiancheng.me/v1/tags")
        .then(resp => resp.json())
        .then(data => 
        {
            setTagDataList(data)
        })


    });


    
    return (
    <div>
        <TagButtonList tags={tagDataList} onClick={setTag}/>
        <ImgCardList  tag={tag} imgDataList={imgDataList}/>
    </div>
    )
    
}

//populates the tag buttons into a list
const TagButtonList = (props)=>{
    const TagLists = props.tags.map(data=><TagButton key={data.id} tag={data} onClick={props.onClick}/>)
}

//helper that maps tagid with tag object so that when photos are populated,
//each photo component can search this map provided the list of photo's tag id
/*
const TagTable = (props) => {
    const Table = props.tags.map(data => <Map key={data.id} tag={data}/>)
}
*/

//gets the new tag name typed by user, pass it to /v1/tags
/*const TagText = () => {
    return <input
    type="text"
    placeholder="tag name..."
    onChange={event => {this.setState({tagText: event.target.value})}}
    onKeyPress={event => {
                if (event.key === 'Enter') {
                  this.search()
                }
              }}
/>
}*/

//individual button showing each tag
const TagButton = (props)=>{
    return <button onClick={props.onClick(props.tag.id)}>props.tag.name</button>
}

//populates the image cards into a list
const ImgCardList = (props) => {
    var photoListItems;
    if (props.tag === -1) {
        photoListItems = props.imgDataList.map((data) => {
            <ImgCard key={data.id} img={data}/>
        });
    } else {
        photoListItems = props.imgDataList.map((data) => {
            if (data.tags.includes(props.tag)){
                <ImgCard key={data.id} img={data}/>
            }
        });
    }
    
    return (photoListItems)
}

//individual image card, still missing function that when clicked,
//lets user add a tag for this photo
const ImgCard = (props) =>{
    
    
    
    return(<div>
        <img src={props.img.url}/>
        <button onClick={}>Add Tag</button>
     </div>)
    
    
}

//TODO: constructs a map with tag id as the key that points to the 
//corresponding tag object as the value. this table is used to search tag ID and show
//tag name under relevant image.
const MapTagIDName = ()=>{
    
}

//TODO: search tag names with each image so that tags can show
//under each image.
//Note that not all image has a tag
const SearchImageTag = ()=>{
    
}

//TODO: bind image to a new tag name input by user and post 
//it to /v1/tags with payload: {"name": "tagname"}
const BindTagImg = ()=>{
    
}


