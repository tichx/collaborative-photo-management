import React, { Component } from 'react';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import Errors from '../../../Errors/Errors';
import Alert from 'react-bootstrap/Alert';

class UploadImage extends Component {
    constructor(props) {
        super(props);
        this.state = {
            file: null,
            error: '',
            show:false
        }
    }

    sendRequest = async (e) => {
        e.preventDefault();
        const { file } = this.state;
        var fileName = file.name;
        let data = new FormData()
        data.append('file', file);
        const response = await fetch(api.base + api.handlers.upload, {
            method: "POST",
            body: data,
            headers: new Headers({
                // 'Content-Type': 'application/json',
                "Authorization": localStorage.getItem("Authorization"),
            })
        });
        if (response.status >= 300) {
            const error = await response.text();
            this.setError(error);
            return;
        } else if (response.status === 200) {
            fetch("https://api.xutiancheng.me/v1/photos", {
                method: "POST",
                body: JSON.stringify({"url": "https://image-441.s3.amazonaws.com/"+fileName}),
                headers: new Headers({
                    'Content-Type': 'application/json',
                    "Authorization": localStorage.getItem("Authorization"),
                })
            })
            .then(resp =>{
                if (resp.status == 201){
                    const show = true;
                    this.setShow(show);
                } else{
                    const error = resp.text();
                    this.setError(error);
            return;
                }
            });
        }
    }

    handleFile = (e) => {
        this.setState({
            file: e.target.files[0]
        });
    }

    setError = (error) => {
        this.setState({ error })
    }

    setShow = (show) => {
        this.setState({ show })
    }

    
    render() {
        const { firstName, lastName, error, show } = this.state;
        if (show) {
        return <>
            <Errors error={error} setError={this.setError} />
            <div>Upload a new image</div>
            <form
                target="_blank"
                encType="multipart/form-data"
                // action="https://api.xutiancheng.me/v1/upload/"
                method="post"
                >
                <input type="file" name="file" onChange={(e) => this.handleFile(e)}/>
                {/* <input type="submit" value="upload"/>    */}
                <input type="submit" value="upload" onClick={(e) => this.sendRequest(e)} />   
            </form>
            <Alert variant="danger" onClose={() => this.setShow(false)} dismissible>
            <Alert.Heading>photo uploaded</Alert.Heading>
            <p>
                to check the photo, return to main page, then go to
                the photo management page.
            </p>
          </Alert>
        </>
        } else {
            return <>
            <Errors error={error} setError={this.setError} />
            <div>Upload a new image</div>
            <form
                target="_blank"
                encType="multipart/form-data"
                // action="https://api.xutiancheng.me/v1/upload/"
                method="post"
                >
                <input type="file" name="file" onChange={(e) => this.handleFile(e)}/>
                {/* <input type="submit" value="upload"/>    */}
                <input type="submit" value="upload" onClick={(e) => this.sendRequest(e)} />   
            </form>
        </>
        }
    }

}

export default UploadImage;