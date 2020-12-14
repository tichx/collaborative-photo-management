import React, { Component } from 'react';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import Errors from '../../../Errors/Errors';

class UploadImage extends Component {
    constructor(props) {
        super(props);
        this.state = {
            file: null,
            error: ''
        }
    }

    // sendRequest = async (e) => {
    //     e.preventDefault();
    //     const { file } = this.state;
    //     let data = new FormData()
    //     data.append('uploadfile', file);
    //     const response = await fetch(api.base + api.handlers.upload, {
    //         method: "POST",
    //         body: data,
    //         headers: new Headers({
    //             "Authorization": localStorage.getItem("Authorization"),
    //         })
    //     });
    //     if (response.status >= 300) {
    //         const error = await response.text();
    //         console.log(error);
    //         this.setError(error);
    //         return;
    //     }
    //     alert("Upload Successful"); // TODO make this better by refactoring errors
    // }

    handleFile = (e) => {
        this.setState({
            file: e.target.files[0]
        })
    }

    setError = (error) => {
        this.setState({ error })
    }

    // submitPhoto = ()
    
    render() {
        const { firstName, lastName, error } = this.state;
        return <>
            <Errors error={error} setError={this.setError} />
            <div>Upload a new image</div>
            <form
                target="_blank"
                encType="multipart/form-data"
                action="https://api.xutiancheng.me/v1/upload/"
                method="post"
                >
                <input type="file" name="file" />
                <input type="submit" value="upload"/>   
                {/* <input type="submit" value="upload" onClick={} />    */}
            </form>
        </>
    }

}

export default UploadImage;